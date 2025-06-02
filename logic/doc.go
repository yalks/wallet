// Package logic implements the business logic layer for the wallet module.
//
// This package contains the core business rules, validations, and operations
// that power the wallet system. It acts as an intermediary between the DAO
// layer and the public interfaces, ensuring data consistency and business
// rule enforcement.
//
// # Key Components
//
//   - IUserLogic: User-related business operations
//   - ITokenLogic: Token/currency management logic
//   - IBalanceLogic: Balance calculation and verification
//   - IOperationLogic: Transaction operation processing
//
// # Shared Context
//
// The package uses SharedLogicContext to minimize redundant initializations
// and improve performance by reusing DAO instances and SDK clients:
//
//	ctx := &SharedLogicContext{
//	    userDAO:   dao.NewUserDAO(),
//	    tokenDAO:  dao.NewTokenDAO(),
//	    // ... other DAOs
//	}
//
// # Transaction Safety
//
// All operations that modify wallet state are designed to be transaction-safe,
// with proper rollback mechanisms in case of failures. Critical operations
// use database transactions to ensure atomicity.
//
// # Validation
//
// The logic layer performs comprehensive validation including:
//   - Amount validation (positive values, decimal precision)
//   - Business ID format validation
//   - Balance sufficiency checks
//   - Token existence verification
package logic