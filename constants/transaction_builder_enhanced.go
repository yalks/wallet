package constants

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

const (
	// MaxAmountDigits 最大金额位数（包括小数点前后）
	MaxAmountDigits = 36
	// MaxDecimalPlaces 最大小数位数
	MaxDecimalPlaces = 18
	// MaxTransactionAmount 单笔交易最大金额
	MaxTransactionAmount = "1000000" // 100万
	// MinTransactionAmount 单笔交易最小金额
	MinTransactionAmount = "0.000000000000000001" // 最小精度
)

var (
	// amountRegex 金额格式正则表达式
	amountRegex = regexp.MustCompile(`^-?\d+(\.\d+)?$`)
	// referenceRegex 引用号格式正则
	referenceRegex = regexp.MustCompile(`^[a-zA-Z0-9_\-]{1,100}$`)
)

// TransactionRequestEnhanced 增强版交易请求，包含更多验证字段
type TransactionRequestEnhanced struct {
	UserID       int64
	WalletID     int64
	Amount       string
	TokenID      int64
	FundType     FundType
	Reference    string
	Description  string
	RelatedID    int64
	Metadata     map[string]interface{}
	
	// 新增字段
	IdempotencyKey string                 // 幂等键，防止重复提交
	ExpireAt       *time.Time             // 交易过期时间
	Priority       int                    // 优先级（1-10）
	Tags           []string               // 标签列表
	Callback       map[string]string      // 回调信息
}

// TransactionBuilderEnhanced 增强版交易构建器
type TransactionBuilderEnhanced struct {
	request TransactionRequestEnhanced
	mu      sync.Mutex
	errors  []string // 收集验证错误
}

// NewTransactionBuilderEnhanced 创建增强版交易构建器
func NewTransactionBuilderEnhanced() *TransactionBuilderEnhanced {
	return &TransactionBuilderEnhanced{
		request: TransactionRequestEnhanced{
			Metadata: make(map[string]interface{}),
			Tags:     make([]string, 0),
			Callback: make(map[string]string),
			Priority: 5, // 默认中等优先级
		},
		errors: make([]string, 0),
	}
}

// WithUser 设置用户ID（线程安全）
func (b *TransactionBuilderEnhanced) WithUser(userID int64) *TransactionBuilderEnhanced {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	if userID <= 0 {
		b.errors = append(b.errors, "user ID must be positive")
	}
	b.request.UserID = userID
	return b
}

// WithWallet 设置钱包ID（线程安全）
func (b *TransactionBuilderEnhanced) WithWallet(walletID int64) *TransactionBuilderEnhanced {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	if walletID <= 0 {
		b.errors = append(b.errors, "wallet ID must be positive")
	}
	b.request.WalletID = walletID
	return b
}

// WithAmount 设置金额（支持字符串和decimal类型）
func (b *TransactionBuilderEnhanced) WithAmount(amount interface{}) *TransactionBuilderEnhanced {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	var amountStr string
	
	switch v := amount.(type) {
	case string:
		amountStr = strings.TrimSpace(v)
	case decimal.Decimal:
		amountStr = v.String()
	case float64:
		// 使用 decimal 避免浮点数精度问题
		amountStr = decimal.NewFromFloat(v).String()
	case int64:
		amountStr = fmt.Sprintf("%d", v)
	case int:
		amountStr = fmt.Sprintf("%d", v)
	default:
		b.errors = append(b.errors, fmt.Sprintf("unsupported amount type: %T", amount))
		return b
	}
	
	// 验证金额格式
	if err := b.validateAmount(amountStr); err != nil {
		b.errors = append(b.errors, err.Error())
	}
	
	b.request.Amount = amountStr
	return b
}

// validateAmount 验证金额
func (b *TransactionBuilderEnhanced) validateAmount(amount string) error {
	// 检查格式
	if !amountRegex.MatchString(amount) {
		return fmt.Errorf("invalid amount format: %s", amount)
	}
	
	// 使用 decimal 进行精确验证
	d, err := decimal.NewFromString(amount)
	if err != nil {
		return fmt.Errorf("invalid amount: %v", err)
	}
	
	// 检查是否为负数
	if d.IsNegative() {
		return fmt.Errorf("amount cannot be negative: %s", amount)
	}
	
	// 检查是否为零
	if d.IsZero() {
		return fmt.Errorf("amount cannot be zero")
	}
	
	// 检查最小值
	minAmount, _ := decimal.NewFromString(MinTransactionAmount)
	if d.LessThan(minAmount) {
		return fmt.Errorf("amount %s is less than minimum %s", amount, MinTransactionAmount)
	}
	
	// 检查最大值
	maxAmount, _ := decimal.NewFromString(MaxTransactionAmount)
	if d.GreaterThan(maxAmount) {
		return fmt.Errorf("amount %s exceeds maximum %s", amount, MaxTransactionAmount)
	}
	
	// 检查小数位数
	parts := strings.Split(amount, ".")
	if len(parts) == 2 && len(parts[1]) > MaxDecimalPlaces {
		return fmt.Errorf("amount has too many decimal places (max %d)", MaxDecimalPlaces)
	}
	
	// 检查总位数
	cleanAmount := strings.Replace(amount, ".", "", 1)
	if len(cleanAmount) > MaxAmountDigits {
		return fmt.Errorf("amount has too many digits (max %d)", MaxAmountDigits)
	}
	
	return nil
}

// WithToken 设置代币ID
func (b *TransactionBuilderEnhanced) WithToken(tokenID int64) *TransactionBuilderEnhanced {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	if tokenID <= 0 {
		b.errors = append(b.errors, "token ID must be positive")
	}
	b.request.TokenID = tokenID
	return b
}

// WithFundType 设置资金类型
func (b *TransactionBuilderEnhanced) WithFundType(fundType FundType) *TransactionBuilderEnhanced {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	if !IsValidFundType(fundType) {
		b.errors = append(b.errors, fmt.Sprintf("invalid fund type: %s", fundType))
	}
	b.request.FundType = fundType
	return b
}

// WithReference 设置引用号（支持自定义验证）
func (b *TransactionBuilderEnhanced) WithReference(reference string) *TransactionBuilderEnhanced {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	reference = strings.TrimSpace(reference)
	if reference != "" && !referenceRegex.MatchString(reference) {
		b.errors = append(b.errors, fmt.Sprintf("invalid reference format: %s", reference))
	}
	b.request.Reference = reference
	return b
}

// WithDescription 设置描述（限制长度）
func (b *TransactionBuilderEnhanced) WithDescription(description string) *TransactionBuilderEnhanced {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	description = strings.TrimSpace(description)
	if len(description) > 500 {
		b.errors = append(b.errors, "description exceeds 500 characters")
		description = description[:500]
	}
	b.request.Description = description
	return b
}

// WithRelatedID 设置关联ID
func (b *TransactionBuilderEnhanced) WithRelatedID(id int64) *TransactionBuilderEnhanced {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	if id < 0 {
		b.errors = append(b.errors, "related ID cannot be negative")
	}
	b.request.RelatedID = id
	return b
}

// WithMetadata 添加元数据（支持批量）
func (b *TransactionBuilderEnhanced) WithMetadata(data map[string]interface{}) *TransactionBuilderEnhanced {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	for key, value := range data {
		// 验证 key
		if key == "" {
			b.errors = append(b.errors, "metadata key cannot be empty")
			continue
		}
		if len(key) > 100 {
			b.errors = append(b.errors, fmt.Sprintf("metadata key '%s' exceeds 100 characters", key))
			continue
		}
		
		// 验证 value 类型
		switch v := value.(type) {
		case nil:
			// 跳过 nil 值
			continue
		case string:
			if len(v) > 1000 {
				b.errors = append(b.errors, fmt.Sprintf("metadata value for key '%s' exceeds 1000 characters", key))
				continue
			}
		}
		
		b.request.Metadata[key] = value
	}
	return b
}

// WithIdempotencyKey 设置幂等键
func (b *TransactionBuilderEnhanced) WithIdempotencyKey(key string) *TransactionBuilderEnhanced {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	key = strings.TrimSpace(key)
	if key == "" {
		b.errors = append(b.errors, "idempotency key cannot be empty")
	}
	if len(key) > 128 {
		b.errors = append(b.errors, "idempotency key exceeds 128 characters")
	}
	b.request.IdempotencyKey = key
	return b
}

// WithExpireAt 设置过期时间
func (b *TransactionBuilderEnhanced) WithExpireAt(expireAt time.Time) *TransactionBuilderEnhanced {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	if expireAt.Before(time.Now()) {
		b.errors = append(b.errors, "expire time cannot be in the past")
	}
	b.request.ExpireAt = &expireAt
	return b
}

// WithPriority 设置优先级（1-10）
func (b *TransactionBuilderEnhanced) WithPriority(priority int) *TransactionBuilderEnhanced {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	if priority < 1 || priority > 10 {
		b.errors = append(b.errors, "priority must be between 1 and 10")
		priority = 5 // 设置默认值
	}
	b.request.Priority = priority
	return b
}

// WithTags 添加标签
func (b *TransactionBuilderEnhanced) WithTags(tags ...string) *TransactionBuilderEnhanced {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		if len(tag) > 50 {
			b.errors = append(b.errors, fmt.Sprintf("tag '%s' exceeds 50 characters", tag))
			continue
		}
		// 去重
		exists := false
		for _, t := range b.request.Tags {
			if t == tag {
				exists = true
				break
			}
		}
		if !exists {
			b.request.Tags = append(b.request.Tags, tag)
		}
	}
	return b
}

// WithCallback 设置回调信息
func (b *TransactionBuilderEnhanced) WithCallback(url, method string) *TransactionBuilderEnhanced {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	if url != "" {
		b.request.Callback["url"] = url
	}
	if method != "" {
		method = strings.ToUpper(method)
		if method != "GET" && method != "POST" {
			b.errors = append(b.errors, "callback method must be GET or POST")
		} else {
			b.request.Callback["method"] = method
		}
	}
	return b
}

// Build 构建并验证交易请求
func (b *TransactionBuilderEnhanced) Build() (*TransactionRequestEnhanced, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	// 检查构建过程中的错误
	if len(b.errors) > 0 {
		return nil, fmt.Errorf("validation errors: %s", strings.Join(b.errors, "; "))
	}
	
	// 验证必填字段
	if b.request.UserID == 0 {
		return nil, fmt.Errorf("user ID is required")
	}
	if b.request.WalletID == 0 {
		return nil, fmt.Errorf("wallet ID is required")
	}
	if b.request.Amount == "" {
		return nil, fmt.Errorf("amount is required")
	}
	if b.request.TokenID == 0 {
		return nil, fmt.Errorf("token ID is required")
	}
	if !IsValidFundType(b.request.FundType) {
		return nil, fmt.Errorf("valid fund type is required")
	}
	
	// 生成默认值
	if b.request.Reference == "" {
		b.request.Reference = b.generateReference()
	}
	
	if b.request.Description == "" {
		if info, exists := GetFundTypeInfo(b.request.FundType); exists {
			b.request.Description = info.Description
		}
	}
	
	if b.request.IdempotencyKey == "" {
		b.request.IdempotencyKey = b.generateIdempotencyKey()
	}
	
	// 返回副本，避免后续修改
	result := b.request
	result.Metadata = make(map[string]interface{})
	for k, v := range b.request.Metadata {
		result.Metadata[k] = v
	}
	result.Tags = make([]string, len(b.request.Tags))
	copy(result.Tags, b.request.Tags)
	result.Callback = make(map[string]string)
	for k, v := range b.request.Callback {
		result.Callback[k] = v
	}
	
	return &result, nil
}

// generateReference 生成唯一引用号
func (b *TransactionBuilderEnhanced) generateReference() string {
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%s_%d_%d", b.request.FundType, b.request.UserID, timestamp)
}

// generateIdempotencyKey 生成幂等键
func (b *TransactionBuilderEnhanced) generateIdempotencyKey() string {
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%d_%d_%s_%d", b.request.UserID, b.request.WalletID, b.request.FundType, timestamp)
}

// Reset 重置构建器（可重用）
func (b *TransactionBuilderEnhanced) Reset() *TransactionBuilderEnhanced {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	b.request = TransactionRequestEnhanced{
		Metadata: make(map[string]interface{}),
		Tags:     make([]string, 0),
		Callback: make(map[string]string),
		Priority: 5,
	}
	b.errors = make([]string, 0)
	return b
}

// Validate 仅验证不构建（用于预检查）
func (b *TransactionBuilderEnhanced) Validate() error {
	_, err := b.Build()
	return err
}

// Clone 克隆构建器（用于创建相似交易）
func (b *TransactionBuilderEnhanced) Clone() *TransactionBuilderEnhanced {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	newBuilder := NewTransactionBuilderEnhanced()
	newBuilder.request = b.request
	
	// 深拷贝 map 和 slice
	newBuilder.request.Metadata = make(map[string]interface{})
	for k, v := range b.request.Metadata {
		newBuilder.request.Metadata[k] = v
	}
	
	newBuilder.request.Tags = make([]string, len(b.request.Tags))
	copy(newBuilder.request.Tags, b.request.Tags)
	
	newBuilder.request.Callback = make(map[string]string)
	for k, v := range b.request.Callback {
		newBuilder.request.Callback[k] = v
	}
	
	return newBuilder
}

// 便捷构建器函数

// BuildBatchTransfer 批量转账构建器
func BuildBatchTransfer(fromUserID, fromWalletID int64, transfers []struct {
	ToUserID   int64
	ToWalletID int64
	Amount     string
	TokenID    int64
}) ([]*TransactionRequestEnhanced, error) {
	requests := make([]*TransactionRequestEnhanced, 0, len(transfers)*2)
	
	for i, transfer := range transfers {
		// 转出交易
		outReq, err := NewTransactionBuilderEnhanced().
			WithUser(fromUserID).
			WithWallet(fromWalletID).
			WithAmount(transfer.Amount).
			WithToken(transfer.TokenID).
			WithFundType(FundTypeTransferOut).
			WithMetadata(map[string]interface{}{
				"batch_index": i,
				"to_user_id":  transfer.ToUserID,
				"to_wallet_id": transfer.ToWalletID,
			}).
			WithTags("batch_transfer").
			Build()
		if err != nil {
			return nil, fmt.Errorf("build out transaction %d: %v", i, err)
		}
		requests = append(requests, outReq)
		
		// 转入交易
		inReq, err := NewTransactionBuilderEnhanced().
			WithUser(transfer.ToUserID).
			WithWallet(transfer.ToWalletID).
			WithAmount(transfer.Amount).
			WithToken(transfer.TokenID).
			WithFundType(FundTypeTransferIn).
			WithMetadata(map[string]interface{}{
				"batch_index": i,
				"from_user_id": fromUserID,
				"from_wallet_id": fromWalletID,
			}).
			WithTags("batch_transfer").
			Build()
		if err != nil {
			return nil, fmt.Errorf("build in transaction %d: %v", i, err)
		}
		requests = append(requests, inReq)
	}
	
	return requests, nil
}

// BuildScheduledTransaction 构建定时交易
func BuildScheduledTransaction(builder *TransactionBuilderEnhanced, scheduleAt time.Time) (*TransactionRequestEnhanced, error) {
	return builder.
		WithExpireAt(scheduleAt.Add(24 * time.Hour)). // 24小时后过期
		WithMetadata(map[string]interface{}{
			"scheduled_at": scheduleAt.Format(time.RFC3339),
			"status": "scheduled",
		}).
		WithTags("scheduled").
		Build()
}

// BuildRecurringTransaction 构建循环交易
func BuildRecurringTransaction(builder *TransactionBuilderEnhanced, interval time.Duration, count int) (*TransactionRequestEnhanced, error) {
	return builder.
		WithMetadata(map[string]interface{}{
			"recurring_interval": interval.String(),
			"recurring_count": count,
			"recurring_status": "active",
		}).
		WithTags("recurring").
		WithPriority(7). // 较高优先级
		Build()
}