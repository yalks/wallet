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
	
	// New fields for request context
	RequestSource    string // Request source (telegram, web, api, admin)
	RequestIP        string // User's IP address
	RequestUserAgent string // User's User-Agent
	
	// New fields for transaction details
	FeeAmount      string // Fee amount (as string to maintain precision)
	FeeType        string // Fee type (fixed, percentage)
	TargetUserID   int64  // Target user ID for transfers
	TargetUsername string // Target username for transfers
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

// WithRequestSource sets the request source
func (b *TransactionBuilder) WithRequestSource(source string) *TransactionBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.request.RequestSource = strings.TrimSpace(source)
	return b
}

// WithRequestIP sets the request IP
func (b *TransactionBuilder) WithRequestIP(ip string) *TransactionBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.request.RequestIP = strings.TrimSpace(ip)
	return b
}

// WithRequestUserAgent sets the request user agent
func (b *TransactionBuilder) WithRequestUserAgent(userAgent string) *TransactionBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.request.RequestUserAgent = strings.TrimSpace(userAgent)
	return b
}

// WithFeeAmount sets the fee amount
func (b *TransactionBuilder) WithFeeAmount(feeAmount string) *TransactionBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.request.FeeAmount = strings.TrimSpace(feeAmount)
	return b
}

// WithFeeType sets the fee type
func (b *TransactionBuilder) WithFeeType(feeType string) *TransactionBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.request.FeeType = strings.TrimSpace(feeType)
	return b
}

// WithTargetUser sets the target user information
func (b *TransactionBuilder) WithTargetUser(userID int64, username string) *TransactionBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.request.TargetUserID = userID
	b.request.TargetUsername = strings.TrimSpace(username)
	return b
}

// WithTokenSymbol sets the token symbol for fund operations
func (b *TransactionBuilder) WithTokenSymbol(symbol string) *TransactionBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.request.Metadata["token_symbol"] = strings.TrimSpace(symbol)
	return b
}

// WithBusinessID sets the business ID for idempotency
func (b *TransactionBuilder) WithBusinessID(businessID string) *TransactionBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.request.Metadata["business_id"] = strings.TrimSpace(businessID)
	return b
}

// WithDecimalAmount sets the amount as decimal.Decimal
func (b *TransactionBuilder) WithDecimalAmount(amount decimal.Decimal) *TransactionBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.request.Amount = amount.String()
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

// FundOperationBuilder 资金操作构建器
type FundOperationBuilder struct {
	userID           uint64
	tokenSymbol      string
	amount           decimal.Decimal
	businessID       string
	fundType         FundType
	description      string
	metadata         map[string]string
	relatedID        int64
	requestSource    string
	requestIP        string
	requestUserAgent string
	mu               sync.Mutex
}

// NewFundOperationBuilder 创建新的资金操作构建器
func NewFundOperationBuilder() *FundOperationBuilder {
	return &FundOperationBuilder{
		metadata: make(map[string]string),
	}
}

// WithUser 设置用户ID
func (b *FundOperationBuilder) WithUser(userID uint64) *FundOperationBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.userID = userID
	return b
}

// WithTokenSymbol 设置代币符号
func (b *FundOperationBuilder) WithTokenSymbol(symbol string) *FundOperationBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.tokenSymbol = strings.TrimSpace(symbol)
	return b
}

// WithAmount 设置操作金额
func (b *FundOperationBuilder) WithAmount(amount decimal.Decimal) *FundOperationBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.amount = amount
	return b
}

// WithBusinessID 设置业务ID
func (b *FundOperationBuilder) WithBusinessID(businessID string) *FundOperationBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.businessID = strings.TrimSpace(businessID)
	return b
}

// WithFundType 设置资金类型
func (b *FundOperationBuilder) WithFundType(fundType FundType) *FundOperationBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.fundType = fundType
	return b
}

// WithDescription 设置操作描述
func (b *FundOperationBuilder) WithDescription(description string) *FundOperationBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.description = strings.TrimSpace(description)
	return b
}

// WithMetadata 添加元数据
func (b *FundOperationBuilder) WithMetadata(key, value string) *FundOperationBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	if key != "" {
		b.metadata[key] = value
	}
	return b
}

// WithRelatedID 设置关联ID
func (b *FundOperationBuilder) WithRelatedID(id int64) *FundOperationBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.relatedID = id
	return b
}

// WithRequestSource 设置请求来源
func (b *FundOperationBuilder) WithRequestSource(source string) *FundOperationBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.requestSource = strings.TrimSpace(source)
	return b
}

// WithRequestIP 设置请求IP
func (b *FundOperationBuilder) WithRequestIP(ip string) *FundOperationBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.requestIP = strings.TrimSpace(ip)
	return b
}

// WithRequestUserAgent 设置请求User-Agent
func (b *FundOperationBuilder) WithRequestUserAgent(userAgent string) *FundOperationBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.requestUserAgent = strings.TrimSpace(userAgent)
	return b
}

// FundOperationRequest 资金操作请求结构
type FundOperationRequest struct {
	UserID      uint64             `json:"user_id" validate:"required"`
	TokenSymbol string             `json:"token_symbol" validate:"required"`
	Amount      decimal.Decimal    `json:"amount" validate:"required,gt=0"`
	BusinessID  string             `json:"business_id" validate:"required"`
	FundType    FundType           `json:"fund_type" validate:"required"`
	Description string             `json:"description"`
	Metadata    map[string]string  `json:"metadata,omitempty"`
	RelatedID   int64              `json:"related_id,omitempty"`
	
	RequestSource    string `json:"request_source,omitempty"`
	RequestIP        string `json:"request_ip,omitempty"`
	RequestUserAgent string `json:"request_user_agent,omitempty"`
}

// Build 构建并验证资金操作请求
func (b *FundOperationBuilder) Build() (*FundOperationRequest, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	if b.userID == 0 {
		return nil, fmt.Errorf("user ID is required")
	}
	if b.tokenSymbol == "" {
		return nil, fmt.Errorf("token symbol is required")
	}
	if b.amount.LessThanOrEqual(decimal.Zero) {
		return nil, fmt.Errorf("amount must be greater than zero")
	}
	if b.businessID == "" {
		return nil, fmt.Errorf("business ID is required")
	}
	if !IsValidFundType(b.fundType) {
		return nil, fmt.Errorf("invalid fund type: %s", b.fundType)
	}
	
	// 创建元数据副本
	metadata := make(map[string]string)
	for k, v := range b.metadata {
		metadata[k] = v
	}
	
	return &FundOperationRequest{
		UserID:           b.userID,
		TokenSymbol:      b.tokenSymbol,
		Amount:           b.amount,
		BusinessID:       b.businessID,
		FundType:         b.fundType,
		Description:      b.description,
		Metadata:         metadata,
		RelatedID:        b.relatedID,
		RequestSource:    b.requestSource,
		RequestIP:        b.requestIP,
		RequestUserAgent: b.requestUserAgent,
	}, nil
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
func BuildTransferOutTransaction(userID, walletID int64, amount string, tokenID int64, transferID int64, recipientUserID int64, recipientUsername string) (*TransactionRequest, error) {
	return NewTransactionBuilder().
		WithUser(userID).
		WithWallet(walletID).
		WithAmount(amount).
		WithToken(tokenID).
		WithFundType(FundTypeTransferOut).
		WithRelatedID(transferID).
		WithTargetUser(recipientUserID, recipientUsername).
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