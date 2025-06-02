// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/shopspring/decimal"
)

// Users is the golang structure for table users.
type Users struct {
	Id                             uint64          `json:"id"                             orm:"id"                                description:""`                           //
	TelegramId                     int64           `json:"telegramId"                     orm:"telegram_id"                       description:""`                           //
	Name                           string          `json:"name"                           orm:"name"                              description:""`                           //
	Email                          string          `json:"email"                          orm:"email"                             description:""`                           //
	EmailVerifiedAt                *gtime.Time     `json:"emailVerifiedAt"                orm:"email_verified_at"                 description:""`                           //
	Password                       string          `json:"password"                       orm:"password"                          description:"用户登录密码（经过哈希处理）"`             // 用户登录密码（经过哈希处理）
	RememberToken                  string          `json:"rememberToken"                  orm:"remember_token"                    description:""`                           //
	Account                        string          `json:"account"                        orm:"account"                           description:"唯一用户账户标识符（例如，用于登录）"`         // 唯一用户账户标识符（例如，用于登录）
	AccountType                    int             `json:"accountType"                    orm:"account_type"                      description:"账户类型 1 用户 2商户 3 代理"`         // 账户类型 1 用户 2商户 3 代理
	InviteCode                     string          `json:"inviteCode"                     orm:"invite_code"                       description:"用户唯一邀请码（用于分享给其他人）"`          // 用户唯一邀请码（用于分享给其他人）
	AreaCode                       string          `json:"areaCode"                       orm:"area_code"                         description:"电话区号"`                       // 电话区号
	Phone                          string          `json:"phone"                          orm:"phone"                             description:"用户电话号码"`                     // 用户电话号码
	Avatar                         string          `json:"avatar"                         orm:"avatar"                            description:"用户头像URL或路径"`                 // 用户头像URL或路径
	Nickname                       string          `json:"nickname"                       orm:"nickname"                          description:"用户显示昵称"`                     // 用户显示昵称
	PaymentPassword                string          `json:"paymentPassword"                orm:"payment_password"                  description:"支付密码（经过哈希处理，未设置时为NULL）"`     // 支付密码（经过哈希处理，未设置时为NULL）
	RecommendId                    uint64          `json:"recommendId"                    orm:"recommend_id"                      description:"推荐该用户的用户ID（关联users.id）"`     // 推荐该用户的用户ID（关联users.id）
	Deep                           int             `json:"deep"                           orm:"deep"                              description:"推荐结构中的层级深度"`                 // 推荐结构中的层级深度
	RecommendRelationship          string          `json:"recommendRelationship"          orm:"recommend_relationship"            description:"表示客户推荐层级关系的路径（例如，/1/5/10/）"` // 表示客户推荐层级关系的路径（例如，/1/5/10/）
	AgentRelationship              string          `json:"agentRelationship"              orm:"agent_relationship"                description:"表示代理推荐层级关系的路径"`              // 表示代理推荐层级关系的路径
	FirstId                        uint            `json:"firstId"                        orm:"first_id"                          description:"关联的代理一级ID"`                  // 关联的代理一级ID
	SecondId                       uint            `json:"secondId"                       orm:"second_id"                         description:"关联的代理二级ID"`                  // 关联的代理二级ID
	ThirdId                        uint            `json:"thirdId"                        orm:"third_id"                          description:"关联的代理三级ID"`                  // 关联的代理三级ID
	IsStop                         int             `json:"isStop"                         orm:"is_stop"                           description:"用户账户暂停状态：0=活跃，1=已暂停"`        // 用户账户暂停状态：0=活跃，1=已暂停
	CurrentToken                   string          `json:"currentToken"                   orm:"current_token"                     description:"最后使用的认证令牌（例如，API令牌）"`        // 最后使用的认证令牌（例如，API令牌）
	LastLoginTime                  *gtime.Time     `json:"lastLoginTime"                  orm:"last_login_time"                   description:"最后一次成功登录的时间戳"`               // 最后一次成功登录的时间戳
	Reason                         string          `json:"reason"                         orm:"reason"                            description:"暂停或其他状态变更的原因"`               // 暂停或其他状态变更的原因
	CreatedAt                      *gtime.Time     `json:"createdAt"                      orm:"created_at"                        description:""`                           //
	UpdatedAt                      *gtime.Time     `json:"updatedAt"                      orm:"updated_at"                        description:""`                           //
	DeletedAt                      *gtime.Time     `json:"deletedAt"                      orm:"deleted_at"                        description:"软删除的时间戳"`                    // 软删除的时间戳
	Google2FaSecret                string          `json:"google2FaSecret"                orm:"google2fa_secret"                  description:"谷歌2fa密钥"`                    // 谷歌2fa密钥
	Google2FaEnabled               int             `json:"google2FaEnabled"               orm:"google2fa_enabled"                 description:"谷歌2fa是否启用"`                  // 谷歌2fa是否启用
	IsPaymentPassword              int             `json:"isPaymentPassword"              orm:"is_payment_password"               description:"是否开启免密支付"`                   // 是否开启免密支付
	PaymentPasswordAmount          decimal.Decimal `json:"paymentPasswordAmount"          orm:"payment_password_amount"           description:"免密支付额度"`                     // 免密支付额度
	BackupAccounts                 string          `json:"backupAccounts"                 orm:"backup_accounts"                   description:"备用账户"`                       // 备用账户
	Language                       string          `json:"language"                       orm:"language"                          description:"语言"`                         // 语言
	RedPacketPermission            int             `json:"redPacketPermission"            orm:"red_packet_permission"             description:"红包权限"`                       // 红包权限
	TransferPermission             int             `json:"transferPermission"             orm:"transfer_permission"               description:"转账权限"`                       // 转账权限
	WithdrawPermission             int             `json:"withdrawPermission"             orm:"withdraw_permission"               description:"提现权限"`                       // 提现权限
	FlashTradePermission           int             `json:"flashTradePermission"           orm:"flash_trade_permission"            description:"闪兑权限"`                       // 闪兑权限
	RechargePermission             int             `json:"rechargePermission"             orm:"recharge_permission"               description:"充值权限"`                       // 充值权限
	ReceivePermission              int             `json:"receivePermission"              orm:"receive_permission"                description:"收款权限"`                       // 收款权限
	ResetPaymentPasswordPermission int             `json:"resetPaymentPasswordPermission" orm:"reset_payment_password_permission" description:"重置支付密码权限：0=无权限，1=有权限"`       // 重置支付密码权限：0=无权限，1=有权限
	MainWalletId                   string          `json:"mainWalletId"                   orm:"main_wallet_id"                    description:"钱包id"`                       // 钱包id
}
