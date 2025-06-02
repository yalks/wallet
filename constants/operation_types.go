package constants

// OperationType represents the type of wallet operation
type OperationType string

const (
	// Basic operation types
	OperationTypeCredit   OperationType = "credit"   // 增加余额
	OperationTypeDebit    OperationType = "debit"    // 减少余额
	OperationTypeTransfer OperationType = "transfer" // 转账操作
	OperationTypeFreeze   OperationType = "freeze"   // 冻结资金
	OperationTypeUnfreeze OperationType = "unfreeze" // 解冻资金
)

// WalletType represents different wallet types
type WalletType string

const (
	WalletTypeAvailable WalletType = "available" // 可用余额钱包
	WalletTypeFrozen    WalletType = "frozen"    // 冻结余额钱包
	WalletTypePending   WalletType = "pending"   // 待处理余额钱包
)

// TransactionStatus represents the status of a transaction
type TransactionStatus string

const (
	TransactionStatusPending    TransactionStatus = "pending"    // 待处理
	TransactionStatusProcessing TransactionStatus = "processing" // 处理中
	TransactionStatusCompleted  TransactionStatus = "completed"  // 已完成
	TransactionStatusFailed     TransactionStatus = "failed"     // 失败
	TransactionStatusCancelled  TransactionStatus = "cancelled"  // 已取消
	TransactionStatusExpired    TransactionStatus = "expired"    // 已过期
	TransactionStatusRefunded   TransactionStatus = "refunded"   // 已退款
)

// OperationContext contains metadata for an operation
type OperationContext struct {
	OperationType OperationType
	FundType      FundType
	FromWallet    WalletType
	ToWallet      WalletType
	Status        TransactionStatus
}

// GetOperationTypeForFundType determines the operation type based on fund type
func GetOperationTypeForFundType(fundType FundType) OperationType {
	// Check for special cases first
	switch fundType {
	case FundTypeTransferOut, FundTypeTransferIn:
		return OperationTypeTransfer
	default:
		// Use direction for other fund types
		direction := GetFundDirection(fundType)
		switch direction {
		case FundDirectionIn:
			return OperationTypeCredit
		case FundDirectionOut:
			return OperationTypeDebit
		default:
			return OperationTypeCredit
		}
	}
}

// IsValidOperationType checks if an operation type is valid
func IsValidOperationType(opType OperationType) bool {
	switch opType {
	case OperationTypeCredit, OperationTypeDebit, OperationTypeTransfer, 
		 OperationTypeFreeze, OperationTypeUnfreeze:
		return true
	default:
		return false
	}
}

// IsValidWalletType checks if a wallet type is valid
func IsValidWalletType(walletType WalletType) bool {
	switch walletType {
	case WalletTypeAvailable, WalletTypeFrozen, WalletTypePending:
		return true
	default:
		return false
	}
}

// IsValidTransactionStatus checks if a transaction status is valid
func IsValidTransactionStatus(status TransactionStatus) bool {
	switch status {
	case TransactionStatusPending, TransactionStatusProcessing, TransactionStatusCompleted,
		 TransactionStatusFailed, TransactionStatusCancelled, TransactionStatusExpired,
		 TransactionStatusRefunded:
		return true
	default:
		return false
	}
}

// IsFinalStatus checks if a transaction status is final (cannot be changed)
func IsFinalStatus(status TransactionStatus) bool {
	switch status {
	case TransactionStatusCompleted, TransactionStatusFailed, 
		 TransactionStatusCancelled, TransactionStatusRefunded:
		return true
	default:
		return false
	}
}