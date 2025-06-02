package logic

import (
	"github.com/shopspring/decimal"
	"github.com/yalks/wallet/constants"
)

// FeeConfig represents fee configuration for different operations
type FeeConfig struct {
	FeeType   string          // "fixed" or "percentage"
	FeeAmount decimal.Decimal // Amount for fixed fee or percentage (0.01 = 1%)
	MinFee    decimal.Decimal // Minimum fee amount
	MaxFee    decimal.Decimal // Maximum fee amount
}

// FeeCalculator calculates fees for different transaction types
type FeeCalculator struct {
	feeConfigs map[constants.FundType]FeeConfig
}

// NewFeeCalculator creates a new fee calculator with default configurations
func NewFeeCalculator() *FeeCalculator {
	return &FeeCalculator{
		feeConfigs: getDefaultFeeConfigs(),
	}
}

// CalculateFee calculates the fee for a given fund type and amount
func (fc *FeeCalculator) CalculateFee(fundType constants.FundType, amount decimal.Decimal) (fee decimal.Decimal, feeType string) {
	config, exists := fc.feeConfigs[fundType]
	if !exists {
		// No fee for operations without fee config
		return decimal.Zero, ""
	}

	feeType = config.FeeType
	
	switch config.FeeType {
	case "fixed":
		fee = config.FeeAmount
	case "percentage":
		fee = amount.Mul(config.FeeAmount)
		
		// Apply min/max limits
		if !config.MinFee.IsZero() && fee.LessThan(config.MinFee) {
			fee = config.MinFee
		}
		if !config.MaxFee.IsZero() && fee.GreaterThan(config.MaxFee) {
			fee = config.MaxFee
		}
	default:
		fee = decimal.Zero
		feeType = ""
	}
	
	return fee, feeType
}

// getDefaultFeeConfigs returns default fee configurations
func getDefaultFeeConfigs() map[constants.FundType]FeeConfig {
	return map[constants.FundType]FeeConfig{
		// Withdrawal fees
		constants.FundTypeWithdraw: {
			FeeType:   "percentage",
			FeeAmount: decimal.NewFromFloat(0.002), // 0.2%
			MinFee:    decimal.NewFromFloat(0.1),   // Minimum 0.1
			MaxFee:    decimal.NewFromFloat(100),   // Maximum 100
		},
		
		// Transfer fees (optional, can be removed if transfers are free)
		constants.FundTypeTransferOut: {
			FeeType:   "fixed",
			FeeAmount: decimal.NewFromFloat(0.01), // Fixed 0.01 fee
			MinFee:    decimal.Zero,
			MaxFee:    decimal.Zero,
		},
		
		// Exchange fees (currency conversion)
		constants.FundTypeExchangeOut: {
			FeeType:   "percentage",
			FeeAmount: decimal.NewFromFloat(0.001), // 0.1%
			MinFee:    decimal.Zero,
			MaxFee:    decimal.Zero,
		},
		
		// No fees for these operations
		constants.FundTypeDeposit:         {},
		constants.FundTypeTransferIn:      {},
		constants.FundTypeRedPacketCreate: {},
		constants.FundTypeRedPacketClaim:  {},
		constants.FundTypeRedPacketRefund: {},
		constants.FundTypeAdminAdd:        {},
		constants.FundTypeAdminDeduct:     {},
		constants.FundTypeSystemAdjustment:{},
	}
}

// SetFeeConfig allows updating fee configuration for a specific fund type
func (fc *FeeCalculator) SetFeeConfig(fundType constants.FundType, config FeeConfig) {
	fc.feeConfigs[fundType] = config
}

// GetFeeConfig returns the fee configuration for a specific fund type
func (fc *FeeCalculator) GetFeeConfig(fundType constants.FundType) (FeeConfig, bool) {
	config, exists := fc.feeConfigs[fundType]
	return config, exists
}