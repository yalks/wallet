// Package dao provides the Data Access Object layer for the wallet module.
//
// This package contains interfaces and implementations for database operations
// related to users, wallets, tokens, and transactions. All database operations
// are abstracted through well-defined interfaces to ensure testability and
// maintainability.
//
// # Key Interfaces
//
//   - IUserDAO: User data operations
//   - IWalletDAO: Wallet management operations
//   - ITokenDAO: Token/currency operations
//   - ITransactionDAO: Transaction record operations
//
// # Transaction Support
//
// All DAO methods support database transactions through the gdb.TX parameter,
// allowing for atomic operations across multiple tables:
//
//	tx, err := g.DB().Begin(ctx)
//	defer tx.Rollback()
//	
//	// Perform multiple operations
//	err = walletDAO.UpdateBalance(ctx, tx, walletID, amount)
//	err = transactionDAO.CreateTransaction(ctx, tx, transaction)
//	
//	tx.Commit()
//
// # Error Handling
//
// All methods return detailed errors with context about the operation that failed,
// making debugging and error tracking easier in production environments.
package dao