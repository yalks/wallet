package wallet

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/shopspring/decimal"

	"github.com/yalks/wallet/constants"
	"github.com/yalks/wallet/entity"
	"github.com/yalks/wallet/logic"
)

// transactionManager 实现事务管理器接口
type transactionManager struct {
	mu            sync.RWMutex
	logic         *logic.SharedLogicContext
	tokenLogic    logic.ITokenLogic
	userLogic     logic.IUserLogic
	balanceLogic  logic.IBalanceLogic
	feeCalculator *logic.FeeCalculator
	validator     *logic.TransactionValidator
}

// NewTransactionManager 创建事务管理器
func NewTransactionManager() ITransactionManager {
	return &transactionManager{
		logic:         logic.GetSharedContext(),
		tokenLogic:    logic.NewTokenLogic(),
		userLogic:     logic.NewUserLogic(),
		balanceLogic:  logic.NewBalanceLogic(),
		feeCalculator: logic.NewFeeCalculator(),
		validator:     logic.NewTransactionValidator(),
	}
}

// CreateTransaction 创建交易记录
func (tm *transactionManager) CreateTransaction(ctx context.Context, req *constants.TransactionRequest) (int64, error) {
	// 参数验证
	if err := tm.validateTransactionRequest(req); err != nil {
		return 0, gerror.Wrap(err, "交易请求验证失败")
	}

	// 获取数据库连接
	db := g.DB()
	var transactionID int64

	// 开启事务
	err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 检查幂等性
		existingTx, err := tm.getTransactionByReference(ctx, tx, req.Reference)
		if err != nil {
			return gerror.Wrap(err, "检查交易幂等性失败")
		}
		if existingTx != nil {
			transactionID = int64(existingTx.TransactionId)
			g.Log().Infof(ctx, "交易已存在，返回现有交易: Reference=%s, TransactionID=%d", req.Reference, transactionID)
			return nil
		}

		// 获取用户和代币信息
		user, err := tm.userLogic.GetUserByID(ctx, uint64(req.UserID))
		if err != nil {
			return gerror.Wrapf(err, "获取用户信息失败: UserID=%d", req.UserID)
		}

		token, err := tm.tokenLogic.GetTokenByID(ctx, uint(req.TokenID))
		if err != nil {
			return gerror.Wrapf(err, "获取代币信息失败: TokenID=%d", req.TokenID)
		}

		// 获取当前余额
		currentBalance, _, err := tm.balanceLogic.GetBalance(ctx, uint64(req.UserID), token.Symbol)
		if err != nil {
			return gerror.Wrap(err, "获取当前余额失败")
		}

		// 转换金额
		amount, err := decimal.NewFromString(req.Amount)
		if err != nil {
			return gerror.Wrapf(err, "无效的金额格式: %s", req.Amount)
		}

		// 计算新余额
		var newBalance decimal.Decimal
		direction := constants.GetFundDirection(req.FundType)
		switch direction {
		case constants.FundDirectionIn:
			newBalance = currentBalance.Add(amount)
		case constants.FundDirectionOut:
			newBalance = currentBalance.Sub(amount)
			// 检查余额是否充足
			if newBalance.LessThan(decimal.Zero) {
				return gerror.Newf("余额不足: 当前余额=%s, 需要金额=%s", currentBalance.String(), amount.String())
			}
		default:
			return gerror.Newf("未知的资金方向: %s", direction)
		}

		// 创建交易记录
		transaction := &entity.Transactions{
			UserId:            uint(req.UserID),
			TokenId:           uint(req.TokenID),
			Amount:            amount,
			BalanceBefore:     currentBalance,
			BalanceAfter:      newBalance,
			Type:              string(req.FundType),
			Status:            1, // 1表示成功
			Direction:         string(direction),
			Memo:              req.Description,
			Symbol:            token.Symbol,
			WalletType:        "available",
			BusinessId:        req.Reference,
			RelatedEntityId:   uint64(req.RelatedID),
			RelatedEntityType: "fund_operation",
			CreatedAt:         gtime.Now(),
			UpdatedAt:         gtime.Now(),
			
			// New fields - User request information
			RequestAmount:    amount, // Using the same amount as the transaction amount
			RequestReference: req.Reference,
			RequestMetadata:  tm.convertMetadataToJSON(req.Metadata),
			RequestSource:    tm.validator.SanitizeRequestSource(req.RequestSource),
			RequestIp:        req.RequestIP,
			RequestUserAgent: req.RequestUserAgent,
			RequestTimestamp: gtime.Now(),
			ProcessedAt:      gtime.Now(),
			
			// Calculate fee based on fund type
			FeeAmount: tm.calculateFee(req.FundType, amount),
			FeeType:   tm.getFeeType(req.FundType),
			
			// Exchange rate (set to 1 if no conversion)
			ExchangeRate: decimal.NewFromInt(1),
			
			// Target user fields for transfers
			TargetUserId:   uint(req.TargetUserID),
			TargetUsername: req.TargetUsername,
		}

		// 插入交易记录
		id, err := tm.logic.GetTransactionDAO().CreateTransaction(ctx, tx, transaction)
		if err != nil {
			return gerror.Wrap(err, "创建交易记录失败")
		}
		transactionID = id

		// 更新本地余额
		err = tm.balanceLogic.UpdateLocalBalance(ctx, tx, uint64(req.UserID), token.Symbol, newBalance, decimal.Zero)
		if err != nil {
			return gerror.Wrap(err, "更新本地余额失败")
		}

		// 执行远程钱包操作
		err = tm.executeRemoteOperation(ctx, user, token, amount, req.FundType, req.Metadata)
		if err != nil {
			return gerror.Wrap(err, "执行远程钱包操作失败")
		}

		g.Log().Infof(ctx, "交易创建成功: TransactionID=%d, UserID=%d, Amount=%s %s, FundType=%s",
			transactionID, req.UserID, amount.String(), token.Symbol, req.FundType)

		return nil
	})

	if err != nil {
		return 0, err
	}

	return transactionID, nil
}

// GetTransactionByID 根据ID获取交易记录
func (tm *transactionManager) GetTransactionByID(ctx context.Context, transactionID int64) (*TransactionRecord, error) {
	// 查询交易记录
	var tx entity.Transactions
	err := g.Model("transactions").Ctx(ctx).
		Where("transaction_id = ?", transactionID).
		Scan(&tx)
	if err != nil {
		return nil, gerror.Wrapf(err, "获取交易记录失败: TransactionID=%d", transactionID)
	}
	if tx.TransactionId == 0 {
		return nil, gerror.Newf("交易记录不存在: TransactionID=%d", transactionID)
	}

	return tm.convertToTransactionRecord(&tx), nil
}

// GetTransactionByReference 根据引用获取交易记录
func (tm *transactionManager) GetTransactionByReference(ctx context.Context, reference string) (*TransactionRecord, error) {
	tx, err := tm.getTransactionByReference(ctx, nil, reference)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, nil
	}

	return tm.convertToTransactionRecord(tx), nil
}

// UpdateTransactionStatus 更新交易状态
func (tm *transactionManager) UpdateTransactionStatus(ctx context.Context, transactionID int64, status constants.TransactionStatus) error {
	// 验证状态是否有效
	if !constants.IsValidTransactionStatus(status) {
		return gerror.Newf("无效的交易状态: %s", status)
	}

	// 获取现有交易
	var existingTx entity.Transactions
	err := g.Model("transactions").Ctx(ctx).
		Where("transaction_id = ?", transactionID).
		Scan(&existingTx)
	if err != nil {
		return gerror.Wrap(err, "获取交易记录失败")
	}
	if existingTx.TransactionId == 0 {
		return gerror.Newf("交易记录不存在: TransactionID=%d", transactionID)
	}

	// 检查当前状态是否为最终状态
	currentStatus := constants.TransactionStatus(fmt.Sprintf("%d", existingTx.Status))
	if constants.IsFinalStatus(currentStatus) {
		return gerror.Newf("交易已处于最终状态，无法更新: CurrentStatus=%s", currentStatus)
	}

	// 更新状态
	_, err = g.Model("transactions").Ctx(ctx).
		Where("transaction_id = ?", transactionID).
		Update(g.Map{
			"status":     statusToInt(status),
			"updated_at": gtime.Now(),
		})

	if err != nil {
		return gerror.Wrapf(err, "更新交易状态失败: TransactionID=%d, Status=%s", transactionID, status)
	}

	g.Log().Infof(ctx, "交易状态更新成功: TransactionID=%d, OldStatus=%s, NewStatus=%s",
		transactionID, currentStatus, status)

	return nil
}

// GetUserTransactionHistory 获取用户交易历史
func (tm *transactionManager) GetUserTransactionHistory(ctx context.Context, userID int64, fundType constants.FundType, limit, offset int) ([]*TransactionRecord, error) {
	model := g.Model("transactions").Ctx(ctx).
		Where("user_id = ?", userID).
		OrderDesc("created_at")

	// 如果指定了资金类型，添加过滤条件
	if fundType != "" && constants.IsValidFundType(fundType) {
		model = model.Where("type = ?", string(fundType))
	}

	// 添加分页
	if limit > 0 {
		model = model.Limit(limit)
	}
	if offset > 0 {
		model = model.Offset(offset)
	}

	var transactions []*entity.Transactions
	err := model.Scan(&transactions)
	if err != nil {
		return nil, gerror.Wrap(err, "查询交易历史失败")
	}

	// 转换为交易记录
	records := make([]*TransactionRecord, 0, len(transactions))
	for _, tx := range transactions {
		records = append(records, tm.convertToTransactionRecord(tx))
	}

	return records, nil
}

// validateTransactionRequest 验证交易请求
func (tm *transactionManager) validateTransactionRequest(req *constants.TransactionRequest) error {
	if req.UserID == 0 {
		return gerror.New("用户ID不能为空")
	}
	if req.WalletID == 0 {
		return gerror.New("钱包ID不能为空")
	}
	if req.TokenID == 0 {
		return gerror.New("代币ID不能为空")
	}
	if req.Amount == "" || req.Amount == "0" {
		return gerror.New("金额必须大于0")
	}
	if !constants.IsValidFundType(req.FundType) {
		return gerror.Newf("无效的资金类型: %s", req.FundType)
	}
	if req.Reference == "" {
		return gerror.New("交易引用不能为空")
	}
	
	// 验证新字段
	if req.RequestSource != "" && !tm.validator.ValidateRequestSource(req.RequestSource) {
		return gerror.Newf("无效的请求来源: %s (允许的值: %v)", req.RequestSource, tm.validator.GetAllowedSources())
	}
	if req.RequestIP != "" && !tm.validator.ValidateIP(req.RequestIP) {
		return gerror.Newf("无效的IP地址格式: %s", req.RequestIP)
	}
	if req.FeeType != "" && !tm.validator.ValidateFeeType(req.FeeType) {
		return gerror.Newf("无效的手续费类型: %s (允许的值: fixed, percentage)", req.FeeType)
	}
	if req.RequestUserAgent != "" && !tm.validator.ValidateUserAgent(req.RequestUserAgent) {
		return gerror.New("User-Agent过长，最大长度为1000字符")
	}

	return nil
}

// getTransactionByReference 根据引用获取交易（内部方法）
func (tm *transactionManager) getTransactionByReference(ctx context.Context, tx gdb.TX, reference string) (*entity.Transactions, error) {
	var transaction entity.Transactions

	model := g.Model("transactions").Ctx(ctx).Where("reference = ? OR business_id = ?", reference, reference)
	if tx != nil {
		model = model.TX(tx)
	}

	err := model.Scan(&transaction)
	if err != nil {
		return nil, gerror.Wrapf(err, "查询交易记录失败: Reference=%s", reference)
	}

	if transaction.TransactionId == 0 {
		return nil, nil
	}

	return &transaction, nil
}

// executeRemoteOperation 执行远程钱包操作
func (tm *transactionManager) executeRemoteOperation(ctx context.Context, user *entity.Users, token *entity.Tokens, amount decimal.Decimal, fundType constants.FundType, metadata map[string]interface{}) error {
	walletSDK := tm.logic.GetWalletSDK()
	if walletSDK == nil {
		return gerror.New("钱包SDK未初始化")
	}

	// 获取主钱包ID
	if user.MainWalletId == "" {
		return gerror.Newf("用户没有主钱包: UserID=%d", user.Id)
	}

	// 转换金额为原始格式（验证格式正确性）
	_, err := tm.tokenLogic.ConvertBalanceToRaw(ctx, amount, token.Symbol)
	if err != nil {
		return gerror.Wrap(err, "转换金额格式失败")
	}

	// 构建SDK元数据
	sdkMetadata := make(map[string]string)
	for k, v := range metadata {
		sdkMetadata[k] = fmt.Sprintf("%v", v)
	}
	sdkMetadata["fund_type"] = string(fundType)

	// 根据资金方向执行操作
	direction := constants.GetFundDirection(fundType)

	// 直接调用logic层的operation逻辑，避免重复实现SDK调用
	operationReq := &logic.FinancialOperationRequest{
		UserID:      uint64(user.Id),
		TokenSymbol: token.Symbol,
		Amount:      amount,
		BusinessID:  fmt.Sprintf("%s_%d", fundType, gtime.Now().UnixNano()),
		Description: fmt.Sprintf("Transaction for %s", fundType),
		Metadata:    sdkMetadata,
	}

	// 根据方向设置操作类型
	switch direction {
	case constants.FundDirectionIn:
		operationReq.OperationType = logic.OperationTypeCredit
	case constants.FundDirectionOut:
		operationReq.OperationType = logic.OperationTypeDebit
	default:
		return gerror.Newf("未知的资金方向: %s", direction)
	}

	// 使用现有的operation logic执行远程操作
	operationLogic := logic.NewOperationLogic()
	_, err = operationLogic.ExecuteInTx(ctx, nil, operationReq)
	if err != nil {
		return gerror.Wrap(err, "执行远程钱包操作失败")
	}

	return nil
}

// convertToTransactionRecord 转换实体为交易记录
func (tm *transactionManager) convertToTransactionRecord(tx *entity.Transactions) *TransactionRecord {
	record := &TransactionRecord{
		ID:          int64(tx.TransactionId),
		UserID:      int64(tx.UserId),
		WalletID:    0, // 钱包ID从实体中移除了
		TokenID:     int64(tx.TokenId),
		Amount:      tx.Amount.String(),
		FundType:    constants.FundType(tx.Type),
		Direction:   constants.FundDirection(tx.Direction),
		Status:      intToStatus(tx.Status),
		Reference:   tx.BusinessId, // 使用BusinessId作为Reference
		Description: tx.Memo,
		RelatedID:   int64(tx.RelatedEntityId),
		CreatedAt:   tx.CreatedAt.String(),
		UpdatedAt:   tx.UpdatedAt.String(),
	}

	// 元数据暂时为空，因为entity中没有Metadata字段
	record.Metadata = make(map[string]interface{})

	return record
}

// statusToInt 将状态转换为整数
func statusToInt(status constants.TransactionStatus) int {
	switch status {
	case constants.TransactionStatusCompleted:
		return 1
	case constants.TransactionStatusFailed, constants.TransactionStatusCancelled, constants.TransactionStatusExpired, constants.TransactionStatusRefunded:
		return 0
	default:
		return 0
	}
}

// intToStatus 将整数转换为状态
func intToStatus(status uint) constants.TransactionStatus {
	switch status {
	case 1:
		return constants.TransactionStatusCompleted
	case 0:
		return constants.TransactionStatusFailed
	default:
		return constants.TransactionStatusFailed
	}
}

// convertMetadataToJSON 将元数据map转换为JSON字符串
func (tm *transactionManager) convertMetadataToJSON(metadata map[string]interface{}) string {
	if metadata == nil || len(metadata) == 0 {
		return "{}"
	}
	
	data, err := json.Marshal(metadata)
	if err != nil {
		g.Log().Errorf(context.Background(), "转换元数据到JSON失败: %v", err)
		return "{}"
	}
	
	return string(data)
}

// calculateFee 计算手续费
func (tm *transactionManager) calculateFee(fundType constants.FundType, amount decimal.Decimal) decimal.Decimal {
	fee, _ := tm.feeCalculator.CalculateFee(fundType, amount)
	return fee
}

// getFeeType 获取手续费类型
func (tm *transactionManager) getFeeType(fundType constants.FundType) string {
	_, feeType := tm.feeCalculator.CalculateFee(fundType, decimal.Zero)
	return feeType
}
