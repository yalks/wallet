package logic

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shopspring/decimal"

	"github.com/yalks/wallet/entity"
)

// ITokenLogic 代币业务逻辑接口
type ITokenLogic interface {
	// GetTokenBySymbol 通过符号获取代币
	GetTokenBySymbol(ctx context.Context, symbol string) (*entity.Tokens, error)
	// GetTokenByID 通过ID获取代币
	GetTokenByID(ctx context.Context, tokenID uint) (*entity.Tokens, error)
	// GetActiveTokens 获取所有活跃代币
	GetActiveTokens(ctx context.Context) ([]*entity.Tokens, error)
	// ConvertBalanceToRaw 将用户余额转换为原始值（乘以10^decimals）
	ConvertBalanceToRaw(ctx context.Context, balance decimal.Decimal, symbol string) (decimal.Decimal, error)
	// ConvertRawToBalance 将原始值转换回用户余额（除以10^decimals）
	ConvertRawToBalance(ctx context.Context, rawBalance decimal.Decimal, symbol string) (decimal.Decimal, error)
	// ConvertBalanceToDBStorage 将余额转换为数据库存储的int64格式
	ConvertBalanceToDBStorage(ctx context.Context, balance decimal.Decimal, symbol string) (int64, error)
	// ConvertDBStorageToBalance 将数据库存储的int64格式转换为余额
	ConvertDBStorageToBalance(ctx context.Context, storageAmount int64, symbol string) (decimal.Decimal, error)
	// FormatUserBalance 格式化用户余额并移除多余小数点和零
	FormatUserBalance(ctx context.Context, balance decimal.Decimal, symbol string) (decimal.Decimal, error)
}

type tokenLogic struct {
	context *SharedLogicContext
}

// NewTokenLogic 创建代币业务逻辑实例
func NewTokenLogic() ITokenLogic {
	return &tokenLogic{
		context: GetSharedContext(),
	}
}

// GetActiveTokens 获取所有活跃代币
func (l *tokenLogic) GetActiveTokens(ctx context.Context) ([]*entity.Tokens, error) {
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

// GetTokenBySymbol 通过符号获取代币
func (l *tokenLogic) GetTokenBySymbol(ctx context.Context, symbol string) (*entity.Tokens, error) {
	token, err := l.context.GetTokenDAO().GetTokenBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, gerror.Newf("代币不存在: Symbol=%s", symbol)
	}
	// Removed duplicate log block, one is sufficient.
	if token != nil { // 确保 token 不是 nil 才打印
		g.Log().Debugf(ctx, "Logic GetTokenBySymbol for symbol '%s': Received Token ID: %d, Symbol: %s, Decimals: %d from DAO",
			symbol, token.TokenId, token.Symbol, token.Decimals)
	}
	return token, nil
}

// GetTokenByID 通过ID获取代币
func (l *tokenLogic) GetTokenByID(ctx context.Context, tokenID uint) (*entity.Tokens, error) {
	token, err := l.context.GetTokenDAO().GetTokenByID(ctx, tokenID)
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, gerror.Newf("代币不存在: TokenID=%d", tokenID)
	}
	return token, nil
}

// ConvertBalanceToRaw 将用户余额转换为原始值（乘以10^decimals）
// 这是系统中唯一应该使用的余额到原始值转换方法
func (l *tokenLogic) ConvertBalanceToRaw(ctx context.Context, balance decimal.Decimal, symbol string) (decimal.Decimal, error) {
	g.Log().Debugf(ctx, "ConvertBalanceToRaw: Called with balance='%s', symbol='%s'", balance.String(), symbol)
	if balance.IsNegative() {
		return decimal.Zero, gerror.Newf("余额不能为负数: %s", balance.String())
	}

	token, err := l.GetTokenBySymbol(ctx, symbol)
	if err != nil {
		return decimal.Zero, gerror.Wrapf(err, "获取代币信息失败: Symbol=%s", symbol)
	}
	g.Log().Debugf(ctx, "ConvertBalanceToRaw: For symbol '%s', got token.Decimals=%d", symbol, token.Decimals)

	// 1. 先将余额截断/舍入到代币本身定义的小数位数
	// 使用 Truncate 来直接截断多余的小数位，确保后续 Shift 后是整数
	adjustedBalance := balance.Truncate(int32(token.Decimals))
	g.Log().Debugf(ctx, "ConvertBalanceToRaw: For symbol '%s', balance='%s' adjusted to '%s' based on token.Decimals=%d",
		symbol, balance.String(), adjustedBalance.String(), token.Decimals)

	// 2. 乘以10^decimals转换为原始值
	rawBalance := adjustedBalance.Shift(int32(token.Decimals))
	g.Log().Debugf(ctx, "ConvertBalanceToRaw: For symbol '%s', adjustedBalance='%s', token.Decimals=%d, Shift result rawBalance='%s'",
		symbol, adjustedBalance.String(), token.Decimals, rawBalance.String())

	// 验证转换结果
	if rawBalance.IsNegative() {
		return decimal.Zero, gerror.Newf("转换后的原始值不能为负数: Balance=%s, Symbol=%s", balance.String(), symbol)
	}
	// 确保结果是整数 (由于先进行了 Truncate，这里应该总是整数了)
	if !rawBalance.IsInteger() {
		// 这一步理论上不应该发生，如果发生了说明逻辑有更深层问题或decimal库行为诡异
		g.Log().Errorf(ctx, "ConvertBalanceToRaw: For symbol '%s', rawBalance '%s' is still not an integer after adjustment and shift. This should not happen if Truncate worked as expected.", symbol, rawBalance.String())
		return decimal.Zero, gerror.Newf("调整和转换后，原始余额仍然不是整数，这不符合预期: %s", rawBalance.String())
	}

	return rawBalance, nil
}

// ConvertRawToBalance 将原始值转换回用户余额（除以10^decimals）
// 这是系统中唯一应该使用的原始值到余额转换方法
func (l *tokenLogic) ConvertRawToBalance(ctx context.Context, rawBalance decimal.Decimal, symbol string) (decimal.Decimal, error) {
	if rawBalance.IsNegative() {
		return decimal.Zero, gerror.Newf("原始余额不能为负数: %s", rawBalance.String())
	}

	token, err := l.GetTokenBySymbol(ctx, symbol)
	if err != nil {
		return decimal.Zero, gerror.Wrapf(err, "获取代币信息失败: Symbol=%s", symbol)
	}

	// 除以10^decimals转换为用户余额
	balance := rawBalance.Shift(-int32(token.Decimals))

	// 验证转换结果
	if balance.IsNegative() {
		return decimal.Zero, gerror.Newf("转换后的余额不能为负数: RawBalance=%s, Symbol=%s", rawBalance.String(), symbol)
	}

	return balance, nil
}

// ConvertBalanceToDBStorage 将余额转换为数据库存储的int64格式
// 这是系统中唯一应该用于数据库存储转换的方法
func (l *tokenLogic) ConvertBalanceToDBStorage(ctx context.Context, balance decimal.Decimal, symbol string) (int64, error) {
	g.Log().Debugf(ctx, "ConvertBalanceToDBStorage: Called with balance='%s', symbol='%s'", balance.String(), symbol)
	rawBalance, err := l.ConvertBalanceToRaw(ctx, balance, symbol)
	g.Log().Debugf(ctx, "ConvertBalanceToDBStorage: For symbol '%s', ConvertBalanceToRaw returned rawBalance='%s', err=%v",
		symbol, rawBalance.String(), err)
	if err != nil {
		return 0, gerror.Wrapf(err, "转换余额到原始格式失败: Balance=%s, Symbol=%s", balance.String(), symbol)
	}

	// 确保转换为int64不会溢出
	// Removed duplicate log block, one is sufficient.
	g.Log().Debugf(ctx, "ConvertBalanceToDBStorage: For symbol '%s', before IsInteger check, rawBalance='%s', rawBalance.IsInteger()=%t",
		symbol, rawBalance.String(), rawBalance.IsInteger())
	if !rawBalance.IsInteger() {
		return 0, gerror.Newf("原始余额不是整数: %s", rawBalance.String())
	}

	storageValue := rawBalance.IntPart()
	return storageValue, nil
}

// ConvertDBStorageToBalance 将数据库存储的int64格式转换为余额
// 这是系统中唯一应该用于从数据库存储转换的方法
func (l *tokenLogic) ConvertDBStorageToBalance(ctx context.Context, storageAmount int64, symbol string) (decimal.Decimal, error) {
	if storageAmount < 0 {
		return decimal.Zero, gerror.Newf("存储金额不能为负数: %d", storageAmount)
	}

	rawBalance := decimal.NewFromInt(storageAmount)
	balance, err := l.ConvertRawToBalance(ctx, rawBalance, symbol)
	if err != nil {
		return decimal.Zero, gerror.Wrapf(err, "转换存储金额到余额失败: StorageAmount=%d, Symbol=%s", storageAmount, symbol)
	}

	return balance, nil
}

// FormatUserBalance 格式化用户余额并移除多余小数点和零
func (l *tokenLogic) FormatUserBalance(ctx context.Context, balance decimal.Decimal, symbol string) (decimal.Decimal, error) {
	token, err := l.GetTokenBySymbol(ctx, symbol)
	if err != nil {
		return decimal.Zero, gerror.Wrapf(err, "获取代币信息失败: Symbol=%s", symbol)
	}

	// 根据代币精度格式化余额
	formatted := balance.Truncate(int32(token.Decimals))
	return formatted, nil
}
