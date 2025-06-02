package constants

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionBuilder_Basic(t *testing.T) {
	builder := NewTransactionBuilder()
	
	req, err := builder.
		WithUser(123).
		WithWallet(456).
		WithAmount("100.50").
		WithToken(1).
		WithFundType(FundTypeDeposit).
		WithReference("test_ref_123").
		WithDescription("Test deposit").
		WithRelatedID(789).
		WithMetadata("key1", "value1").
		WithMetadata("key2", 123).
		Build()
	
	require.NoError(t, err)
	assert.Equal(t, int64(123), req.UserID)
	assert.Equal(t, int64(456), req.WalletID)
	assert.Equal(t, "100.50", req.Amount)
	assert.Equal(t, int64(1), req.TokenID)
	assert.Equal(t, FundTypeDeposit, req.FundType)
	assert.Equal(t, "test_ref_123", req.Reference)
	assert.Equal(t, "Test deposit", req.Description)
	assert.Equal(t, int64(789), req.RelatedID)
	assert.Equal(t, "value1", req.Metadata["key1"])
	assert.Equal(t, 123, req.Metadata["key2"])
}

func TestTransactionBuilder_Validation(t *testing.T) {
	tests := []struct {
		name      string
		build     func() (*TransactionRequest, error)
		wantError string
	}{
		{
			name: "missing user ID",
			build: func() (*TransactionRequest, error) {
				return NewTransactionBuilder().
					WithWallet(456).
					WithAmount("100").
					WithToken(1).
					WithFundType(FundTypeDeposit).
					Build()
			},
			wantError: "user ID must be positive",
		},
		{
			name: "negative user ID",
			build: func() (*TransactionRequest, error) {
				return NewTransactionBuilder().
					WithUser(-1).
					WithWallet(456).
					WithAmount("100").
					WithToken(1).
					WithFundType(FundTypeDeposit).
					Build()
			},
			wantError: "user ID must be positive",
		},
		{
			name: "missing wallet ID",
			build: func() (*TransactionRequest, error) {
				return NewTransactionBuilder().
					WithUser(123).
					WithAmount("100").
					WithToken(1).
					WithFundType(FundTypeDeposit).
					Build()
			},
			wantError: "wallet ID must be positive",
		},
		{
			name: "missing amount",
			build: func() (*TransactionRequest, error) {
				return NewTransactionBuilder().
					WithUser(123).
					WithWallet(456).
					WithToken(1).
					WithFundType(FundTypeDeposit).
					Build()
			},
			wantError: "amount is required",
		},
		{
			name: "zero amount",
			build: func() (*TransactionRequest, error) {
				return NewTransactionBuilder().
					WithUser(123).
					WithWallet(456).
					WithAmount("0").
					WithToken(1).
					WithFundType(FundTypeDeposit).
					Build()
			},
			wantError: "must be greater than zero",
		},
		{
			name: "negative amount",
			build: func() (*TransactionRequest, error) {
				return NewTransactionBuilder().
					WithUser(123).
					WithWallet(456).
					WithAmount("-100").
					WithToken(1).
					WithFundType(FundTypeDeposit).
					Build()
			},
			wantError: "must be a positive number",
		},
		{
			name: "invalid amount format",
			build: func() (*TransactionRequest, error) {
				return NewTransactionBuilder().
					WithUser(123).
					WithWallet(456).
					WithAmount("abc").
					WithToken(1).
					WithFundType(FundTypeDeposit).
					Build()
			},
			wantError: "must be a positive number",
		},
		{
			name: "amount exceeds maximum",
			build: func() (*TransactionRequest, error) {
				return NewTransactionBuilder().
					WithUser(123).
					WithWallet(456).
					WithAmount("1000001").
					WithToken(1).
					WithFundType(FundTypeDeposit).
					Build()
			},
			wantError: "exceeds maximum amount",
		},
		{
			name: "too many decimal places",
			build: func() (*TransactionRequest, error) {
				return NewTransactionBuilder().
					WithUser(123).
					WithWallet(456).
					WithAmount("0.1234567890123456789").
					WithToken(1).
					WithFundType(FundTypeDeposit).
					Build()
			},
			wantError: "too many decimal places",
		},
		{
			name: "missing token ID",
			build: func() (*TransactionRequest, error) {
				return NewTransactionBuilder().
					WithUser(123).
					WithWallet(456).
					WithAmount("100").
					WithFundType(FundTypeDeposit).
					Build()
			},
			wantError: "token ID must be positive",
		},
		{
			name: "invalid fund type",
			build: func() (*TransactionRequest, error) {
				return NewTransactionBuilder().
					WithUser(123).
					WithWallet(456).
					WithAmount("100").
					WithToken(1).
					WithFundType("invalid_type").
					Build()
			},
			wantError: "invalid fund type",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.build()
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantError)
		})
	}
}

func TestTransactionBuilder_AmountEdgeCases(t *testing.T) {
	tests := []struct {
		name   string
		amount string
		valid  bool
	}{
		{"minimum precision", "0.000000000000000001", true},
		{"below minimum", "0.0000000000000000001", false},
		{"maximum amount", "1000000", true},
		{"just below maximum", "999999.999999999999999999", true},
		{"with trailing zeros", "100.5000", true},
		{"scientific notation", "1e10", false},
		{"with spaces", "  100.50  ", true}, // should be trimmed
		{"multiple dots", "100.50.25", false},
		{"leading zeros", "00100.50", true},
		{"only decimal", ".50", false},
		{"ending with dot", "100.", true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewTransactionBuilder().
				WithUser(123).
				WithWallet(456).
				WithAmount(tt.amount).
				WithToken(1).
				WithFundType(FundTypeDeposit)
			
			_, err := builder.Build()
			if tt.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestTransactionBuilder_AutoGeneration(t *testing.T) {
	builder := NewTransactionBuilder()
	
	req, err := builder.
		WithUser(123).
		WithWallet(456).
		WithAmount("100").
		WithToken(1).
		WithFundType(FundTypeDeposit).
		Build()
	
	require.NoError(t, err)
	
	// Check auto-generated reference
	assert.NotEmpty(t, req.Reference)
	assert.True(t, strings.HasPrefix(req.Reference, string(FundTypeDeposit)))
	assert.Contains(t, req.Reference, "123") // user ID
	
	// Check auto-generated description
	assert.NotEmpty(t, req.Description)
}

func TestTransactionBuilder_ThreadSafety(t *testing.T) {
	builder := NewTransactionBuilder()
	
	// Set initial values
	builder.
		WithUser(1).
		WithWallet(1).
		WithToken(1).
		WithFundType(FundTypeDeposit)
	
	var wg sync.WaitGroup
	errors := make([]error, 100)
	
	// Concurrent modifications
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			
			// Try to modify the builder concurrently
			builder.
				WithUser(int64(idx)).
				WithAmount(fmt.Sprintf("%d.50", idx)).
				WithMetadata(fmt.Sprintf("key%d", idx), idx)
			
			// Build should not panic
			_, err := builder.Build()
			errors[idx] = err
		}(i)
	}
	
	wg.Wait()
	
	// At least some builds should succeed
	successCount := 0
	for _, err := range errors {
		if err == nil {
			successCount++
		}
	}
	assert.Greater(t, successCount, 0)
}

func TestTransactionBuilder_Clone(t *testing.T) {
	original := NewTransactionBuilder().
		WithUser(123).
		WithWallet(456).
		WithAmount("100").
		WithToken(1).
		WithFundType(FundTypeDeposit).
		WithMetadata("key1", "value1")
	
	// Clone the builder
	cloned := original.Clone()
	
	// Modify the clone
	cloned.WithAmount("200").WithMetadata("key2", "value2")
	
	// Build both
	origReq, err := original.Build()
	require.NoError(t, err)
	
	clonedReq, err := cloned.Build()
	require.NoError(t, err)
	
	// Original should not be affected
	assert.Equal(t, "100", origReq.Amount)
	assert.Equal(t, "value1", origReq.Metadata["key1"])
	assert.Nil(t, origReq.Metadata["key2"])
	
	// Clone should have new values
	assert.Equal(t, "200", clonedReq.Amount)
	assert.Equal(t, "value1", clonedReq.Metadata["key1"])
	assert.Equal(t, "value2", clonedReq.Metadata["key2"])
}

func TestTransactionBuilder_Reset(t *testing.T) {
	builder := NewTransactionBuilder().
		WithUser(123).
		WithWallet(456).
		WithAmount("100").
		WithToken(1).
		WithFundType(FundTypeDeposit)
	
	// Reset the builder
	builder.Reset()
	
	// Should fail validation due to missing fields
	_, err := builder.Build()
	assert.Error(t, err)
	
	// Can build new transaction after reset
	req, err := builder.
		WithUser(789).
		WithWallet(101).
		WithAmount("50").
		WithToken(2).
		WithFundType(FundTypeWithdraw).
		Build()
	
	require.NoError(t, err)
	assert.Equal(t, int64(789), req.UserID)
	assert.Equal(t, "50", req.Amount)
	assert.Equal(t, FundTypeWithdraw, req.FundType)
}

func TestTransactionBuilder_SpecificBuilders(t *testing.T) {
	t.Run("RedPacketCreate", func(t *testing.T) {
		req, err := BuildRedPacketCreateTransaction(123, 456, "100", 1, 789)
		require.NoError(t, err)
		assert.Equal(t, FundTypeRedPacketCreate, req.FundType)
		assert.Equal(t, int64(789), req.RelatedID)
		assert.Equal(t, int64(789), req.Metadata["red_packet_id"])
	})
	
	t.Run("TransferOut", func(t *testing.T) {
		req, err := BuildTransferOutTransaction(123, 456, "50", 1, 999, "recipient_user")
		require.NoError(t, err)
		assert.Equal(t, FundTypeTransferOut, req.FundType)
		assert.Equal(t, int64(999), req.RelatedID)
		assert.Equal(t, int64(999), req.Metadata["transfer_id"])
		assert.Equal(t, "recipient_user", req.Metadata["recipient"])
	})
	
	t.Run("Deposit", func(t *testing.T) {
		req, err := BuildDepositTransaction(123, 456, "200", 1, "0x1234567890")
		require.NoError(t, err)
		assert.Equal(t, FundTypeDeposit, req.FundType)
		assert.Equal(t, "0x1234567890", req.Reference)
		assert.Equal(t, "0x1234567890", req.Metadata["tx_hash"])
	})
	
	t.Run("Withdraw", func(t *testing.T) {
		req, err := BuildWithdrawTransaction(123, 456, "150", 1, 888, "0xabcdef")
		require.NoError(t, err)
		assert.Equal(t, FundTypeWithdraw, req.FundType)
		assert.Equal(t, int64(888), req.RelatedID)
		assert.Equal(t, int64(888), req.Metadata["withdraw_id"])
		assert.Equal(t, "0xabcdef", req.Metadata["address"])
	})
}

func TestTransactionBuilder_MetadataSafety(t *testing.T) {
	builder := NewTransactionBuilder().
		WithUser(123).
		WithWallet(456).
		WithAmount("100").
		WithToken(1).
		WithFundType(FundTypeDeposit).
		WithMetadata("key1", "value1")
	
	// Build first request
	req1, err := builder.Build()
	require.NoError(t, err)
	
	// Modify metadata after build
	req1.Metadata["key2"] = "value2"
	
	// Build second request
	req2, err := builder.Build()
	require.NoError(t, err)
	
	// Second request should not have the modification
	assert.Equal(t, "value1", req2.Metadata["key1"])
	assert.Nil(t, req2.Metadata["key2"])
}

func TestTransactionBuilder_EmptyMetadataKey(t *testing.T) {
	builder := NewTransactionBuilder().
		WithUser(123).
		WithWallet(456).
		WithAmount("100").
		WithToken(1).
		WithFundType(FundTypeDeposit).
		WithMetadata("", "should be ignored")
	
	req, err := builder.Build()
	require.NoError(t, err)
	
	// Empty key should be ignored
	assert.NotContains(t, req.Metadata, "")
}

// Benchmark tests
func BenchmarkTransactionBuilder_Build(b *testing.B) {
	builder := NewTransactionBuilder()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builder.
			WithUser(123).
			WithWallet(456).
			WithAmount("100.50").
			WithToken(1).
			WithFundType(FundTypeDeposit).
			WithReference("test_ref").
			WithDescription("Test transaction").
			WithMetadata("key1", "value1").
			WithMetadata("key2", i).
			Build()
		
		builder.Reset()
	}
}

func BenchmarkTransactionBuilder_ConcurrentBuild(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		builder := NewTransactionBuilder()
		i := 0
		for pb.Next() {
			builder.
				WithUser(int64(i)).
				WithWallet(int64(i+1)).
				WithAmount(fmt.Sprintf("%d.50", i)).
				WithToken(1).
				WithFundType(FundTypeDeposit).
				Build()
			
			builder.Reset()
			i++
		}
	})
}