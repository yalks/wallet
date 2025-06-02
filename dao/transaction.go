package dao

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/yalks/wallet/entity"
)

// ITransactionDAO 交易数据访问接口
type ITransactionDAO interface {
	// CreateTransaction 创建交易记录
	CreateTransaction(ctx context.Context, tx gdb.TX, transaction *entity.Transactions) (int64, error)
	// GetTransactionByBusinessID 通过业务ID获取交易
	GetTransactionByBusinessID(ctx context.Context, businessID string) (*entity.Transactions, error)
	// UpdateTransactionStatus 更新交易状态
	UpdateTransactionStatus(ctx context.Context, tx gdb.TX, transactionID int64, status string) error
}

type transactionDAO struct{}

// NewTransactionDAO 创建交易DAO实例
func NewTransactionDAO() ITransactionDAO {
	return &transactionDAO{}
}

// CreateTransaction 创建交易记录
func (d *transactionDAO) CreateTransaction(ctx context.Context, tx gdb.TX, transaction *entity.Transactions) (int64, error) {
	var db *gdb.Model
	if tx != nil {
		db = g.Model("transactions").Ctx(ctx).TX(tx)
	} else {
		db = g.Model("transactions").Ctx(ctx)
	}

	result, err := db.Insert(transaction)
	if err != nil {
		return 0, gerror.Wrapf(err, "创建交易记录失败: UserID=%d", transaction.UserId)
	}

	transactionID, err := result.LastInsertId()
	if err != nil {
		return 0, gerror.Wrap(err, "获取交易ID失败")
	}

	return transactionID, nil
}

// GetTransactionByBusinessID 通过业务ID获取交易
func (d *transactionDAO) GetTransactionByBusinessID(ctx context.Context, businessID string) (*entity.Transactions, error) {
	var transaction *entity.Transactions
	err := g.Model("transactions").Ctx(ctx).
		Where("business_id = ? AND status = 1", businessID). // 只查询成功的交易
		OrderDesc("created_at").                             // 按创建时间倒序
		Scan(&transaction)
	if err != nil {
		return nil, gerror.Wrapf(err, "查询交易失败: BusinessID=%s", businessID)
	}
	return transaction, nil
}

// UpdateTransactionStatus 更新交易状态
func (d *transactionDAO) UpdateTransactionStatus(ctx context.Context, tx gdb.TX, transactionID int64, status string) error {
	var db *gdb.Model
	if tx != nil {
		db = g.Model("transactions").Ctx(ctx).TX(tx)
	} else {
		db = g.Model("transactions").Ctx(ctx)
	}

	_, err := db.Where("transaction_id = ?", transactionID).Update(map[string]any{
		"status":     status,
		"updated_at": gtime.Now(),
	})
	if err != nil {
		return gerror.Wrapf(err, "更新交易状态失败: TransactionID=%d", transactionID)
	}
	return nil
}
