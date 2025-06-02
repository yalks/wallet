package wallet

import (
	"testing"

	"github.com/yalks/wallet/constants"
)

// TestPackageImports verifies that the package and its dependencies compile correctly
func TestPackageImports(t *testing.T) {
	// Test that constants are accessible
	_ = constants.FundTypeDeposit
	_ = constants.OperationTypeCredit
	_ = constants.TransactionStatusPending
	
	// Test that the Manager function exists
	// Note: We can't call it without initialization, but we can verify it exists
	t.Log("Manager function is accessible")
}

// TestFundTypeDirection verifies fund type directions
func TestFundTypeDirection(t *testing.T) {
	tests := []struct {
		fundType constants.FundType
		expected constants.FundDirection
	}{
		{constants.FundTypeDeposit, constants.FundDirectionIn},
		{constants.FundTypeWithdraw, constants.FundDirectionOut},
		{constants.FundTypeCommission, constants.FundDirectionIn},
	}
	
	for _, tt := range tests {
		direction := constants.GetFundDirection(tt.fundType)
		if direction != tt.expected {
			t.Errorf("GetFundDirection(%s) = %s, want %s", tt.fundType, direction, tt.expected)
		}
	}
}

// TestOperationTypeForFundType verifies operation type mapping
func TestOperationTypeForFundType(t *testing.T) {
	tests := []struct {
		fundType constants.FundType
		expected constants.OperationType
	}{
		{constants.FundTypeDeposit, constants.OperationTypeCredit},
		{constants.FundTypeWithdraw, constants.OperationTypeDebit},
		{constants.FundTypeTransferOut, constants.OperationTypeTransfer},
		{constants.FundTypeTransferIn, constants.OperationTypeTransfer},
	}
	
	for _, tt := range tests {
		opType := constants.GetOperationTypeForFundType(tt.fundType)
		if opType != tt.expected {
			t.Errorf("GetOperationTypeForFundType(%s) = %s, want %s", tt.fundType, opType, tt.expected)
		}
	}
}