package logic

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/shopspring/decimal"

	"github.com/yalks/wallet/entity"
)

// IBalanceLogic 余额业务逻辑接口
type IBalanceLogic interface {
	// GetBalance 获取用户余额
	GetBalance(ctx context.Context, userID uint64, tokenSymbol string) (availableBalance, frozenBalance decimal.Decimal, err error)
	// GetLocalBalance 获取本地钱包余额
	GetLocalBalance(ctx context.Context, userID uint64, tokenSymbol string) (availableBalance, frozenBalance decimal.Decimal, err error)
	// GetRemoteBalance 获取远程钱包余额
	GetRemoteBalance(ctx context.Context, userID uint64, tokenSymbol string) (decimal.Decimal, error)
	// UpdateLocalBalance 更新本地钱包余额
	UpdateLocalBalance(ctx context.Context, tx gdb.TX, userID uint64, tokenSymbol string, availableBalance, frozenBalance decimal.Decimal) error
	// SyncBalanceFromRemote 从远程同步余额到本地
	// SyncBalanceFromRemote(ctx context.Context, userID uint64, tokenSymbol string) error
}

type balanceLogic struct {
	userLogic  IUserLogic
	tokenLogic ITokenLogic
	context    *SharedLogicContext
}

// NewBalanceLogic 创建余额业务逻辑实例
func NewBalanceLogic() IBalanceLogic {
	return &balanceLogic{
		userLogic:  NewUserLogic(),
		tokenLogic: NewTokenLogic(),
		context:    GetSharedContext(),
	}
}

// GetBalance 获取用户余额（优先本地，本地没有则创建）- 暂时不查询远程
func (l *balanceLogic) GetBalance(ctx context.Context, userID uint64, tokenSymbol string) (availableBalance, frozenBalance decimal.Decimal, err error) {
	// 1. 首先尝试从本地获取
	availableBalance, frozenBalance, err = l.GetLocalBalance(ctx, userID, tokenSymbol)
	if err == nil {
		return availableBalance, frozenBalance, nil
	}

	g.Log().Infof(ctx, "本地钱包记录不存在，创建初始余额为0的记录: UserID=%d, Symbol=%s", userID, tokenSymbol)

	// 2. 如果本地没有记录，创建一个初始余额为0的记录
	// 暂时注释远程查询逻辑
	// remoteBalance, err := l.GetRemoteBalance(ctx, userID, tokenSymbol)
	// if err != nil {
	// 	return decimal.Zero, decimal.Zero, err
	// }

	// 3. 创建本地钱包记录，初始余额为0
	err = l.createLocalWalletRecord(ctx, userID, tokenSymbol, decimal.Zero)
	if err != nil {
		g.Log().Warningf(ctx, "创建本地钱包记录失败: UserID=%d, Symbol=%s, Error=%v", userID, tokenSymbol, err)
		// 不阻止操作，但记录警告
	}

	return decimal.Zero, decimal.Zero, nil
}

// GetLocalBalance 获取本地钱包余额
func (l *balanceLogic) GetLocalBalance(ctx context.Context, userID uint64, tokenSymbol string) (availableBalance, frozenBalance decimal.Decimal, err error) {
	wallet, err := l.context.GetWalletDAO().GetWalletByUserIDAndSymbol(ctx, userID, tokenSymbol)
	if err != nil {
		return decimal.Zero, decimal.Zero, gerror.Wrapf(err, "获取本地钱包记录失败: UserID=%d, Symbol=%s", userID, tokenSymbol)
	}
	if wallet == nil {
		return decimal.Zero, decimal.Zero, gerror.Newf("本地钱包记录不存在: UserID=%d, Symbol=%s", userID, tokenSymbol)
	}

	// 转换余额为decimal格式 - 使用统一的TokenLogic数据库转换方法
	availableBalance, err = l.tokenLogic.ConvertDBStorageToBalance(ctx, wallet.AvailableBalance, tokenSymbol)
	if err != nil {
		return decimal.Zero, decimal.Zero, gerror.Wrapf(err, "转换可用余额失败: UserID=%d, Symbol=%s", userID, tokenSymbol)
	}

	frozenBalance, err = l.tokenLogic.ConvertDBStorageToBalance(ctx, wallet.FrozenBalance, tokenSymbol)
	if err != nil {
		return decimal.Zero, decimal.Zero, gerror.Wrapf(err, "转换冻结余额失败: UserID=%d, Symbol=%s", userID, tokenSymbol)
	}

	return availableBalance, frozenBalance, nil
}

// GetRemoteBalance 获取远程钱包余额 - 暂时返回错误，使用本地事务
func (l *balanceLogic) GetRemoteBalance(ctx context.Context, userID uint64, tokenSymbol string) (decimal.Decimal, error) {
	// 暂时不支持远程查询，直接返回错误
	return decimal.Zero, gerror.New("远程钱包查询暂时禁用，请使用本地余额")
	
	// 原实现已注释，待后续恢复远程钱包功能时启用
	/*
	walletSDK := l.context.GetWalletSDK()
	if walletSDK == nil {
		return decimal.Zero, gerror.New("钱包SDK未初始化")
	}

	// 获取用户主钱包ID
	mainWalletID, err := l.userLogic.GetMainWalletID(ctx, userID)
	if err != nil {
		return decimal.Zero, err
	}

	// 从SDK查询余额
	walletSummary, err := walletSDK.GetWalletSummary(ctx, mainWalletID)
	if err != nil {
		return decimal.Zero, gerror.Wrapf(err, "查询远程钱包摘要失败: WalletID=%s", mainWalletID)
	}
	if walletSummary == nil {
		return decimal.Zero, gerror.Newf("钱包摘要为空: WalletID=%s", mainWalletID)
	}

	// 从SDK响应中提取指定代币的余额
	if walletSummary.AvailableFunds != nil {
		if availableAmountInt64, exists := walletSummary.AvailableFunds[tokenSymbol]; exists {
			// 从原始值转换为decimal（使用统一的TokenLogic方法）
			rawBalance := decimal.NewFromInt(availableAmountInt64)
			balance, err := l.tokenLogic.ConvertRawToBalance(ctx, rawBalance, tokenSymbol)
			if err != nil {
				return decimal.Zero, gerror.Wrapf(err, "转换远程余额格式失败: RawBalance=%s, Symbol=%s", rawBalance.String(), tokenSymbol)
			}
			return balance, nil
		}
	}

	return decimal.Zero, nil
	*/
}

// UpdateLocalBalance 更新本地钱包余额
func (l *balanceLogic) UpdateLocalBalance(ctx context.Context, tx gdb.TX, userID uint64, tokenSymbol string, availableBalance, frozenBalance decimal.Decimal) error {
	// 获取本地钱包记录
	wallet, err := l.context.GetWalletDAO().GetWalletByUserIDAndSymbol(ctx, userID, tokenSymbol)
	if err != nil {
		return gerror.Wrapf(err, "获取本地钱包记录失败: UserID=%d, Symbol=%s", userID, tokenSymbol)
	}
	if wallet == nil {
		return gerror.Newf("本地钱包记录不存在: UserID=%d, Symbol=%s", userID, tokenSymbol)
	}

	// 转换余额为存储格式（最小单位）- 使用统一的TokenLogic数据库转换方法
	g.Log().Debugf(ctx, "UpdateLocalBalance: Calling ConvertBalanceToDBStorage for symbol '%s' with availableBalance='%s'",
		tokenSymbol, availableBalance.String())
	availableBalanceInt, err := l.tokenLogic.ConvertBalanceToDBStorage(ctx, availableBalance, tokenSymbol)
	if err != nil {
		return gerror.Wrapf(err, "转换可用余额到存储格式失败: UserID=%d, Symbol=%s", userID, tokenSymbol)
	}

	frozenBalanceInt, err := l.tokenLogic.ConvertBalanceToDBStorage(ctx, frozenBalance, tokenSymbol)
	if err != nil {
		return gerror.Wrapf(err, "转换冻结余额到存储格式失败: UserID=%d, Symbol=%s", userID, tokenSymbol)
	}

	// 更新余额
	err = l.context.GetWalletDAO().UpdateWalletBalance(ctx, tx, uint(wallet.WalletId), availableBalanceInt, frozenBalanceInt)
	if err != nil {
		return err
	}

	g.Log().Debugf(ctx, "更新本地钱包余额: UserID=%d, Symbol=%s, Available=%s, Frozen=%s",
		userID, tokenSymbol, availableBalance.String(), frozenBalance.String())
	return nil
}

// // SyncBalanceFromRemote 从远程同步余额到本地
// func (l *balanceLogic) SyncBalanceFromRemote(ctx context.Context, userID uint64, tokenSymbol string) error {
// 	// 获取远程余额
// 	remoteBalance, err := l.GetRemoteBalance(ctx, userID, tokenSymbol)
// 	if err != nil {
// 		return gerror.Wrap(err, "获取远程余额失败")
// 	}

// 	// 获取本地余额
// 	localBalance, _, err := l.GetLocalBalance(ctx, userID, tokenSymbol)
// 	if err != nil {
// 		g.Log().Warningf(ctx, "获取本地余额失败，将创建新记录: %v", err)
// 		// 如果本地记录不存在，这里可以考虑创建新记录
// 		// 但为了简化，暂时只更新已存在的记录
// 		return gerror.Wrap(err, "本地钱包记录不存在，无法同步")
// 	}

// 	// 如果余额一致，无需同步
// 	if localBalance.Equal(remoteBalance) {
// 		g.Log().Infof(ctx, "余额已一致，无需同步: UserID=%d, Symbol=%s", userID, tokenSymbol)
// 		return nil
// 	}

// 	// 更新本地钱包余额
// 	err = l.UpdateLocalBalance(ctx, nil, userID, tokenSymbol, remoteBalance, decimal.Zero)
// 	if err != nil {
// 		return gerror.Wrap(err, "更新本地钱包余额失败")
// 	}

// 	g.Log().Infof(ctx, "余额同步完成: UserID=%d, Symbol=%s, 从 %s 同步到 %s",
// 		userID, tokenSymbol, localBalance.String(), remoteBalance.String())

// 	return nil
// }

// createLocalWalletRecord 创建本地钱包记录
func (l *balanceLogic) createLocalWalletRecord(ctx context.Context, userID uint64, tokenSymbol string, initialBalance decimal.Decimal) error {
	// 获取代币信息
	token, err := l.tokenLogic.GetTokenBySymbol(ctx, tokenSymbol)
	if err != nil {
		return gerror.Wrapf(err, "获取代币信息失败: Symbol=%s", tokenSymbol)
	}

	// 转换初始余额为存储格式
	initialBalanceInt, err := l.tokenLogic.ConvertBalanceToDBStorage(ctx, initialBalance, tokenSymbol)
	if err != nil {
		return gerror.Wrapf(err, "转换初始余额到存储格式失败: Balance=%s, Symbol=%s", initialBalance.String(), tokenSymbol)
	}

	// 创建钱包记录
	walletRecord := &entity.Wallets{
		UserId:           uint(userID),
		TokenId:          uint(token.TokenId),
		Symbol:           tokenSymbol,
		AvailableBalance: initialBalanceInt,
		FrozenBalance:    0,
		DecimalPlaces:    uint(token.Decimals),
		CreatedAt:        gtime.Now(),
		UpdatedAt:        gtime.Now(),
	}

	err = l.context.GetWalletDAO().CreateWalletRecord(ctx, nil, walletRecord)
	if err != nil {
		return gerror.Wrapf(err, "创建本地钱包记录失败: UserID=%d, TokenID=%d", userID, token.TokenId)
	}

	g.Log().Infof(ctx, "本地钱包记录创建成功: UserID=%d, Symbol=%s, InitialBalance=%s",
		userID, tokenSymbol, initialBalance.String())
	return nil
}
