// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Wallets is the golang structure for table wallets.
type Wallets struct {
	WalletId         int         `json:"walletId"         orm:"wallet_id"         description:"钱包记录唯一 ID"`                    // 钱包记录唯一 ID
	LedgerId         string      `json:"ledgerId"         orm:"ledger_id"         description:"ledger 钱包id"`                  // ledger 钱包id
	UserId           uint        `json:"userId"           orm:"user_id"           description:"用户 ID (外键, 关联 users.user_id)"` // 用户 ID (外键, 关联 users.user_id)
	TokenId          uint        `json:"tokenId"          orm:"token_id"          description:"用户 ID (外键, 关联 users.user_id)"` // 用户 ID (外键, 关联 users.user_id)
	AvailableBalance int64       `json:"availableBalance" orm:"available_balance" description:"可用余额"`                         // 可用余额
	FrozenBalance    int64       `json:"frozenBalance"    orm:"frozen_balance"    description:"冻结余额 (例如: 挂单中, 提现处理中)"`        // 冻结余额 (例如: 挂单中, 提现处理中)
	DecimalPlaces    uint        `json:"decimalPlaces"    orm:"decimal_places"    description:"精度"`                           // 精度
	CreatedAt        *gtime.Time `json:"createdAt"        orm:"created_at"        description:"钱包记录创建时间"`                     // 钱包记录创建时间
	UpdatedAt        *gtime.Time `json:"updatedAt"        orm:"updated_at"        description:"余额最后更新时间"`                     // 余额最后更新时间
	DeletedAt        *gtime.Time `json:"deletedAt"        orm:"deleted_at"        description:"软删除的时间戳"`                      // 软删除的时间戳
	TelegramId       int64       `json:"telegramId"       orm:"telegram_id"       description:""`                             //
	Type             string      `json:"type"             orm:"type"              description:"类型"`                           // 类型
	Symbol           string      `json:"symbol"           orm:"symbol"            description:"代币符号 (例如: USDT, BTC, ETH)"`    // 代币符号 (例如: USDT, BTC, ETH)
}
