package dao

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/yalks/wallet/entity"
)

// IWalletDAO 钱包数据访问接口
type IWalletDAO interface {
	// GetWalletByUserIDAndSymbol 通过用户ID和代币符号获取钱包
	GetWalletByUserIDAndSymbol(ctx context.Context, userID uint64, symbol string) (*entity.Wallets, error)
	// CreateWalletRecord 创建钱包记录
	CreateWalletRecord(ctx context.Context, tx gdb.TX, wallet *entity.Wallets) error
	// UpdateWalletBalance 更新钱包余额
	UpdateWalletBalance(ctx context.Context, tx gdb.TX, walletID uint, availableBalance, frozenBalance int64) error
	// GetAllWallets 获取所有钱包记录
	GetAllWallets(ctx context.Context) ([]*entity.Wallets, error)
}

type walletDAO struct{}

// NewWalletDAO 创建钱包DAO实例
func NewWalletDAO() IWalletDAO {
	return &walletDAO{}
}

// GetWalletByUserIDAndSymbol 通过用户ID和代币符号获取钱包
func (d *walletDAO) GetWalletByUserIDAndSymbol(ctx context.Context, userID uint64, symbol string) (*entity.Wallets, error) {
	var wallet *entity.Wallets
	err := g.Model("wallets").Ctx(ctx).
		Where("user_id = ? AND symbol = ? AND deleted_at IS NULL", userID, symbol).
		Scan(&wallet)
	if err != nil {
		return nil, gerror.Wrapf(err, "查询钱包失败: UserID=%d, Symbol=%s", userID, symbol)
	}
	return wallet, nil
}

// CreateWalletRecord 创建钱包记录
func (d *walletDAO) CreateWalletRecord(ctx context.Context, tx gdb.TX, wallet *entity.Wallets) error {
	var db *gdb.Model
	if tx != nil {
		db = g.Model("wallets").Ctx(ctx).TX(tx)
	} else {
		db = g.Model("wallets").Ctx(ctx)
	}

	_, err := db.Insert(wallet)
	if err != nil {
		return gerror.Wrapf(err, "创建钱包记录失败: UserID=%d, Symbol=%s", wallet.UserId, wallet.Symbol)
	}
	return nil
}

// UpdateWalletBalance 更新钱包余额
func (d *walletDAO) UpdateWalletBalance(ctx context.Context, tx gdb.TX, walletID uint, availableBalance, frozenBalance int64) error {
	var db *gdb.Model
	if tx != nil {
		db = g.Model("wallets").Ctx(ctx).TX(tx)
	} else {
		db = g.Model("wallets").Ctx(ctx)
	}

	_, err := db.Where("wallet_id = ?", walletID).Update(map[string]any{
		"available_balance": availableBalance,
		"frozen_balance":    frozenBalance,
		"updated_at":        gtime.Now(),
	})
	if err != nil {
		return gerror.Wrapf(err, "更新钱包余额失败: WalletID=%d", walletID)
	}
	return nil
}

// GetAllWallets 获取所有钱包记录
func (d *walletDAO) GetAllWallets(ctx context.Context) ([]*entity.Wallets, error) {
	var wallets []*entity.Wallets
	err := g.Model("wallets").Ctx(ctx).
		Where("deleted_at IS NULL").
		Scan(&wallets)
	if err != nil {
		return nil, gerror.Wrap(err, "查询所有钱包记录失败")
	}
	return wallets, nil
}
