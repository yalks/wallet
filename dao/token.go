package dao

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"github.com/yalks/wallet/entity"
)

// ITokenDAO 代币数据访问接口
type ITokenDAO interface {
	// GetTokenBySymbol 通过符号获取代币
	GetTokenBySymbol(ctx context.Context, symbol string) (*entity.Tokens, error)
	// GetTokenByID 通过ID获取代币
	GetTokenByID(ctx context.Context, tokenID uint) (*entity.Tokens, error)
	// GetActiveTokens 获取所有活跃代币
	GetActiveTokens(ctx context.Context) ([]*entity.Tokens, error)
}

type tokenDAO struct{}

// NewTokenDAO 创建代币DAO实例
func NewTokenDAO() ITokenDAO {
	return &tokenDAO{}
}

// GetTokenBySymbol 通过符号获取代币
func (d *tokenDAO) GetTokenBySymbol(ctx context.Context, symbol string) (*entity.Tokens, error) {
	var token *entity.Tokens
	err := g.Model("tokens").Ctx(ctx).
		Where("symbol = ? AND is_active = 1 AND status = 1", symbol).
		OrderAsc("token_id").
		Scan(&token)
	if err != nil {
		return nil, gerror.Wrapf(err, "查询代币失败: Symbol=%s", symbol)
	}
	if token != nil { // 确保 token 不是 nil 才打印，避免 panic
		g.Log().Debugf(ctx, "DAO GetTokenBySymbol for symbol '%s': Fetched Token ID: %d, Symbol: %s, Network: %s, Decimals: %d, IsActive: %d, Status: %d",
			symbol, token.TokenId, token.Symbol, token.Network, token.Decimals, token.IsActive, token.Status)
	}
	return token, nil
}

// GetTokenByID 通过ID获取代币
func (d *tokenDAO) GetTokenByID(ctx context.Context, tokenID uint) (*entity.Tokens, error) {
	var token *entity.Tokens
	err := g.Model("tokens").Ctx(ctx).
		Where("token_id = ? AND is_active = 1 AND status = 1", tokenID).
		Scan(&token)
	if err != nil {
		return nil, gerror.Wrapf(err, "查询代币失败: TokenID=%d", tokenID)
	}
	return token, nil
}

// GetActiveTokens 获取所有活跃代币
func (d *tokenDAO) GetActiveTokens(ctx context.Context) ([]*entity.Tokens, error) {
	var tokens []*entity.Tokens
	err := g.Model("tokens").Ctx(ctx).
		Where("is_active = 1 AND status = 1").
		OrderAsc("token_id").
		Scan(&tokens)
	if err != nil {
		return nil, gerror.Wrap(err, "查询活跃代币失败")
	}
	return tokens, nil
}
