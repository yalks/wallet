package dao

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"github.com/yalks/wallet/entity"
)

// IUserDAO 用户数据访问接口
type IUserDAO interface {
	// GetUserByID 通过用户ID获取用户
	GetUserByID(ctx context.Context, userID uint64) (*entity.Users, error)
	// GetUserByTelegramID 通过Telegram ID获取用户
	GetUserByTelegramID(ctx context.Context, telegramID int64) (*entity.Users, error)
	// GetUserByUsername 通过用户名获取用户
	GetUserByUsername(ctx context.Context, username string) (*entity.Users, error)
}

type userDAO struct{}

// NewUserDAO 创建用户DAO实例
func NewUserDAO() IUserDAO {
	return &userDAO{}
}

// GetUserByID 通过用户ID获取用户
func (d *userDAO) GetUserByID(ctx context.Context, userID uint64) (*entity.Users, error) {
	var user *entity.Users
	err := g.Model("users").Ctx(ctx).Where("id = ? AND deleted_at IS NULL", userID).Scan(&user)
	if err != nil {
		return nil, gerror.Wrapf(err, "查询用户失败: UserID=%d", userID)
	}
	return user, nil
}

// GetUserByTelegramID 通过Telegram ID获取用户
func (d *userDAO) GetUserByTelegramID(ctx context.Context, telegramID int64) (*entity.Users, error) {
	var user *entity.Users
	err := g.Model("users").Ctx(ctx).Where("telegram_id = ? AND deleted_at IS NULL", telegramID).Scan(&user)
	if err != nil {
		return nil, gerror.Wrapf(err, "查询用户失败: TelegramID=%d", telegramID)
	}
	return user, nil
}

// GetUserByUsername 通过用户名获取用户
func (d *userDAO) GetUserByUsername(ctx context.Context, username string) (*entity.Users, error) {
	var user *entity.Users
	err := g.Model("users").Ctx(ctx).Where("username = ? AND deleted_at IS NULL", username).Scan(&user)
	if err != nil {
		return nil, gerror.Wrapf(err, "查询用户失败: Username=%s", username)
	}
	return user, nil
}
