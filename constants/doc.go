// Package constants defines all constants, enumerations, and helper utilities
// used throughout the wallet module.
//
// This package centralizes all constant definitions to ensure consistency
// and prevent magic numbers/strings in the codebase. It includes operation
// types, fund types, status codes, and transaction builders.
//
// # Operation Types
//
// Defines the types of operations that can be performed on wallets:
//   - OperationTypeDeposit: Adding funds to a wallet
//   - OperationTypeWithdraw: Removing funds from a wallet
//   - OperationTypeTransfer: Moving funds between wallets
//
// # Fund Types
//
// Categorizes different types of fund movements for accounting purposes:
//   - FundTypeRedPacket: Red packet related operations
//   - FundTypeGame: Gaming related transactions
//   - FundTypeCommission: Commission payments
//   - And many more...
//
// # Transaction Builder
//
// Provides a fluent interface for building transaction requests:
//
//	req, err := constants.NewTransactionBuilder().
//	    WithUser(userID).
//	    WithAmount("100.50").
//	    WithFundType(constants.FundTypeDeposit).
//	    Build()
//
// # Usage Guidelines
//
// Always use constants from this package instead of hardcoding values.
// This ensures consistency and makes refactoring easier.
package constants