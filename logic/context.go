package logic

import (
	"context"
	"sync"

	ledgerwalletsdk "github.com/a19ba14d/ledger-wallet-sdk"
	"github.com/a19ba14d/ledger-wallet-sdk/pkg/sdkconfig"
	"github.com/gogf/gf/v2/frame/g"

	"github.com/yalks/wallet/dao"
)

// SharedLogicContext 共享逻辑上下文，减少重复初始化
type SharedLogicContext struct {
	// DAO层
	userDAO        dao.IUserDAO
	tokenDAO       dao.ITokenDAO
	walletDAO      dao.IWalletDAO
	transactionDAO dao.ITransactionDAO

	// 钱包SDK
	walletSDK ledgerwalletsdk.IWallet

	// 初始化标志
	initialized bool
	mu          sync.RWMutex
}

var (
	sharedContext *SharedLogicContext
	contextOnce   sync.Once
)

// GetSharedContext 获取共享逻辑上下文单例
func GetSharedContext() *SharedLogicContext {
	contextOnce.Do(func() {
		sharedContext = &SharedLogicContext{
			userDAO:        dao.NewUserDAO(),
			tokenDAO:       dao.NewTokenDAO(),
			walletDAO:      dao.NewWalletDAO(),
			transactionDAO: dao.NewTransactionDAO(),
		}
		sharedContext.initWalletSDK()
		sharedContext.initialized = true
	})
	return sharedContext
}

// initWalletSDK 初始化钱包SDK
func (c *SharedLogicContext) initWalletSDK() {
	baseURL := g.Cfg().MustGet(context.Background(), "walletsApi.baseUrl").String()
	if client, err := ledgerwalletsdk.New(sdkconfig.WithBaseURL(baseURL)); err != nil {
		g.Log().Errorf(context.Background(), "初始化钱包SDK失败: %v", err)
	} else {
		c.walletSDK = client
	}
}

// GetUserDAO 获取用户DAO
func (c *SharedLogicContext) GetUserDAO() dao.IUserDAO {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.userDAO
}

// GetTokenDAO 获取代币DAO
func (c *SharedLogicContext) GetTokenDAO() dao.ITokenDAO {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.tokenDAO
}

// GetWalletDAO 获取钱包DAO
func (c *SharedLogicContext) GetWalletDAO() dao.IWalletDAO {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.walletDAO
}

// GetTransactionDAO 获取交易DAO
func (c *SharedLogicContext) GetTransactionDAO() dao.ITransactionDAO {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.transactionDAO
}

// GetWalletSDK 获取钱包SDK
func (c *SharedLogicContext) GetWalletSDK() ledgerwalletsdk.IWallet {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.walletSDK
}

// IsInitialized 检查是否已初始化
func (c *SharedLogicContext) IsInitialized() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.initialized
}
