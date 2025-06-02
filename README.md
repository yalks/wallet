# Wallet Module

[![Go Reference](https://pkg.go.dev/badge/github.com/yalks/wallet.svg)](https://pkg.go.dev/github.com/yalks/wallet)
[![Go Report Card](https://goreportcard.com/badge/github.com/yalks/wallet)](https://goreportcard.com/report/github.com/yalks/wallet)

A comprehensive Go module for managing digital wallets, transactions, and user balances. This module provides a clean interface for wallet operations with support for multiple currencies, transaction management, and balance tracking.

## Features

- **User Management**: Create and manage user accounts with wallet associations
- **Multi-Currency Support**: Handle multiple token types and currencies
- **Transaction Management**: Complete transaction lifecycle management with atomic operations
- **Balance Operations**: Real-time balance tracking and updates
- **Thread-Safe Operations**: Concurrent-safe implementations using proper locking mechanisms
- **Database Abstraction**: Clean DAO layer for database operations
- **Extensible Architecture**: Well-defined interfaces for easy extension

## Installation

```bash
go get github.com/yalks/wallet
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/yalks/wallet"
    "github.com/gogf/gf/v2/frame/g"
)

func main() {
    // Initialize the wallet manager
    manager := wallet.Manager()
    
    // Create a new user
    user, err := manager.CreateUser(ctx, &wallet.CreateUserReq{
        Uid:        "user123",
        BusinessId: "business456",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Get user balance
    balance, err := manager.GetBalance(ctx, &wallet.GetBalanceReq{
        WalletId: user.WalletId,
        TokenId:  1, // USDT
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("User balance: %s\n", balance.Amount)
}
```

## Architecture

The module follows a layered architecture:

```
├── interfaces.go       # Public interfaces (IWalletManager, ITransactionManager)
├── manager.go          # Main wallet manager implementation
├── transaction_manager.go # Transaction management
├── singleton.go        # Singleton pattern for manager instance
├── dao/                # Data Access Object layer
│   ├── token.go        # Token DAO
│   ├── transaction.go  # Transaction DAO
│   ├── user.go         # User DAO
│   └── wallet.go       # Wallet DAO
├── entity/             # Data models
│   ├── tokens.go       # Token entities
│   ├── transactions.go # Transaction entities
│   ├── users.go        # User entities
│   └── wallets.go      # Wallet entities
├── logic/              # Business logic layer
│   ├── balance.go      # Balance calculations
│   ├── context.go      # Context utilities
│   ├── operation.go    # Operation logic
│   ├── token.go        # Token operations
│   └── user.go         # User operations
└── constants/          # Constants and enums
    ├── constants.go    # General constants
    ├── fund_types.go   # Fund type definitions
    ├── operation_types.go # Operation type definitions
    └── transaction_builder.go # Transaction builder utilities
```

## Core Interfaces

### IWalletManager

The main interface for wallet operations:

```go
type IWalletManager interface {
    // User operations
    CreateUser(ctx context.Context, req *CreateUserReq) (*CreateUserResp, error)
    GetBalance(ctx context.Context, req *GetBalanceReq) (*GetBalanceResp, error)
    
    // Transaction operations
    CreateTransaction(ctx context.Context, req *CreateTransactionReq) (*CreateTransactionResp, error)
    ExecuteTransaction(ctx context.Context, req *ExecuteTransactionReq) (*ExecuteTransactionResp, error)
    
    // Query operations
    GetUserInfo(ctx context.Context, req *GetUserInfoReq) (*GetUserInfoResp, error)
    GetTransactionHistory(ctx context.Context, req *GetTransactionHistoryReq) (*GetTransactionHistoryResp, error)
}
```

### ITransactionManager

Interface for advanced transaction management:

```go
type ITransactionManager interface {
    CreateAndExecuteTransaction(ctx context.Context, param CreateAndExecuteTransactionParam) error
    CreatePendingTransaction(ctx context.Context, param CreatePendingTransactionParam) (int64, error)
    ExecutePendingTransaction(ctx context.Context, transactionId int64) error
}
```

## Usage Examples

### Creating a User

```go
user, err := manager.CreateUser(ctx, &wallet.CreateUserReq{
    Uid:        "unique_user_id",
    BusinessId: "business_id",
    UserType:   1, // Regular user
})
```

### Checking Balance

```go
balance, err := manager.GetBalance(ctx, &wallet.GetBalanceReq{
    WalletId: walletId,
    TokenId:  tokenId,
})
```

### Creating a Transaction

```go
tx, err := manager.CreateTransaction(ctx, &wallet.CreateTransactionReq{
    WalletId:      walletId,
    TokenId:       tokenId,
    Amount:        "100.50",
    OperationType: constants.OperationTypeDeposit,
    BusinessId:    "business_id",
})
```

### Executing a Transaction

```go
result, err := manager.ExecuteTransaction(ctx, &wallet.ExecuteTransactionReq{
    TransactionId: tx.TransactionId,
})
```

## Configuration

The module uses GoFrame's configuration system. Database configuration should be set up in your application:

```yaml
database:
  default:
    link: "mysql:user:password@tcp(127.0.0.1:3306)/wallet_db"
```

## Error Handling

All methods return detailed errors with context:

```go
if err != nil {
    // Errors include context about what went wrong
    log.Printf("Failed to create transaction: %v", err)
}
```

## Thread Safety

The wallet manager is thread-safe and can be used concurrently. The singleton pattern ensures a single instance is used throughout the application:

```go
// Get the singleton instance
manager := wallet.Manager()
```

## Database Schema

The module expects the following database tables:
- `user` - User information
- `wallet` - Wallet accounts
- `transaction` - Transaction records
- `token` - Token/currency definitions

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Testing

Run the test suite:

```bash
go test ./...
```

Run with coverage:

```bash
go test -cover ./...
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support and questions:
- Open an issue on GitHub
- Contact the maintainers

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for a list of changes.