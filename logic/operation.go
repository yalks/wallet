package logic

import (
	"context"
	"encoding/json"
	"fmt"

	// sdkTypes "github.com/a19ba14d/ledger-wallet-sdk/pkg/types" // 暂时注释，不使用远程 SDK
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/shopspring/decimal"

	"github.com/yalks/wallet/entity"
)

// OperationType 操作类型
type OperationType string

const (
	OperationTypeCredit   OperationType = "credit"   // 增加资金
	OperationTypeDebit    OperationType = "debit"    // 减少资金
	OperationTypeTransfer OperationType = "transfer" // 转账
)

// FinancialOperationRequest 财务操作请求
type FinancialOperationRequest struct {
	UserID        uint64            `json:"user_id"`
	TokenSymbol   string            `json:"token_symbol"`
	Amount        decimal.Decimal   `json:"amount"`
	OperationType OperationType     `json:"operation_type"`
	BusinessID    string            `json:"business_id"`
	Description   string            `json:"description"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	FeeAmount     decimal.Decimal   `json:"fee_amount,omitempty"` // Fee amount (trusted from request)
	FeeType       string            `json:"fee_type,omitempty"`   // Fee type (fixed, percentage)
}

// FinancialOperationResult 财务操作结果
type FinancialOperationResult struct {
	TransactionID  int64           `json:"transaction_id"`
	BalanceBefore  decimal.Decimal `json:"balance_before"`
	BalanceAfter   decimal.Decimal `json:"balance_after"`
	WalletResponse any             `json:"wallet_response,omitempty"`
}

// IOperationLogic 操作业务逻辑接口
type IOperationLogic interface {
	// ExecuteInTx 在事务中执行财务操作
	ExecuteInTx(ctx context.Context, tx gdb.TX, req *FinancialOperationRequest) (*FinancialOperationResult, error)
	// GetUserBalance 获取用户余额
	GetUserBalance(ctx context.Context, userID uint64, tokenSymbol string) (decimal.Decimal, error)
	// ValidateOperation 验证操作是否可执行
	ValidateOperation(ctx context.Context, req *FinancialOperationRequest) error
}

type operationLogic struct {
	userLogic    IUserLogic
	tokenLogic   ITokenLogic
	balanceLogic IBalanceLogic
	context      *SharedLogicContext
}

// NewOperationLogic 创建操作业务逻辑实例
func NewOperationLogic() IOperationLogic {
	return &operationLogic{
		userLogic:    NewUserLogic(),
		tokenLogic:   NewTokenLogic(),
		balanceLogic: NewBalanceLogic(),
		context:      GetSharedContext(),
	}
}

// ExecuteInTx 在事务中执行财务操作
func (l *operationLogic) ExecuteInTx(ctx context.Context, tx gdb.TX, req *FinancialOperationRequest) (*FinancialOperationResult, error) {
	// 1. 参数验证
	if err := l.validateRequest(req); err != nil {
		return nil, gerror.Wrap(err, "参数验证失败")
	}

	// 2. 幂等性检查
	if existingTx, err := l.context.GetTransactionDAO().GetTransactionByBusinessID(ctx, req.BusinessID); err != nil {
		return nil, gerror.Wrap(err, "幂等性检查失败")
	} else if existingTx != nil {
		g.Log().Infof(ctx, "幂等性检查: 操作已存在 BusinessID=%s, TransactionID=%d", req.BusinessID, existingTx.TransactionId)
		return &FinancialOperationResult{
			TransactionID: int64(existingTx.TransactionId),
			BalanceBefore: existingTx.BalanceBefore,
			BalanceAfter:  existingTx.BalanceAfter,
		}, nil
	}

	// 3. 业务验证
	if err := l.ValidateOperation(ctx, req); err != nil {
		return nil, err
	}

	// 4. 获取操作前余额
	balanceBefore, _, err := l.balanceLogic.GetBalance(ctx, req.UserID, req.TokenSymbol)
	if err != nil {
		return nil, gerror.Wrap(err, "获取操作前余额失败")
	}

	// 5. 计算操作后余额
	var balanceAfter decimal.Decimal
	switch req.OperationType {
	case OperationTypeCredit:
		balanceAfter = balanceBefore.Add(req.Amount)
	case OperationTypeDebit:
		balanceAfter = balanceBefore.Sub(req.Amount)
	default:
		return nil, gerror.Newf("不支持的操作类型: %s", req.OperationType)
	}

	// 6. 执行远程钱包操作 - 暂时注释，使用本地事务
	// walletResponse, err := l.executeRemoteWalletOperation(ctx, req, balanceAfter)
	// if err != nil {
	// 	return nil, gerror.Wrap(err, "执行远程钱包操作失败")
	// }
	var walletResponse any = nil // 暂时设为nil

	// 7. 更新本地余额 - 提前到远程操作位置，确保事务原子性
	err = l.balanceLogic.UpdateLocalBalance(ctx, tx, req.UserID, req.TokenSymbol, balanceAfter, decimal.Zero)
	if err != nil {
		return nil, gerror.Wrap(err, "更新本地余额失败")
	}

	// 8. 创建交易记录
	transactionID, err := l.createTransactionRecord(ctx, tx, req, balanceBefore, balanceAfter)
	if err != nil {
		return nil, gerror.Wrap(err, "创建交易记录失败")
	}

	// 回滚逻辑暂时注释，因为使用本地事务
	// if err != nil {
	// 	g.Log().Errorf(ctx, "更新本地余额失败: %v", err)
	// 	// 尝试回滚远程钱包操作
	// 	rollbackErr := l.rollbackRemoteWalletOperation(ctx, req, balanceBefore, walletResponse)
	// 	if rollbackErr != nil {
	// 		g.Log().Errorf(ctx, "回滚远程钱包操作失败: %v", rollbackErr)
	// 	}
	// 	return nil, gerror.Wrap(err, "更新本地余额失败，操作已回滚")
	// }

	result := &FinancialOperationResult{
		TransactionID:  transactionID,
		BalanceBefore:  balanceBefore,
		BalanceAfter:   balanceAfter,
		WalletResponse: walletResponse,
	}

	g.Log().Infof(ctx, "财务操作成功: BusinessID=%s, TransactionID=%d, UserID=%d, Amount=%s %s",
		req.BusinessID, transactionID, req.UserID, req.Amount.String(), req.TokenSymbol)

	return result, nil
}

// GetUserBalance 获取用户余额
func (l *operationLogic) GetUserBalance(ctx context.Context, userID uint64, tokenSymbol string) (decimal.Decimal, error) {
	availableBalance, _, err := l.balanceLogic.GetBalance(ctx, userID, tokenSymbol)
	if err != nil {
		return decimal.Zero, gerror.Wrapf(err, "获取用户余额失败: UserID=%d, Symbol=%s", userID, tokenSymbol)
	}
	return availableBalance, nil
}

// ValidateOperation 验证操作是否可执行
func (l *operationLogic) ValidateOperation(ctx context.Context, req *FinancialOperationRequest) error {
	// 1. 验证用户是否存在
	user, err := l.userLogic.GetUserByID(ctx, req.UserID)
	if err != nil {
		return gerror.Wrapf(err, "用户验证失败: UserID=%d", req.UserID)
	}

	// 2. 验证代币是否存在
	_, err = l.tokenLogic.GetTokenBySymbol(ctx, req.TokenSymbol)
	if err != nil {
		return gerror.Wrapf(err, "代币验证失败: Symbol=%s", req.TokenSymbol)
	}

	// 3. 检测和创建钱包（远程和本地）
	err = l.ensureWalletExists(ctx, user, req.TokenSymbol)
	if err != nil {
		return gerror.Wrap(err, "钱包检测和创建失败")
	}

	// 4. 校验余额一致性
	// err = l.validateBalanceConsistency(ctx, req.UserID, req.TokenSymbol)
	// if err != nil {
	// 	return gerror.Wrap(err, "余额一致性校验失败")
	// }

	// 5. 对于扣款操作，检查余额是否充足
	if req.OperationType == OperationTypeDebit {
		balance, err := l.GetUserBalance(ctx, req.UserID, req.TokenSymbol)
		if err != nil {
			return gerror.Wrap(err, "获取余额失败")
		}

		if balance.LessThan(req.Amount) {
			return gerror.Newf("余额不足: 当前余额=%s, 需要金额=%s", balance.String(), req.Amount.String())
		}
	}

	return nil
}

// validateRequest 验证请求参数
func (l *operationLogic) validateRequest(req *FinancialOperationRequest) error {
	if req.UserID == 0 {
		return gerror.New("用户ID不能为空")
	}
	if req.TokenSymbol == "" {
		return gerror.New("代币符号不能为空")
	}
	if req.Amount.LessThanOrEqual(decimal.Zero) {
		return gerror.New("金额必须大于0")
	}
	if req.BusinessID == "" {
		return gerror.New("业务ID不能为空")
	}
	return nil
}

// executeRemoteWalletOperation 执行远程钱包操作 - 暂时注释，使用本地事务
/*
func (l *operationLogic) executeRemoteWalletOperation(ctx context.Context, req *FinancialOperationRequest, newBalance decimal.Decimal) (any, error) {
	walletSDK := l.context.GetWalletSDK()
	if walletSDK == nil {
		return nil, gerror.New("钱包SDK未初始化")
	}

	// 获取用户主钱包ID
	mainWalletID, err := l.userLogic.GetMainWalletID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	// 验证钱包是否存在
	summary, err := walletSDK.GetWalletSummary(ctx, mainWalletID)
	if err != nil {
		return nil, gerror.Wrapf(err, "验证钱包存在性失败: WalletID=%s", mainWalletID)
	}
	if summary == nil {
		return nil, gerror.Newf("钱包不存在: WalletID=%s", mainWalletID)
	}

	// 转换金额为原始格式（最小单位）
	rawAmount, err := l.tokenLogic.ConvertBalanceToRaw(ctx, req.Amount, req.TokenSymbol)
	if err != nil {
		return nil, gerror.Wrap(err, "转换金额格式失败")
	}

	// 转换新余额为原始格式
	rawNewBalance, err := l.tokenLogic.ConvertBalanceToRaw(ctx, newBalance, req.TokenSymbol)
	if err != nil {
		return nil, gerror.Wrap(err, "转换新余额格式失败")
	}

	// 根据操作类型调用不同的SDK方法
	switch req.OperationType {
	case OperationTypeCredit:
		// 调用专门处理充值的方法
		return l.executeRemoteWalletCreditOperation(ctx, mainWalletID, req, rawAmount, rawNewBalance)

	case OperationTypeDebit:
		// 调用专门处理扣款的方法
		return l.executeRemoteWalletDebitOperation(ctx, mainWalletID, req, rawAmount, rawNewBalance)

	default:
		return nil, gerror.Newf("不支持的远程操作类型: %s", req.OperationType)
	}
}
*/

// executeRemoteWalletCreditOperation 执行远程钱包充值操作 - 暂时注释，使用本地事务
/*
func (l *operationLogic) executeRemoteWalletCreditOperation(ctx context.Context, walletID string, req *FinancialOperationRequest, rawAmount, rawNewBalance decimal.Decimal) (any, error) {
	// 增加资金 - 调用SDK余额更新方法
	g.Log().Infof(ctx, "执行SDK充值操作: WalletID=%s, Symbol=%s, Amount=%s, NewBalance=%s, BusinessID=%s",
		walletID, req.TokenSymbol, rawAmount.String(), rawNewBalance.String(), req.BusinessID)

	// 调用SDK的实际更新方法 - 这里需要根据SDK文档实现具体的方法调用

	var amount = sdkTypes.Monetary{
		Asset:  req.TokenSymbol,
		Amount: rawAmount.IntPart(), // Corrected: Use rawAmount (actual amount to credit) instead of rawNewBalance
	}
	walletSDK := l.context.GetWalletSDK()
	if walletSDK == nil {
		return nil, gerror.New("钱包SDK未初始化")
	}

	err := walletSDK.CreditWallet(ctx, walletID, amount, nil, nil, req.Metadata, nil, nil)
	if err != nil {
		return nil, gerror.Wrapf(err, "SDK充值操作失败: WalletID=%s, Symbol=%s, Amount=%s",
			walletID, req.TokenSymbol, rawAmount.String())
	}

	// 构造实际的响应格式
	response := map[string]any{
		"wallet_id":        walletID,
		"symbol":           req.TokenSymbol,
		"previous_balance": req.Amount,
		"type":             "credit",
		"new_balance":      rawNewBalance,
		"business_id":      req.BusinessID,
		"status":           "success",
		"timestamp":        gtime.Now().String(),
		"operation":        "balance_update",
		"metadata":         req.Metadata,
	}

	g.Log().Infof(ctx, "远程钱包余额更新完成: WalletID=%s, Symbol=%s, NewBalance=%d",
		walletID, req.TokenSymbol, rawNewBalance)

	return response, nil
}
*/

// executeRemoteWalletDebitOperation 执行远程钱包扣款操作 - 暂时注释，使用本地事务
/*
func (l *operationLogic) executeRemoteWalletDebitOperation(ctx context.Context, walletID string, req *FinancialOperationRequest, rawAmount, rawNewBalance decimal.Decimal) (any, error) {

	walletSDK := l.context.GetWalletSDK()
	if walletSDK == nil {
		return nil, gerror.New("钱包SDK未初始化")
	}

	var amount = sdkTypes.Monetary{
		Asset:  req.TokenSymbol,
		Amount: rawAmount.IntPart(), // Corrected: Use rawAmount (actual amount to debit) instead of rawNewBalance
	}

	_, err := walletSDK.DebitWallet(ctx, walletID, amount, nil, req.Metadata, nil, nil, nil, nil)
	if err != nil {
		return nil, gerror.Wrapf(err, "SDK扣款操作失败: WalletID=%s, Symbol=%s, Amount=%s",
			walletID, req.TokenSymbol, rawAmount.String())
	}

	// 构造实际的响应格式
	response := map[string]any{
		"wallet_id":        walletID,
		"symbol":           req.TokenSymbol,
		"previous_balance": req.Amount,
		"type":             "debit",
		"new_balance":      rawNewBalance,
		"business_id":      req.BusinessID,
		"status":           "success",
		"timestamp":        gtime.Now().String(),
		"operation":        "balance_update",
		"metadata":         req.Metadata,
	}

	g.Log().Infof(ctx, "远程钱包余额更新完成: WalletID=%s, Symbol=%s, NewBalance=%d",
		walletID, req.TokenSymbol, rawNewBalance)

	return response, nil
}
*/

// rollbackRemoteWalletOperation 回滚远程钱包操作 - 暂时注释，使用本地事务
/*
func (l *operationLogic) rollbackRemoteWalletOperation(ctx context.Context, req *FinancialOperationRequest, originalBalance decimal.Decimal, _ any) error {
	walletSDK := l.context.GetWalletSDK()
	if walletSDK == nil {
		return gerror.New("钱包SDK未初始化，无法回滚")
	}

	g.Log().Warningf(ctx, "开始回滚远程钱包操作: UserID=%d, Symbol=%s, OriginalBalance=%s",
		req.UserID, req.TokenSymbol, originalBalance.String())

	// 获取用户主钱包ID
	mainWalletID, err := l.userLogic.GetMainWalletID(ctx, req.UserID)
	if err != nil {
		return gerror.Wrap(err, "获取主钱包ID失败")
	}

	// 构建回滚业务ID
	rollbackBusinessID := req.BusinessID + "_ROLLBACK_" + gtime.Now().Format("YmdHis")

	// 转换原始余额为存储格式
	rawOriginalBalance, err := l.tokenLogic.ConvertBalanceToRaw(ctx, originalBalance, req.TokenSymbol)
	if err != nil {
		return gerror.Wrap(err, "转换原始余额格式失败")
	}

	// 执行回滚操作（恢复到原始余额）
	rollbackResponse, err := l.executeRemoteWalletDebitOperation(ctx, mainWalletID, req, rawOriginalBalance, rawOriginalBalance)
	if err != nil {
		return gerror.Wrap(err, "执行回滚操作失败")
	}

	g.Log().Infof(ctx, "远程钱包操作回滚成功: UserID=%d, Symbol=%s, RestoredBalance=%s, RollbackBusinessID=%s",
		req.UserID, req.TokenSymbol, originalBalance.String(), rollbackBusinessID)

	// 记录回滚操作
	g.Log().Infof(ctx, "回滚响应: %+v", rollbackResponse)

	return nil
}
*/

// createTransactionRecord 创建交易记录
func (l *operationLogic) createTransactionRecord(ctx context.Context, tx gdb.TX, req *FinancialOperationRequest, balanceBefore, balanceAfter decimal.Decimal) (int64, error) {
	// 获取代币信息用于精度转换
	token, err := l.tokenLogic.GetTokenBySymbol(ctx, req.TokenSymbol)
	if err != nil {
		return 0, err
	}

	// 从上下文提取请求信息
	reqCtx := ExtractRequestContext(ctx)

	// 创建交易记录（使用decimal格式）
	transaction := &entity.Transactions{
		UserId:        uint(req.UserID),
		TokenId:       token.TokenId,
		Amount:        req.Amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		Type:          string(req.OperationType),
		Status:        1, // 1表示成功
		Memo:          req.Description,
		Symbol:        req.TokenSymbol,
		Direction:     l.getDirection(req.OperationType),
		WalletType:    "available",    // 默认为可用余额
		BusinessId:    req.BusinessID, // 添加业务ID用于幂等性检查
		CreatedAt:     gtime.Now(),
		UpdatedAt:     gtime.Now(),

		// New fields - User request information
		RequestAmount:    req.Amount, // Using the original request amount
		RequestReference: req.BusinessID,
		RequestMetadata:  l.mergeMetadataToJSON(reqCtx.Metadata, req.Metadata),
		RequestSource:    reqCtx.Source,
		RequestIp:        reqCtx.IP,
		RequestUserAgent: reqCtx.UserAgent,
		RequestTimestamp: gtime.Now(),
		ProcessedAt:      gtime.Now(),

		// Use fee from request parameters (trusted from upstream)
		FeeAmount: req.FeeAmount,
		FeeType:   req.FeeType,

		// Exchange rate (set to 1 if no conversion)
		ExchangeRate: decimal.NewFromInt(1),

		// Target user fields (populated from request metadata for transfers)
		TargetUserId:   l.extractTargetUserId(req.Metadata),
		TargetUsername: l.extractTargetUsername(req.Metadata),
	}

	transactionID, err := l.context.GetTransactionDAO().CreateTransaction(ctx, tx, transaction)
	if err != nil {
		return 0, err
	}

	return transactionID, nil
}

// ensureWalletExists 确保用户的远程钱包和本地钱包记录都存在
func (l *operationLogic) ensureWalletExists(ctx context.Context, user *entity.Users, tokenSymbol string) error {
	g.Log().Infof(ctx, "开始检测用户钱包: UserID=%d, Symbol=%s", user.Id, tokenSymbol)

	// 1. 检查用户是否有主钱包ID - 暂时跳过远程钱包创建
	if user.MainWalletId == "" {
		g.Log().Infof(ctx, "用户没有主钱包ID，暂时使用本地钱包: UserID=%d", user.Id)

		// 暂时生成一个本地钱包ID
		mainWalletID := fmt.Sprintf("LOCAL_WALLET_%d", user.Id)

		// 更新用户记录中的主钱包ID
		err := l.updateUserMainWalletID(ctx, uint64(user.Id), mainWalletID)
		if err != nil {
			return gerror.Wrap(err, "更新用户主钱包ID失败")
		}

		user.MainWalletId = mainWalletID
		g.Log().Infof(ctx, "成功创建本地钱包ID: UserID=%d, WalletID=%s", user.Id, mainWalletID)
	}

	// 2. 跳过远程钱包验证
	// err := l.verifyRemoteWalletExists(ctx, user.MainWalletId)
	// if err != nil {
	// 	return gerror.Wrapf(err, "远程钱包验证失败: WalletID=%s", user.MainWalletId)
	// }

	// 3. 检查本地钱包记录是否存在
	err := l.ensureLocalWalletRecord(ctx, uint64(user.Id), tokenSymbol)
	if err != nil {
		return gerror.Wrap(err, "确保本地钱包记录失败")
	}

	g.Log().Infof(ctx, "钱包检测完成: UserID=%d, Symbol=%s", user.Id, tokenSymbol)
	return nil
}

// createRemoteWallet 创建远程钱包 - 暂时注释，使用本地事务
/*
func (l *operationLogic) createRemoteWallet(ctx context.Context, user *entity.Users) (string, error) {
	walletSDK := l.context.GetWalletSDK()
	if walletSDK == nil {
		return "", gerror.New("钱包SDK未初始化")
	}

	// 构造钱包名称和元数据
	walletName := fmt.Sprintf("%d_main_wallet", user.TelegramId)
	metadata := map[string]string{
		"telegram_id":  fmt.Sprintf("%d", user.TelegramId),
		"user_id":      fmt.Sprintf("%d", user.Id),
		"wallet_usage": "primary",
		"created_by":   "financial_operation_validation",
	}

	g.Log().Infof(ctx, "创建远程钱包: UserID=%d, WalletName=%s", user.Id, walletName)

	// 调用SDK创建钱包
	response, err := walletSDK.CreateWallet(ctx, walletName, metadata)
	if err != nil {
		return "", gerror.Wrapf(err, "调用SDK创建钱包失败: UserID=%d", user.Id)
	}

	if response == nil || response.Id == "" {
		return "", gerror.Newf("SDK返回空钱包ID: UserID=%d", user.Id)
	}

	g.Log().Infof(ctx, "远程钱包创建成功: UserID=%d, WalletID=%s", user.Id, response.Id)
	return response.Id, nil
}
*/

// updateUserMainWalletID 更新用户的主钱包ID
func (l *operationLogic) updateUserMainWalletID(ctx context.Context, userID uint64, mainWalletID string) error {
	// 直接使用 GoFrame 的 DAO 来更新主钱包ID
	// 参考 create_user_wallet_v2.go 中的实现方式
	_, err := g.Model("users").Ctx(ctx).
		Where("id = ?", userID).
		Update(g.Map{
			"main_wallet_id": mainWalletID,
			"updated_at":     gtime.Now(),
		})
	if err != nil {
		return gerror.Wrapf(err, "更新用户主钱包ID失败: UserID=%d, WalletID=%s", userID, mainWalletID)
	}
	return nil
}

// verifyRemoteWalletExists 验证远程钱包是否存在 - 暂时跳过
/*
func (l *operationLogic) verifyRemoteWalletExists(ctx context.Context, walletID string) error {
	walletSDK := l.context.GetWalletSDK()
	if walletSDK == nil {
		return gerror.New("钱包SDK未初始化")
	}

	g.Log().Debugf(ctx, "验证远程钱包存在性: WalletID=%s", walletID)

	// 尝试获取钱包摘要来验证钱包是否存在
	summary, err := walletSDK.GetWalletSummary(ctx, walletID)
	if err != nil {
		return gerror.Wrapf(err, "获取远程钱包摘要失败: WalletID=%s", walletID)
	}

	if summary == nil {
		return gerror.Newf("远程钱包不存在: WalletID=%s", walletID)
	}

	g.Log().Debugf(ctx, "远程钱包验证成功: WalletID=%s", walletID)
	return nil
}
*/

// ensureLocalWalletRecord 确保本地钱包记录存在
func (l *operationLogic) ensureLocalWalletRecord(ctx context.Context, userID uint64, tokenSymbol string) error {
	// 检查本地钱包记录是否存在
	_, _, err := l.balanceLogic.GetLocalBalance(ctx, userID, tokenSymbol)
	if err == nil {
		// 本地记录已存在
		g.Log().Debugf(ctx, "本地钱包记录已存在: UserID=%d, Symbol=%s", userID, tokenSymbol)
		return nil
	}

	g.Log().Infof(ctx, "本地钱包记录不存在，需要创建: UserID=%d, Symbol=%s", userID, tokenSymbol)

	// 获取代币信息
	token, err := l.tokenLogic.GetTokenBySymbol(ctx, tokenSymbol)
	if err != nil {
		return gerror.Wrapf(err, "获取代币信息失败: Symbol=%s", tokenSymbol)
	}

	// 创建本地钱包记录（初始余额为0）
	walletRecord := &entity.Wallets{
		UserId:           uint(userID),
		TokenId:          uint(token.TokenId),
		Symbol:           tokenSymbol,
		AvailableBalance: 0,
		FrozenBalance:    0,
		DecimalPlaces:    uint(token.Decimals),
		CreatedAt:        gtime.Now(),
		UpdatedAt:        gtime.Now(),
	}

	err = l.context.GetWalletDAO().CreateWalletRecord(ctx, nil, walletRecord)
	if err != nil {
		return gerror.Wrapf(err, "创建本地钱包记录失败: UserID=%d, TokenID=%d", userID, token.TokenId)
	}

	g.Log().Infof(ctx, "本地钱包记录创建成功: UserID=%d, Symbol=%s", userID, tokenSymbol)
	return nil
}

// validateBalanceConsistency 校验本地和远程余额一致性
func (l *operationLogic) validateBalanceConsistency(ctx context.Context, userID uint64, tokenSymbol string) error {
	g.Log().Debugf(ctx, "开始校验余额一致性: UserID=%d, Symbol=%s", userID, tokenSymbol)

	// 获取本地余额
	localBalance, _, err := l.balanceLogic.GetLocalBalance(ctx, userID, tokenSymbol)
	if err != nil {
		// 如果本地记录不存在，说明是新创建的，无需校验
		g.Log().Infof(ctx, "本地钱包记录不存在，跳过余额一致性校验: UserID=%d, Symbol=%s", userID, tokenSymbol)
		return nil
	}

	// 获取远程余额
	remoteBalance, err := l.balanceLogic.GetRemoteBalance(ctx, userID, tokenSymbol)
	if err != nil {
		g.Log().Warningf(ctx, "获取远程余额失败，跳过一致性校验: UserID=%d, Symbol=%s, Error=%v", userID, tokenSymbol, err)
		return nil // 不阻止操作，只记录警告
	}

	// 比较余额
	if !localBalance.Equal(remoteBalance) {
		g.Log().Warningf(ctx, "余额不一致: UserID=%d, Symbol=%s, Local=%s, Remote=%s",
			userID, tokenSymbol, localBalance.String(), remoteBalance.String())

		// 自动同步余额（使用远程余额为准）
		err = l.balanceLogic.UpdateLocalBalance(ctx, nil, userID, tokenSymbol, remoteBalance, decimal.Zero)
		if err != nil {
			g.Log().Errorf(ctx, "自动同步余额失败: UserID=%d, Symbol=%s, Error=%v", userID, tokenSymbol, err)
			return gerror.Wrap(err, "余额不一致且同步失败")
		}

		g.Log().Infof(ctx, "余额已自动同步: UserID=%d, Symbol=%s, 从 %s 同步到 %s",
			userID, tokenSymbol, localBalance.String(), remoteBalance.String())
	} else {
		g.Log().Debugf(ctx, "余额一致性校验通过: UserID=%d, Symbol=%s, Balance=%s",
			userID, tokenSymbol, localBalance.String())
	}

	return nil
}

// getDirection 根据操作类型获取资金方向
func (l *operationLogic) getDirection(operationType OperationType) string {
	switch operationType {
	case OperationTypeCredit:
		return "in" // 增加
	case OperationTypeDebit:
		return "out" // 减少
	default:
		return "in"
	}
}

// extractTargetUserId 从元数据中提取目标用户ID
func (l *operationLogic) extractTargetUserId(metadata map[string]string) uint {
	if metadata == nil {
		return 0
	}

	if targetUserIdStr, exists := metadata["target_user_id"]; exists {
		var targetUserId uint64
		if _, err := fmt.Sscanf(targetUserIdStr, "%d", &targetUserId); err == nil {
			return uint(targetUserId)
		}
	}

	return 0
}

// extractTargetUsername 从元数据中提取目标用户名
func (l *operationLogic) extractTargetUsername(metadata map[string]string) string {
	if metadata == nil {
		return ""
	}

	if targetUsername, exists := metadata["target_username"]; exists {
		return targetUsername
	}

	return ""
}

// mergeMetadataToJSON 合并请求上下文元数据和业务元数据并转换为JSON
func (l *operationLogic) mergeMetadataToJSON(contextMetadata, businessMetadata map[string]string) string {
	merged := make(map[string]string)

	// 先添加上下文元数据
	for k, v := range contextMetadata {
		merged[k] = v
	}

	// 再添加业务元数据（可能覆盖上下文元数据）
	for k, v := range businessMetadata {
		merged[k] = v
	}

	if len(merged) == 0 {
		return "{}"
	}

	data, err := json.Marshal(merged)
	if err != nil {
		g.Log().Errorf(context.Background(), "合并元数据到JSON失败: %v", err)
		return "{}"
	}

	return string(data)
}
