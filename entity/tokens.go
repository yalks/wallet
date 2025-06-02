// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Tokens is the golang structure for table tokens.
type Tokens struct {
	TokenId                 uint        `json:"tokenId"                 orm:"token_id"                    description:"代币内部 ID (主键)"`                                      // 代币内部 ID (主键)
	Symbol                  string      `json:"symbol"                  orm:"symbol"                      description:"代币符号 (例如: USDT, BTC, ETH)"`                         // 代币符号 (例如: USDT, BTC, ETH)
	Name                    string      `json:"name"                    orm:"name"                        description:"代币名称 (例如: Tether, Bitcoin, Ethereum)"`              // 代币名称 (例如: Tether, Bitcoin, Ethereum)
	Network                 string      `json:"network"                 orm:"network"                     description:"所属网络/链 (例如: Ethereum, Tron, Bitcoin, BSC, Solana)"` // 所属网络/链 (例如: Ethereum, Tron, Bitcoin, BSC, Solana)
	ChainId                 uint        `json:"chainId"                 orm:"chain_id"                    description:"链 ID (主要用于 EVM 兼容链)"`                               // 链 ID (主要用于 EVM 兼容链)
	Confirmations           uint        `json:"confirmations"           orm:"confirmations"               description:"确认次数"`                                              // 确认次数
	ContractAddress         string      `json:"contractAddress"         orm:"contract_address"            description:"代币合约地址 (对于非原生代币). 对于原生代币为 NULL."`                   // 代币合约地址 (对于非原生代币). 对于原生代币为 NULL.
	TokenStandard           string      `json:"tokenStandard"           orm:"token_standard"              description:"代币标准 (例如: native, ERC20, TRC20, BEP20, SPL)"`       // 代币标准 (例如: native, ERC20, TRC20, BEP20, SPL)
	Decimals                uint        `json:"decimals"                orm:"decimals"                    description:"代币精度/小数位数"`                                         // 代币精度/小数位数
	IsStablecoin            int         `json:"isStablecoin"            orm:"is_stablecoin"               description:"是否为稳定币"`                                            // 是否为稳定币
	IsFiat                  int         `json:"isFiat"                  orm:"is_fiat"                     description:"是否为法币"`                                             // 是否为法币
	IsActive                int         `json:"isActive"                orm:"is_active"                   description:"代币是否在系统内激活可用 (总开关)"`                                // 代币是否在系统内激活可用 (总开关)
	AllowDeposit            int         `json:"allowDeposit"            orm:"allow_deposit"               description:"是否允许充值 (设为 FALSE 即表示充值维护中)"`                        // 是否允许充值 (设为 FALSE 即表示充值维护中)
	AllowWithdraw           int         `json:"allowWithdraw"           orm:"allow_withdraw"              description:"是否允许提现 (设为 FALSE 即表示提现维护中)"`                        // 是否允许提现 (设为 FALSE 即表示提现维护中)
	AllowTransfer           int         `json:"allowTransfer"           orm:"allow_transfer"              description:"是否允许内部转账"`                                          // 是否允许内部转账
	AllowReceive            int         `json:"allowReceive"            orm:"allow_receive"               description:"是否允许内部收款 (通常与 allow_transfer 联动或独立控制)"`             // 是否允许内部收款 (通常与 allow_transfer 联动或独立控制)
	AllowRedPacket          int         `json:"allowRedPacket"          orm:"allow_red_packet"            description:"是否允许发红包"`                                           // 是否允许发红包
	AllowTrading            int         `json:"allowTrading"            orm:"allow_trading"               description:"是否允许在交易对中使用"`                                       // 是否允许在交易对中使用
	MaintenanceMessage      string      `json:"maintenanceMessage"      orm:"maintenance_message"         description:"维护信息 (当充值/提现/其他功能禁用时显示)"`                           // 维护信息 (当充值/提现/其他功能禁用时显示)
	WithdrawalFeeType       string      `json:"withdrawalFeeType"       orm:"withdrawal_fee_type"         description:"提币手续费类型: fixed-固定金额, percent-百分比"`                  // 提币手续费类型: fixed-固定金额, percent-百分比
	WithdrawalFeeAmount     string      `json:"withdrawalFeeAmount"     orm:"withdrawal_fee_amount"       description:"提币手续费金额/比例 (根据 fee_type 决定)"`                       // 提币手续费金额/比例 (根据 fee_type 决定)
	MinDepositAmount        string      `json:"minDepositAmount"        orm:"min_deposit_amount"          description:"单笔最小充值金额"`                                          // 单笔最小充值金额
	MaxDepositAmount        string      `json:"maxDepositAmount"        orm:"max_deposit_amount"          description:"单笔最大充值金额 (NULL 表示无限制)"`                             // 单笔最大充值金额 (NULL 表示无限制)
	MinWithdrawalAmount     string      `json:"minWithdrawalAmount"     orm:"min_withdrawal_amount"       description:"单笔最小提币金额"`                                          // 单笔最小提币金额
	MaxWithdrawalAmount     string      `json:"maxWithdrawalAmount"     orm:"max_withdrawal_amount"       description:"单笔最大提币金额 (NULL 表示无限制)"`                             // 单笔最大提币金额 (NULL 表示无限制)
	MinTransferAmount       string      `json:"minTransferAmount"       orm:"min_transfer_amount"         description:"单笔最小转账金额"`                                          // 单笔最小转账金额
	MaxTransferAmount       string      `json:"maxTransferAmount"       orm:"max_transfer_amount"         description:"单笔最大转账金额 (NULL 表示无限制)"`                             // 单笔最大转账金额 (NULL 表示无限制)
	MinReceiveAmount        string      `json:"minReceiveAmount"        orm:"min_receive_amount"          description:"单笔最小收款金额"`                                          // 单笔最小收款金额
	MaxReceiveAmount        string      `json:"maxReceiveAmount"        orm:"max_receive_amount"          description:"单笔最大收款金额 (NULL 表示无限制)"`                             // 单笔最大收款金额 (NULL 表示无限制)
	MinRedPacketAmount      string      `json:"minRedPacketAmount"      orm:"min_red_packet_amount"       description:"单个红包最小金额"`                                          // 单个红包最小金额
	MaxRedPacketAmount      string      `json:"maxRedPacketAmount"      orm:"max_red_packet_amount"       description:"单个红包最大金额 (NULL 表示无限制)"`                             // 单个红包最大金额 (NULL 表示无限制)
	MaxRedPacketCount       int64       `json:"maxRedPacketCount"       orm:"max_red_packet_count"        description:"单次发放红包最大个数 (NULL 表示无限制)"`                           // 单次发放红包最大个数 (NULL 表示无限制)
	MaxRedPacketTotalAmount string      `json:"maxRedPacketTotalAmount" orm:"max_red_packet_total_amount" description:"单次发放红包最大总金额 (NULL 表示无限制)"`                          // 单次发放红包最大总金额 (NULL 表示无限制)
	LogoUrl                 string      `json:"logoUrl"                 orm:"logo_url"                    description:"代币 Logo 图片 URL"`                                    // 代币 Logo 图片 URL
	ProjectWebsite          string      `json:"projectWebsite"          orm:"project_website"             description:"项目官方网站 URL"`                                        // 项目官方网站 URL
	Description             string      `json:"description"             orm:"description"                 description:"代币或项目描述"`                                           // 代币或项目描述
	Order                   uint        `json:"order"                   orm:"order"                       description:"排序字段 (用于前端展示时的排序)"`                                 // 排序字段 (用于前端展示时的排序)
	Status                  int         `json:"status"                  orm:"status"                      description:"状态 0-下架 1-上架"`                                      // 状态 0-下架 1-上架
	CreatedAt               *gtime.Time `json:"createdAt"               orm:"created_at"                  description:"创建时间"`                                              // 创建时间
	UpdatedAt               *gtime.Time `json:"updatedAt"               orm:"updated_at"                  description:"最后更新时间"`                                            // 最后更新时间
	DeletedAt               *gtime.Time `json:"deletedAt"               orm:"deleted_at"                  description:"软删除时间"`                                             // 软删除时间
}
