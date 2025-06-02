package logic

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	"github.com/yalks/wallet/entity"
)

// IUserLogic 用户业务逻辑接口
type IUserLogic interface {
	// GetUserByID 通过用户ID获取用户
	GetUserByID(ctx context.Context, userID uint64) (*entity.Users, error)
	// GetUserByTelegramID 通过Telegram ID获取用户
	GetUserByTelegramID(ctx context.Context, telegramID int64) (*entity.Users, error)
	// GetUserByUsername 通过用户名获取用户
	GetUserByUsername(ctx context.Context, username string) (*entity.Users, error)
	// GetMainWalletID 获取用户主钱包ID
	GetMainWalletID(ctx context.Context, userID uint64) (string, error)
}

type userLogic struct {
	context *SharedLogicContext
}

// NewUserLogic 创建用户业务逻辑实例
func NewUserLogic() IUserLogic {
	return &userLogic{
		context: GetSharedContext(),
	}
}

// GetUserByID 通过用户ID获取用户
func (l *userLogic) GetUserByID(ctx context.Context, userID uint64) (*entity.Users, error) {
	user, err := l.context.GetUserDAO().GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, gerror.Newf("用户不存在: UserID=%d", userID)
	}
	return user, nil
}

// GetUserByTelegramID 通过Telegram ID获取用户
func (l *userLogic) GetUserByTelegramID(ctx context.Context, telegramID int64) (*entity.Users, error) {
	user, err := l.context.GetUserDAO().GetUserByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, gerror.Newf("用户不存在: TelegramID=%d", telegramID)
	}
	return user, nil
}

// GetUserByUsername 通过用户名获取用户
func (l *userLogic) GetUserByUsername(ctx context.Context, username string) (*entity.Users, error) {
	user, err := l.context.GetUserDAO().GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, gerror.Newf("用户不存在: Username=%s", username)
	}
	return user, nil
}

// GetMainWalletID 获取用户主钱包ID
func (l *userLogic) GetMainWalletID(ctx context.Context, userID uint64) (string, error) {
	user, err := l.GetUserByID(ctx, userID)
	if err != nil {
		return "", err
	}
	if user.MainWalletId == "" {
		return "", gerror.Newf("用户没有主钱包: UserID=%d", userID)
	}
	return user.MainWalletId, nil
}
