package wallet

import (
	"context"

	"github.com/yalks/wallet/constants"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/shopspring/decimal"
)

// FundOperationRequest 资金操作请求
type FundOperationRequest struct {
	UserID      uint64             `json:"user_id" validate:"required"`      // 用户ID
	TokenSymbol string             `json:"token_symbol" validate:"required"` // 代币符号
	Amount      decimal.Decimal    `json:"amount" validate:"required,gt=0"`  // 操作金额（用户友好格式）
	BusinessID  string             `json:"business_id" validate:"required"`  // 业务ID（用于幂等性）
	FundType    constants.FundType `json:"fund_type" validate:"required"`    // 资金类型
	Description string             `json:"description"`                      // 操作描述
	Metadata    map[string]string  `json:"metadata,omitempty"`               // 扩展元数据
	RelatedID   int64              `json:"related_id,omitempty"`             // 关联ID（如红包ID、转账ID等）
}

// TransferOperationRequest 转账操作请求（双方交易）
type TransferOperationRequest struct {
	FromUserID  uint64             `json:"from_user_id" validate:"required"` // 发送方用户ID
	ToUserID    uint64             `json:"to_user_id" validate:"required"`   // 接收方用户ID
	TokenSymbol string             `json:"token_symbol" validate:"required"` // 代币符号
	Amount      decimal.Decimal    `json:"amount" validate:"required,gt=0"`  // 转账金额
	BusinessID  string             `json:"business_id" validate:"required"`  // 业务ID（用于幂等性）
	FundType    constants.FundType `json:"fund_type" validate:"required"`    // 资金类型
	Description string             `json:"description"`                      // 操作描述
	Metadata    map[string]string  `json:"metadata,omitempty"`               // 扩展元数据
	RelatedID   int64              `json:"related_id,omitempty"`             // 关联ID
}

// FundOperationResult 资金操作结果
type FundOperationResult struct {
	TransactionID string          `json:"transaction_id"` // 交易ID
	BalanceBefore decimal.Decimal `json:"balance_before"` // 操作前余额
	BalanceAfter  decimal.Decimal `json:"balance_after"`  // 操作后余额
	RawAmount     decimal.Decimal `json:"raw_amount"`     // 原始金额（发送给远程钱包的格式）
	Successful    bool            `json:"successful"`     // 是否成功
}

// BalanceInfo 余额信息
type BalanceInfo struct {
	UserID           uint64          `json:"user_id"`           // 用户ID
	TokenSymbol      string          `json:"token_symbol"`      // 代币符号
	AvailableBalance decimal.Decimal `json:"available_balance"` // 可用余额
	FrozenBalance    decimal.Decimal `json:"frozen_balance"`    // 冻结余额
	TotalBalance     decimal.Decimal `json:"total_balance"`     // 总余额
	LastSyncTime     string          `json:"last_sync_time"`    // 最后同步时间
}

// IWalletManager 统一的钱包管理接口
type IWalletManager interface {
	GetBalance(ctx context.Context, userID uint64, tokenSymbol string) (*BalanceInfo, error)

	// 事务中的操作（已废弃，请使用 ProcessFundOperationInTx）
	// Deprecated: Use ProcessFundOperationInTx instead
	CreditFundsInTx(ctx context.Context, tx gdb.TX, req *FundOperationRequest) (*FundOperationResult, error)
	// Deprecated: Use ProcessFundOperationInTx instead
	DebitFundsInTx(ctx context.Context, tx gdb.TX, req *FundOperationRequest) (*FundOperationResult, error)

	// 推荐使用：基于资金类型的通用操作方法
	ProcessFundOperationInTx(ctx context.Context, tx gdb.TX, req *FundOperationRequest) (*FundOperationResult, error)

	// 推荐使用：转账操作（双方交易的原子操作）
	ProcessTransferInTx(ctx context.Context, tx gdb.TX, req *TransferOperationRequest) (*TransferOperationResult, error)

	// 推荐使用：使用 TransactionBuilder 创建交易
	CreateTransactionWithBuilder(ctx context.Context, txReq *constants.TransactionRequest) (int64, error)
}

// TransferOperationResult 转账操作结果
type TransferOperationResult struct {
	FromTransactionID string          `json:"from_transaction_id"` // 发送方交易ID
	ToTransactionID   string          `json:"to_transaction_id"`   // 接收方交易ID
	FromBalanceAfter  decimal.Decimal `json:"from_balance_after"`  // 发送方操作后余额
	ToBalanceAfter    decimal.Decimal `json:"to_balance_after"`    // 接收方操作后余额
}

// ITransactionManager 事务管理接口
type ITransactionManager interface {
	// 创建交易记录
	CreateTransaction(ctx context.Context, req *constants.TransactionRequest) (int64, error)

	// 获取交易记录
	GetTransactionByID(ctx context.Context, transactionID int64) (*TransactionRecord, error)
	GetTransactionByReference(ctx context.Context, reference string) (*TransactionRecord, error)

	// 更新交易状态
	UpdateTransactionStatus(ctx context.Context, transactionID int64, status constants.TransactionStatus) error

	// 查询交易历史
	GetUserTransactionHistory(ctx context.Context, userID int64, fundType constants.FundType, limit, offset int) ([]*TransactionRecord, error)
}

// TransactionRecord 交易记录
type TransactionRecord struct {
	ID          int64                       `json:"id"`
	UserID      int64                       `json:"user_id"`
	WalletID    int64                       `json:"wallet_id"`
	TokenID     int64                       `json:"token_id"`
	Amount      string                      `json:"amount"`
	FundType    constants.FundType          `json:"fund_type"`
	Direction   constants.FundDirection     `json:"direction"`
	Status      constants.TransactionStatus `json:"status"`
	Reference   string                      `json:"reference"`
	Description string                      `json:"description"`
	RelatedID   int64                       `json:"related_id"`
	Metadata    map[string]interface{}      `json:"metadata"`
	CreatedAt   string                      `json:"created_at"`
	UpdatedAt   string                      `json:"updated_at"`
}
