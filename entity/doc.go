// Package entity defines the data models and structures used throughout the wallet module.
//
// This package contains the core data structures that represent database tables
// and business entities. All structures use appropriate data types and include
// JSON tags for API serialization.
//
// # Core Entities
//
//   - Users: User account information
//   - Wallets: User wallet records with balance information
//   - Tokens: Supported tokens/currencies
//   - Transactions: Transaction history and records
//
// # Data Types
//
// The package uses specific data types for financial operations:
//   - decimal.Decimal for precise monetary calculations
//   - uint64 for IDs to support large-scale systems
//   - gtime.Time for timestamp fields with timezone support
//
// # JSON Serialization
//
// All entities include JSON tags for API responses. Sensitive fields
// use appropriate tags to control visibility in API responses.
package entity