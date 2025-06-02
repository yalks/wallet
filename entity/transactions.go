// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/shopspring/decimal"
)

// Transactions is the golang structure for table transactions.
type Transactions struct {
	TransactionId        uint64          `json:"transactionId"        orm:"transaction_id"         description:"交易记录 ID (主键)"`                                                                              // 交易记录 ID (主键)
	UserId               uint            `json:"userId"               orm:"user_id"                description:"关联用户 ID"`                                                                                   // 关联用户 ID
	Username             uint            `json:"username"             orm:"username"               description:"telegra username"`                                                                          // telegra username
	TokenId              uint            `json:"tokenId"              orm:"token_id"               description:"关联代币 ID"`                                                                                   // 关联代币 ID
	Type                 string          `json:"type"                 orm:"type"                   description:"交易类型: deposit, withdrawal, transfer, red_packet, payment, commission, system_adjust, etc."` // 交易类型: deposit, withdrawal, transfer, red_packet, payment, commission, system_adjust, etc.
	WalletType           string          `json:"walletType"           orm:"wallet_type"            description:"钱包类型:  冻结frozen ，余额 available"`                                                             // 钱包类型:  冻结frozen ，余额 available
	Direction            string          `json:"direction"            orm:"direction"              description:"资金方向: in-增加, out-减少"`                                                                       // 资金方向: in-增加, out-减少
	Amount               decimal.Decimal `json:"amount"               orm:"amount"                 description:"交易金额 (绝对值)"`                                                                                // 交易金额 (绝对值)
	BalanceBefore        decimal.Decimal `json:"balanceBefore"        orm:"balance_before"         description:"交易前余额快照 (对应钱包)"`                                                                            // 交易前余额快照 (对应钱包)
	BalanceAfter         decimal.Decimal `json:"balanceAfter"         orm:"balance_after"          description:"交易后余额快照 (对应钱包)"`                                                                            // 交易后余额快照 (对应钱包)
	RelatedTransactionId uint64          `json:"relatedTransactionId" orm:"related_transaction_id" description:"关联交易 ID (例如: 转账的对方记录)"`                                                                     // 关联交易 ID (例如: 转账的对方记录)
	RelatedEntityId      uint64          `json:"relatedEntityId"      orm:"related_entity_id"      description:"关联实体 ID (例如: 红包 ID, 提现订单 ID)"`                                                              // 关联实体 ID (例如: 红包 ID, 提现订单 ID)
	RelatedEntityType    string          `json:"relatedEntityType"    orm:"related_entity_type"    description:"关联实体类型 (例如: red_packet, withdrawal_order)"`                                                 // 关联实体类型 (例如: red_packet, withdrawal_order)
	Status               uint            `json:"status"               orm:"status"                 description:"交易状态: 1-成功, 0-失败"`                                                                          // 交易状态: 1-成功, 0-失败
	Memo                 string          `json:"memo"                 orm:"memo"                   description:"交易备注/消息 (例如: 管理员调账原因)"`                                                                     // 交易备注/消息 (例如: 管理员调账原因)
	CreatedAt            *gtime.Time     `json:"createdAt"            orm:"created_at"             description:"创建时间"`                                                                                      // 创建时间
	UpdatedAt            *gtime.Time     `json:"updatedAt"            orm:"updated_at"             description:"最后更新时间"`                                                                                    // 最后更新时间
	DeletedAt            *gtime.Time     `json:"deletedAt"            orm:"deleted_at"             description:"软删除时间"`                                                                                     // 软删除时间
	Symbol               string          `json:"symbol"               orm:"symbol"                 description:"代币符号 (例如: USDT, BTC, ETH)"`                                                                 // 代币符号 (例如: USDT, BTC, ETH)
	BusinessId           string          `json:"businessId"           orm:"business_id"            description:"业务唯一标识符，用于幂等性检查"`                                                                          // 业务唯一标识符，用于幂等性检查
}
