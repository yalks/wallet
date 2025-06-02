# Wallet Constants Package

This package defines all fund types and operation types used in the wallet module.

## Fund Types

All possible fund transaction types are defined in `fund_types.go`:

### Red Packet Operations
- `FundTypeRedPacketCreate` - 红包创建扣费
- `FundTypeRedPacketClaim` - 领取红包
- `FundTypeRedPacketRefund` - 红包过期退回
- `FundTypeRedPacketCancel` - 手动撤回红包

### Transfer Operations
- `FundTypeTransferOut` - 转账发起(扣款)
- `FundTypeTransferIn` - 转账收款(入账)
- `FundTypeTransferExpired` - 转账过期退回

### Payment Operations
- `FundTypePaymentRequest` - 收款请求
- `FundTypePaymentOut` - 付款(扣款)
- `FundTypePaymentIn` - 收款(入账)

### Wallet Operations
- `FundTypeDeposit` - 充值
- `FundTypeWithdraw` - 提现
- `FundTypeWithdrawRefund` - 提现退回

### Admin Operations
- `FundTypeAdminAdd` - 后台增加
- `FundTypeAdminDeduct` - 后台减少

### Exchange Operations
- `FundTypeExchangeOut` - 兑换减少
- `FundTypeExchangeIn` - 兑换增加

### Other Operations
- `FundTypeCommission` - 佣金
- `FundTypeReferralBonus` - 推荐奖励
- `FundTypeSystemAdjustment` - 系统调整

## Usage

```go
import "github.com/yalks/wallet/constants"

// Create a transaction
txReq, err := constants.NewTransactionBuilder().
    WithUser(userID).
    WithWallet(walletID).
    WithAmount("100.50").
    WithToken(tokenID).
    WithFundType(constants.FundTypeRedPacketCreate).
    WithRelatedID(redPacketID).
    Build()

// Check fund type properties
if constants.IsValidFundType(fundType) {
    info, _ := constants.GetFundTypeInfo(fundType)
    direction := constants.GetFundDirection(fundType)
}
```

## Files

- `fund_types.go` - Fund type definitions and registry
- `operation_types.go` - Operation types and wallet types
- `transaction_builder.go` - Transaction builder helper
- `usage_example.go` - Usage examples
- `constants.go` - Package entry point