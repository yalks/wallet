package constants

// FundType represents the type of fund transaction
type FundType string

// Fund type constants - all possible fund transaction types
const (
	// Red Packet related fund types
	FundTypeRedPacketCreate     FundType = "red_packet_create"      // 红包创建扣费
	FundTypeRedPacketClaim      FundType = "red_packet_claim"       // 领取红包
	FundTypeRedPacketRefund     FundType = "red_packet_refund"      // 红包过期退回
	FundTypeRedPacketCancel     FundType = "red_packet_cancel"      // 手动撤回红包

	// Transfer related fund types
	FundTypeTransferOut         FundType = "transfer_out"           // 转账发起(扣款)
	FundTypeTransferIn          FundType = "transfer_in"            // 转账收款(入账)
	FundTypeTransferExpired     FundType = "transfer_expired"       // 转账过期退回

	// Payment Request related fund types
	FundTypePaymentRequest      FundType = "payment_request"        // 收款请求
	FundTypePaymentOut          FundType = "payment_out"            // 付款(扣款)
	FundTypePaymentIn           FundType = "payment_in"             // 收款(入账)

	// Deposit and Withdraw fund types
	FundTypeDeposit             FundType = "deposit"                // 充值
	FundTypeWithdraw            FundType = "withdraw"               // 提现
	FundTypeWithdrawRefund      FundType = "withdraw_refund"        // 提现退回

	// Admin operation fund types
	FundTypeAdminAdd            FundType = "admin_add"              // 后台增加
	FundTypeAdminDeduct         FundType = "admin_deduct"           // 后台减少

	// Exchange related fund types
	FundTypeExchangeOut         FundType = "exchange_out"           // 兑换减少
	FundTypeExchangeIn          FundType = "exchange_in"            // 兑换增加

	// Other fund types
	FundTypeCommission          FundType = "commission"             // 佣金
	FundTypeReferralBonus       FundType = "referral_bonus"         // 推荐奖励
	FundTypeSystemAdjustment    FundType = "system_adjustment"      // 系统调整
)

// FundDirection represents the direction of fund flow
type FundDirection string

const (
	FundDirectionIn  FundDirection = "in"  // 资金流入
	FundDirectionOut FundDirection = "out" // 资金流出
)

// FundTypeInfo contains metadata about a fund type
type FundTypeInfo struct {
	Type        FundType
	Direction   FundDirection
	Description string
	Category    string
}

// fundTypeRegistry maps fund types to their info
var fundTypeRegistry = map[FundType]FundTypeInfo{
	// Red Packet operations
	FundTypeRedPacketCreate: {
		Type:        FundTypeRedPacketCreate,
		Direction:   FundDirectionOut,
		Description: "Red packet creation fee",
		Category:    "red_packet",
	},
	FundTypeRedPacketClaim: {
		Type:        FundTypeRedPacketClaim,
		Direction:   FundDirectionIn,
		Description: "Red packet claim",
		Category:    "red_packet",
	},
	FundTypeRedPacketRefund: {
		Type:        FundTypeRedPacketRefund,
		Direction:   FundDirectionIn,
		Description: "Red packet expiry refund",
		Category:    "red_packet",
	},
	FundTypeRedPacketCancel: {
		Type:        FundTypeRedPacketCancel,
		Direction:   FundDirectionIn,
		Description: "Red packet manual cancel refund",
		Category:    "red_packet",
	},

	// Transfer operations
	FundTypeTransferOut: {
		Type:        FundTypeTransferOut,
		Direction:   FundDirectionOut,
		Description: "Transfer initiation",
		Category:    "transfer",
	},
	FundTypeTransferIn: {
		Type:        FundTypeTransferIn,
		Direction:   FundDirectionIn,
		Description: "Transfer receipt",
		Category:    "transfer",
	},
	FundTypeTransferExpired: {
		Type:        FundTypeTransferExpired,
		Direction:   FundDirectionIn,
		Description: "Transfer expiry refund",
		Category:    "transfer",
	},

	// Payment operations
	FundTypePaymentRequest: {
		Type:        FundTypePaymentRequest,
		Direction:   FundDirectionIn,
		Description: "Payment request",
		Category:    "payment",
	},
	FundTypePaymentOut: {
		Type:        FundTypePaymentOut,
		Direction:   FundDirectionOut,
		Description: "Payment made",
		Category:    "payment",
	},
	FundTypePaymentIn: {
		Type:        FundTypePaymentIn,
		Direction:   FundDirectionIn,
		Description: "Payment received",
		Category:    "payment",
	},

	// Deposit and Withdraw
	FundTypeDeposit: {
		Type:        FundTypeDeposit,
		Direction:   FundDirectionIn,
		Description: "Deposit",
		Category:    "wallet",
	},
	FundTypeWithdraw: {
		Type:        FundTypeWithdraw,
		Direction:   FundDirectionOut,
		Description: "Withdrawal",
		Category:    "wallet",
	},
	FundTypeWithdrawRefund: {
		Type:        FundTypeWithdrawRefund,
		Direction:   FundDirectionIn,
		Description: "Withdrawal refund",
		Category:    "wallet",
	},

	// Admin operations
	FundTypeAdminAdd: {
		Type:        FundTypeAdminAdd,
		Direction:   FundDirectionIn,
		Description: "Admin balance addition",
		Category:    "admin",
	},
	FundTypeAdminDeduct: {
		Type:        FundTypeAdminDeduct,
		Direction:   FundDirectionOut,
		Description: "Admin balance deduction",
		Category:    "admin",
	},

	// Exchange operations
	FundTypeExchangeOut: {
		Type:        FundTypeExchangeOut,
		Direction:   FundDirectionOut,
		Description: "Exchange deduction",
		Category:    "exchange",
	},
	FundTypeExchangeIn: {
		Type:        FundTypeExchangeIn,
		Direction:   FundDirectionIn,
		Description: "Exchange addition",
		Category:    "exchange",
	},

	// Others
	FundTypeCommission: {
		Type:        FundTypeCommission,
		Direction:   FundDirectionIn,
		Description: "Commission earned",
		Category:    "bonus",
	},
	FundTypeReferralBonus: {
		Type:        FundTypeReferralBonus,
		Direction:   FundDirectionIn,
		Description: "Referral bonus",
		Category:    "bonus",
	},
	FundTypeSystemAdjustment: {
		Type:        FundTypeSystemAdjustment,
		Direction:   FundDirectionIn, // Can be both, default to in
		Description: "System adjustment",
		Category:    "system",
	},
}

// GetFundTypeInfo returns the info for a given fund type
func GetFundTypeInfo(fundType FundType) (FundTypeInfo, bool) {
	info, exists := fundTypeRegistry[fundType]
	return info, exists
}

// GetFundDirection returns the direction for a given fund type
func GetFundDirection(fundType FundType) FundDirection {
	if info, exists := fundTypeRegistry[fundType]; exists {
		return info.Direction
	}
	return ""
}

// IsValidFundType checks if a fund type is valid
func IsValidFundType(fundType FundType) bool {
	_, exists := fundTypeRegistry[fundType]
	return exists
}

// GetFundTypesByCategory returns all fund types for a given category
func GetFundTypesByCategory(category string) []FundType {
	var types []FundType
	for _, info := range fundTypeRegistry {
		if info.Category == category {
			types = append(types, info.Type)
		}
	}
	return types
}

// GetAllFundTypes returns all registered fund types
func GetAllFundTypes() []FundType {
	var types []FundType
	for fundType := range fundTypeRegistry {
		types = append(types, fundType)
	}
	return types
}