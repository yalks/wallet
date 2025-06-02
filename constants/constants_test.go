package constants

import (
	"testing"
)

func TestOperationTypes(t *testing.T) {
	tests := []struct {
		name     string
		opType   OperationType
		expected OperationType
	}{
		{"credit operation", OperationTypeCredit, OperationTypeCredit},
		{"debit operation", OperationTypeDebit, OperationTypeDebit},
		{"transfer operation", OperationTypeTransfer, OperationTypeTransfer},
		{"freeze operation", OperationTypeFreeze, OperationTypeFreeze},
		{"unfreeze operation", OperationTypeUnfreeze, OperationTypeUnfreeze},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.opType != tt.expected {
				t.Errorf("Operation type mismatch: got %s, want %s", tt.opType, tt.expected)
			}
		})
	}
}

func TestFundTypes(t *testing.T) {
	// Test that fund types are properly defined
	fundTypes := []FundType{
		FundTypeRedPacketCreate,
		FundTypeRedPacketClaim,
		FundTypeRedPacketRefund,
		FundTypeRedPacketCancel,
		FundTypeTransferOut,
		FundTypeTransferIn,
		FundTypeDeposit,
		FundTypeWithdraw,
		FundTypeCommission,
		FundTypeReferralBonus,
		FundTypeSystemAdjustment,
	}
	
	// Verify fund types are unique
	seen := make(map[FundType]bool)
	for _, ft := range fundTypes {
		if seen[ft] {
			t.Errorf("Duplicate fund type: %s", ft)
		}
		seen[ft] = true
	}
	
	// Verify all fund types are valid
	for _, ft := range fundTypes {
		if !IsValidFundType(ft) {
			t.Errorf("Fund type %s is not valid", ft)
		}
	}
	
	// Test GetFundDirection
	testCases := []struct {
		fundType FundType
		expected FundDirection
	}{
		{FundTypeDeposit, FundDirectionIn},
		{FundTypeWithdraw, FundDirectionOut},
		{FundTypeTransferOut, FundDirectionOut},
		{FundTypeTransferIn, FundDirectionIn},
	}
	
	for _, tc := range testCases {
		direction := GetFundDirection(tc.fundType)
		if direction != tc.expected {
			t.Errorf("GetFundDirection(%s) = %s, want %s", tc.fundType, direction, tc.expected)
		}
	}
}

func TestTransactionStatus(t *testing.T) {
	// Test transaction status constants
	statuses := []TransactionStatus{
		TransactionStatusPending,
		TransactionStatusProcessing,
		TransactionStatusCompleted,
		TransactionStatusFailed,
		TransactionStatusCancelled,
		TransactionStatusExpired,
		TransactionStatusRefunded,
	}
	
	// Verify all statuses are non-empty
	for _, status := range statuses {
		if status == "" {
			t.Error("Transaction status should not be empty")
		}
	}
	
	// Verify statuses are unique
	seen := make(map[TransactionStatus]bool)
	for _, status := range statuses {
		if seen[status] {
			t.Errorf("Duplicate transaction status: %s", status)
		}
		seen[status] = true
	}
	
	// Test IsValidTransactionStatus
	for _, status := range statuses {
		if !IsValidTransactionStatus(status) {
			t.Errorf("Status %s should be valid", status)
		}
	}
	
	// Test IsFinalStatus
	finalStatuses := []TransactionStatus{
		TransactionStatusCompleted,
		TransactionStatusFailed,
		TransactionStatusCancelled,
		TransactionStatusRefunded,
	}
	
	for _, status := range finalStatuses {
		if !IsFinalStatus(status) {
			t.Errorf("Status %s should be final", status)
		}
	}
}

func TestWalletTypes(t *testing.T) {
	// Test wallet type constants
	walletTypes := []WalletType{
		WalletTypeAvailable,
		WalletTypeFrozen,
		WalletTypePending,
	}
	
	// Verify all wallet types are valid
	for _, wt := range walletTypes {
		if !IsValidWalletType(wt) {
			t.Errorf("Wallet type %s should be valid", wt)
		}
	}
	
	// Test invalid wallet type
	if IsValidWalletType("invalid") {
		t.Error("Invalid wallet type should not be valid")
	}
}

func TestGetOperationTypeForFundType(t *testing.T) {
	tests := []struct {
		fundType FundType
		expected OperationType
	}{
		{FundTypeDeposit, OperationTypeCredit},
		{FundTypeWithdraw, OperationTypeDebit},
		{FundTypeTransferOut, OperationTypeTransfer},
		{FundTypeTransferIn, OperationTypeTransfer},
		{FundTypeCommission, OperationTypeCredit},
		{FundTypeRedPacketCreate, OperationTypeDebit},
		{FundTypeRedPacketClaim, OperationTypeCredit},
	}
	
	for _, tt := range tests {
		t.Run(string(tt.fundType), func(t *testing.T) {
			opType := GetOperationTypeForFundType(tt.fundType)
			if opType != tt.expected {
				t.Errorf("GetOperationTypeForFundType(%s) = %s, want %s", 
					tt.fundType, opType, tt.expected)
			}
		})
	}
}

func TestValidationFunctions(t *testing.T) {
	// Test IsValidOperationType
	validOps := []OperationType{
		OperationTypeCredit,
		OperationTypeDebit,
		OperationTypeTransfer,
		OperationTypeFreeze,
		OperationTypeUnfreeze,
	}
	
	for _, op := range validOps {
		if !IsValidOperationType(op) {
			t.Errorf("Operation type %s should be valid", op)
		}
	}
	
	if IsValidOperationType("invalid") {
		t.Error("Invalid operation type should not be valid")
	}
}

// Benchmark tests
func BenchmarkOperationTypeComparison(b *testing.B) {
	opType := OperationTypeCredit
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = opType == OperationTypeCredit
	}
}

func BenchmarkFundTypeValidation(b *testing.B) {
	fundType := FundTypeRedPacketCreate
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = IsValidFundType(fundType)
	}
}

func BenchmarkGetFundDirection(b *testing.B) {
	fundType := FundTypeDeposit
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GetFundDirection(fundType)
	}
}