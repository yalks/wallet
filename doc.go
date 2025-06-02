// Package wallet provides a comprehensive solution for managing digital wallets,
// transactions, and user balances in Go applications.
//
// The wallet module offers a clean, well-structured API for financial operations
// with support for multiple currencies, atomic transactions, and concurrent-safe
// operations. It follows best practices for financial systems with proper error
// handling, transaction isolation, and balance consistency.
//
// # Architecture
//
// The module is organized into several layers:
//
//   - Public Interface Layer: Exposes IWalletManager and ITransactionManager interfaces
//   - Business Logic Layer: Implements core business rules and validations
//   - Data Access Layer (DAO): Handles database operations with proper transaction support
//   - Entity Layer: Defines data structures and models
//
// # Core Features
//
//   - User account management with unique wallet assignments
//   - Multi-currency/token support with configurable precision
//   - Atomic transaction processing with rollback support
//   - Real-time balance tracking and updates
//   - Thread-safe operations using proper locking mechanisms
//   - Comprehensive error handling with contextual information
//
// # Usage
//
// To use this module, first obtain the singleton wallet manager instance:
//
//	manager := wallet.Manager()
//
// Then use the manager to perform wallet operations:
//
//	// Create a new user
//	user, err := manager.CreateUser(ctx, &wallet.CreateUserReq{
//	    Uid:        "user123",
//	    BusinessId: "business456",
//	})
//
//	// Check balance
//	balance, err := manager.GetBalance(ctx, &wallet.GetBalanceReq{
//	    WalletId: user.WalletId,
//	    TokenId:  1,
//	})
//
// # Transaction Management
//
// For complex operations requiring multiple steps, use the transaction manager:
//
//	txManager := wallet.NewTransactionManager()
//	err := txManager.CreateAndExecuteTransaction(ctx, params)
//
// # Configuration
//
// The module uses GoFrame's configuration system. Ensure your database
// is properly configured in your application's config file.
//
// # Thread Safety
//
// All public methods are thread-safe and can be called concurrently.
// The singleton pattern ensures consistent state across your application.
package wallet