package constants

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

// TransactionRequest represents a request to create a transaction
type TransactionRequest struct {
	UserID       int64
	WalletID     int64
	Amount       string
	TokenID      int64
	FundType     FundType
	Reference    string
	Description  string
	RelatedID    int64  // Related entity ID (e.g., red packet ID, transfer ID)
	Metadata     map[string]interface{}
}

// TransactionBuilder helps build transaction requests with proper validation
type TransactionBuilder struct {
	request TransactionRequest
	mu      sync.Mutex // 确保线程安全
}

// NewTransactionBuilder creates a new transaction builder
func NewTransactionBuilder() *TransactionBuilder {
	return &TransactionBuilder{
		request: TransactionRequest{
			Metadata: make(map[string]interface{}),
		},
	}
}

// WithUser sets the user ID
func (b *TransactionBuilder) WithUser(userID int64) *TransactionBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.request.UserID = userID
	return b
}

// WithWallet sets the wallet ID
func (b *TransactionBuilder) WithWallet(walletID int64) *TransactionBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.request.WalletID = walletID
	return b
}

// WithAmount sets the amount
func (b *TransactionBuilder) WithAmount(amount string) *TransactionBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.request.Amount = strings.TrimSpace(amount)
	return b
}

// WithToken sets the token ID
func (b *TransactionBuilder) WithToken(tokenID int64) *TransactionBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.request.TokenID = tokenID
	return b
}

// WithFundType sets the fund type
func (b *TransactionBuilder) WithFundType(fundType FundType) *TransactionBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.request.FundType = fundType
	return b
}

// WithReference sets the reference
func (b *TransactionBuilder) WithReference(reference string) *TransactionBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.request.Reference = strings.TrimSpace(reference)
	return b
}

// WithDescription sets the description
func (b *TransactionBuilder) WithDescription(description string) *TransactionBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.request.Description = strings.TrimSpace(description)
	return b
}

// WithRelatedID sets the related entity ID
func (b *TransactionBuilder) WithRelatedID(id int64) *TransactionBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.request.RelatedID = id
	return b
}

// WithMetadata adds metadata
func (b *TransactionBuilder) WithMetadata(key string, value interface{}) *TransactionBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	if key != "" {
		b.request.Metadata[key] = value
	}
	return b
}

// Build validates and returns the transaction request
func (b *TransactionBuilder) Build() (*TransactionRequest, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	// Validate required fields
	if b.request.UserID <= 0 {
		return nil, fmt.Errorf("user ID must be positive")
	}
	if b.request.WalletID <= 0 {
		return nil, fmt.Errorf("wallet ID must be positive")
	}
	if b.request.Amount == "" {
		return nil, fmt.Errorf("amount is required")
	}
	
	// 验证金额格式
	amount := strings.TrimSpace(b.request.Amount)
	if err := b.validateAmount(amount); err != nil {
		return nil, fmt.Errorf("invalid amount: %v", err)
	}
	b.request.Amount = amount
	
	if b.request.TokenID <= 0 {
		return nil, fmt.Errorf("token ID must be positive")
	}
	if !IsValidFundType(b.request.FundType) {
		return nil, fmt.Errorf("invalid fund type: %s", b.request.FundType)
	}

	// Auto-generate reference if not provided
	if b.request.Reference == "" {
		b.request.Reference = b.generateReference()
	}

	// Auto-generate description if not provided
	if b.request.Description == "" {
		if info, exists := GetFundTypeInfo(b.request.FundType); exists {
			b.request.Description = info.Description
		}
	}
	
	// 创建副本返回，避免后续修改影响已构建的请求
	result := b.request
	result.Metadata = make(map[string]interface{})
	for k, v := range b.request.Metadata {
		result.Metadata[k] = v
	}

	return &result, nil
}

// validateAmount 验证金额格式和范围
func (b *TransactionBuilder) validateAmount(amount string) error {
	// 检查基本格式
	if matched, _ := regexp.MatchString(`^[0-9]+(\.[0-9]+)?$`, amount); !matched {
		return fmt.Errorf("must be a positive number")
	}
	
	// 使用 decimal 进行精确验证
	d, err := decimal.NewFromString(amount)
	if err != nil {
		return err
	}
	
	// 检查是否为零或负数
	if d.LessThanOrEqual(decimal.Zero) {
		return fmt.Errorf("must be greater than zero")
	}
	
	// 检查最大值（100万）
	maxAmount, _ := decimal.NewFromString("1000000")
	if d.GreaterThan(maxAmount) {
		return fmt.Errorf("exceeds maximum amount (1,000,000)")
	}
	
	// 检查最小值
	minAmount, _ := decimal.NewFromString("0.000000000000000001")
	if d.LessThan(minAmount) {
		return fmt.Errorf("below minimum precision")
	}
	
	// 检查小数位数（最多18位）
	parts := strings.Split(amount, ".")
	if len(parts) == 2 && len(parts[1]) > 18 {
		return fmt.Errorf("too many decimal places (max 18)")
	}
	
	return nil
}

// generateReference generates a unique reference for a transaction
func (b *TransactionBuilder) generateReference() string {
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%s_%d_%d", b.request.FundType, b.request.UserID, timestamp)
}

// Clone 克隆构建器，用于创建相似的交易
func (b *TransactionBuilder) Clone() *TransactionBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	newBuilder := NewTransactionBuilder()
	newBuilder.request = b.request
	
	// 深拷贝 metadata
	newBuilder.request.Metadata = make(map[string]interface{})
	for k, v := range b.request.Metadata {
		newBuilder.request.Metadata[k] = v
	}
	
	return newBuilder
}

// Reset 重置构建器，允许重用
func (b *TransactionBuilder) Reset() *TransactionBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	b.request = TransactionRequest{
		Metadata: make(map[string]interface{}),
	}
	return b
}


// Common transaction builders for specific fund types

// BuildRedPacketCreateTransaction builds a red packet creation transaction
func BuildRedPacketCreateTransaction(userID, walletID int64, amount string, tokenID int64, redPacketID int64) (*TransactionRequest, error) {
	return NewTransactionBuilder().
		WithUser(userID).
		WithWallet(walletID).
		WithAmount(amount).
		WithToken(tokenID).
		WithFundType(FundTypeRedPacketCreate).
		WithRelatedID(redPacketID).
		WithMetadata("red_packet_id", redPacketID).
		Build()
}

// BuildTransferOutTransaction builds a transfer out transaction
func BuildTransferOutTransaction(userID, walletID int64, amount string, tokenID int64, transferID int64, recipientUsername string) (*TransactionRequest, error) {
	return NewTransactionBuilder().
		WithUser(userID).
		WithWallet(walletID).
		WithAmount(amount).
		WithToken(tokenID).
		WithFundType(FundTypeTransferOut).
		WithRelatedID(transferID).
		WithMetadata("transfer_id", transferID).
		WithMetadata("recipient", recipientUsername).
		Build()
}

// BuildDepositTransaction builds a deposit transaction
func BuildDepositTransaction(userID, walletID int64, amount string, tokenID int64, txHash string) (*TransactionRequest, error) {
	return NewTransactionBuilder().
		WithUser(userID).
		WithWallet(walletID).
		WithAmount(amount).
		WithToken(tokenID).
		WithFundType(FundTypeDeposit).
		WithReference(txHash).
		WithMetadata("tx_hash", txHash).
		Build()
}

// BuildWithdrawTransaction builds a withdrawal transaction
func BuildWithdrawTransaction(userID, walletID int64, amount string, tokenID int64, withdrawID int64, address string) (*TransactionRequest, error) {
	return NewTransactionBuilder().
		WithUser(userID).
		WithWallet(walletID).
		WithAmount(amount).
		WithToken(tokenID).
		WithFundType(FundTypeWithdraw).
		WithRelatedID(withdrawID).
		WithMetadata("withdraw_id", withdrawID).
		WithMetadata("address", address).
		Build()
}