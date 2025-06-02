package wallet

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shopspring/decimal"

	"github.com/yalks/wallet/constants"
	"github.com/yalks/wallet/logic"
)

// walletManager 统一钱包管理器实现
type walletManager struct {
	// 独立的业务逻辑层
	userLogic      logic.IUserLogic
	tokenLogic     logic.ITokenLogic
	balanceLogic   logic.IBalanceLogic
	operationLogic logic.IOperationLogic

	// 事务管理器
	transactionManager ITransactionManager
}

// initialize 初始化钱包管理器的各个组件
func (m *walletManager) initialize(ctx context.Context) error {
	g.Log().Info(ctx, "初始化钱包管理器组件...")

	// 初始化各个逻辑组件
	m.userLogic = logic.NewUserLogic()
	m.tokenLogic = logic.NewTokenLogic()
	m.balanceLogic = logic.NewBalanceLogic()
	m.operationLogic = logic.NewOperationLogic()

	// 初始化事务管理器
	m.transactionManager = NewTransactionManager()

	// 逻辑组件不需要额外的初始化，它们在创建时会自动初始化

	g.Log().Info(ctx, "钱包管理器组件初始化完成")
	return nil
}

// CreditFundsInTx 在事务中增加资金
func (m *walletManager) CreditFundsInTx(ctx context.Context, tx gdb.TX, req *FundOperationRequest) (*FundOperationResult, error) {
	// 参数验证
	if err := m.validateRequest(req); err != nil {
		return nil, err
	}

	// 构建财务操作请求
	financialReq := &logic.FinancialOperationRequest{
		UserID:        req.UserID,
		TokenSymbol:   req.TokenSymbol,
		Amount:        req.Amount,
		OperationType: logic.OperationTypeCredit, // 增加资金
		BusinessID:    req.BusinessID,
		Description:   req.Description,
		Metadata:      req.Metadata,
	}

	// 执行财务操作
	opResult, err := m.operationLogic.ExecuteInTx(ctx, tx, financialReq)
	if err != nil {
		return nil, gerror.Wrap(err, "增加资金操作失败")
	}

	// 转换金额格式为原始格式用于返回
	rawAmount, err := m.tokenLogic.ConvertBalanceToRaw(ctx, req.Amount, req.TokenSymbol)
	if err != nil {
		return nil, gerror.Wrapf(err, "转换金额到原始格式失败: Amount=%s, Symbol=%s", req.Amount.String(), req.TokenSymbol)
	}

	result := &FundOperationResult{
		TransactionID: fmt.Sprintf("%d", opResult.TransactionID),
		BalanceBefore: opResult.BalanceBefore,
		BalanceAfter:  opResult.BalanceAfter,
		RawAmount:     rawAmount,
		Successful:    true,
	}

	g.Log().Infof(ctx, "资金增加成功: UserID=%d, Symbol=%s, Amount=%s, TransactionID=%s",
		req.UserID, req.TokenSymbol, req.Amount.String(), result.TransactionID)

	return result, nil
}

// DebitFundsInTx 在事务中减少资金
func (m *walletManager) DebitFundsInTx(ctx context.Context, tx gdb.TX, req *FundOperationRequest) (*FundOperationResult, error) {
	// 参数验证
	if err := m.validateRequest(req); err != nil {
		return nil, err
	}

	// 余额验证
	balance, err := m.operationLogic.GetUserBalance(ctx, req.UserID, req.TokenSymbol)
	if err != nil {
		return nil, gerror.Wrap(err, "获取用户余额失败")
	}
	if balance.LessThan(req.Amount) {
		return nil, gerror.Newf("余额不足: 当前余额=%s, 需要金额=%s", balance.String(), req.Amount.String())
	}

	// 构建财务操作请求
	financialReq := &logic.FinancialOperationRequest{
		UserID:        req.UserID,
		TokenSymbol:   req.TokenSymbol,
		Amount:        req.Amount,
		OperationType: logic.OperationTypeDebit, // 减少资金
		BusinessID:    req.BusinessID,
		Description:   req.Description,
		Metadata:      req.Metadata,
	}

	// 执行财务操作
	opResult, err := m.operationLogic.ExecuteInTx(ctx, tx, financialReq)
	if err != nil {
		return nil, gerror.Wrap(err, "减少资金操作失败")
	}

	// 转换金额格式为原始格式用于返回
	rawAmount, err := m.tokenLogic.ConvertBalanceToRaw(ctx, req.Amount, req.TokenSymbol)
	if err != nil {
		return nil, gerror.Wrapf(err, "转换金额到原始格式失败: Amount=%s, Symbol=%s", req.Amount.String(), req.TokenSymbol)
	}

	result := &FundOperationResult{
		TransactionID: fmt.Sprintf("%d", opResult.TransactionID),
		BalanceBefore: opResult.BalanceBefore,
		BalanceAfter:  opResult.BalanceAfter,
		RawAmount:     rawAmount,
		Successful:    true,
	}

	g.Log().Infof(ctx, "资金减少成功: UserID=%d, Symbol=%s, Amount=%s, TransactionID=%s",
		req.UserID, req.TokenSymbol, req.Amount.String(), result.TransactionID)

	return result, nil
}

// ProcessFundOperationInTx 基于资金类型的通用操作方法
func (m *walletManager) ProcessFundOperationInTx(ctx context.Context, tx gdb.TX, req *FundOperationRequest) (*FundOperationResult, error) {
	// 验证资金类型
	if !constants.IsValidFundType(req.FundType) {
		return nil, gerror.Newf("无效的资金类型: %s", req.FundType)
	}

	// 根据资金类型的方向决定操作
	direction := constants.GetFundDirection(req.FundType)

	// 如果描述为空，使用资金类型的默认描述
	if req.Description == "" {
		if info, exists := constants.GetFundTypeInfo(req.FundType); exists {
			req.Description = info.Description
		}
	}

	// 根据方向执行相应操作
	switch direction {
	case constants.FundDirectionIn:
		return m.CreditFundsInTx(ctx, tx, req)
	case constants.FundDirectionOut:
		return m.DebitFundsInTx(ctx, tx, req)
	default:
		return nil, gerror.Newf("资金类型 %s 的方向未定义", req.FundType)
	}
}

// GetBalance 获取余额
func (m *walletManager) GetBalance(ctx context.Context, userID uint64, tokenSymbol string) (*BalanceInfo, error) {
	// 获取余额信息
	availableBalance, frozenBalance, err := m.balanceLogic.GetBalance(ctx, userID, tokenSymbol)
	if err != nil {
		return nil, gerror.Wrapf(err, "获取余额失败: UserID=%d, Symbol=%s", userID, tokenSymbol)
	}

	totalBalance := availableBalance.Add(frozenBalance)

	return &BalanceInfo{
		UserID:           userID,
		TokenSymbol:      tokenSymbol,
		AvailableBalance: availableBalance,
		FrozenBalance:    frozenBalance,
		TotalBalance:     totalBalance,
		LastSyncTime:     time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

// validateRequest 验证请求参数并进行安全检查
func (m *walletManager) validateRequest(req *FundOperationRequest) error {
	// 基础参数验证
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

	// 安全检查：金额上限验证
	maxAmount := decimal.NewFromFloat(1000000) // 设置最大金额为100万
	if req.Amount.GreaterThan(maxAmount) {
		return gerror.Newf("单次操作金额不能超过 %s", maxAmount.String())
	}

	// 安全检查：业务ID格式验证（防止SQL注入）
	if len(req.BusinessID) > 255 {
		return gerror.New("业务ID长度不能超过255个字符")
	}

	// 检查业务ID是否包含非法字符
	for _, char := range req.BusinessID {
		if char < 32 || char > 126 {
			return gerror.New("业务ID包含非法字符")
		}
	}

	// 安全检查：代币符号格式验证
	if len(req.TokenSymbol) > 20 {
		return gerror.New("代币符号长度不能超过20个字符")
	}

	// 检查精度是否合理（防止精度攻击）
	if req.Amount.Exponent() < -18 {
		return gerror.New("金额精度过高，最多支持18位小数")
	}

	return nil
}

// ProcessTransferInTx 转账操作（双方交易的原子操作）
func (m *walletManager) ProcessTransferInTx(ctx context.Context, tx gdb.TX, req *TransferOperationRequest) (*TransferOperationResult, error) {
	// 参数验证
	if err := m.validateTransferRequest(req); err != nil {
		return nil, err
	}

	// 添加目标用户信息到元数据
	debitMetadata := make(map[string]string)
	for k, v := range req.Metadata {
		debitMetadata[k] = v
	}
	debitMetadata["target_user_id"] = fmt.Sprintf("%d", req.ToUserID)
	if req.ToUsername != "" {
		debitMetadata["target_username"] = req.ToUsername
	}

	// 构建发送方扣款请求
	debitReq := &FundOperationRequest{
		UserID:           req.FromUserID,
		TokenSymbol:      req.TokenSymbol,
		Amount:           req.Amount,
		BusinessID:       req.BusinessID + "_debit",
		FundType:         req.FundType,
		Description:      fmt.Sprintf("转账给用户%d: %s", req.ToUserID, req.Description),
		Metadata:         debitMetadata,
		RelatedID:        req.RelatedID,
		RequestSource:    req.RequestSource,
		RequestIP:        req.RequestIP,
		RequestUserAgent: req.RequestUserAgent,
	}

	// 构建接收方加款请求
	creditReq := &FundOperationRequest{
		UserID:           req.ToUserID,
		TokenSymbol:      req.TokenSymbol,
		Amount:           req.Amount,
		BusinessID:       req.BusinessID + "_credit",
		FundType:         req.FundType,
		Description:      fmt.Sprintf("收到用户%d转账: %s", req.FromUserID, req.Description),
		Metadata:         req.Metadata,
		RelatedID:        req.RelatedID,
		RequestSource:    req.RequestSource,
		RequestIP:        req.RequestIP,
		RequestUserAgent: req.RequestUserAgent,
	}

	// 执行发送方扣款
	debitResult, err := m.DebitFundsInTx(ctx, tx, debitReq)
	if err != nil {
		return nil, gerror.Wrap(err, "转账扣款失败")
	}

	// 执行接收方加款
	creditResult, err := m.CreditFundsInTx(ctx, tx, creditReq)
	if err != nil {
		return nil, gerror.Wrap(err, "转账加款失败")
	}

	return &TransferOperationResult{
		FromTransactionID: debitResult.TransactionID,
		ToTransactionID:   creditResult.TransactionID,
		FromBalanceAfter:  debitResult.BalanceAfter,
		ToBalanceAfter:    creditResult.BalanceAfter,
	}, nil
}

// CreateTransactionWithBuilder 使用 TransactionBuilder 创建交易
func (m *walletManager) CreateTransactionWithBuilder(ctx context.Context, txReq *constants.TransactionRequest) (int64, error) {
	return m.transactionManager.CreateTransaction(ctx, txReq)
}

// validateTransferRequest 验证转账请求参数
func (m *walletManager) validateTransferRequest(req *TransferOperationRequest) error {
	if req.FromUserID == 0 {
		return gerror.New("发送方用户ID不能为空")
	}
	if req.ToUserID == 0 {
		return gerror.New("接收方用户ID不能为空")
	}
	if req.FromUserID == req.ToUserID {
		return gerror.New("不能向自己转账")
	}
	if req.TokenSymbol == "" {
		return gerror.New("代币符号不能为空")
	}
	if req.Amount.LessThanOrEqual(decimal.Zero) {
		return gerror.New("转账金额必须大于0")
	}
	if req.BusinessID == "" {
		return gerror.New("业务ID不能为空")
	}
	if !constants.IsValidFundType(req.FundType) {
		return gerror.Newf("无效的资金类型: %s", req.FundType)
	}

	return nil
}
