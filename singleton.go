package wallet

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
)

var (
	// 单例实例
	managerInstance IWalletManager
	once            sync.Once
	initialized     bool
)

// Initialize 初始化钱包管理器（由 boot.go 调用）
func Initialize(ctx context.Context) error {
	var initErr error
	once.Do(func() {
		g.Log().Info(ctx, "开始初始化钱包管理器...")

		// 创建钱包管理器实例
		manager := &walletManager{}

		// 初始化各个逻辑组件
		if err := manager.initialize(ctx); err != nil {
			initErr = err
			g.Log().Errorf(ctx, "钱包管理器初始化失败: %v", err)
			return
		}

		managerInstance = manager
		initialized = true
		g.Log().Info(ctx, "钱包管理器初始化完成")
	})

	return initErr
}

// Manager 获取钱包管理器单例
func Manager() IWalletManager {
	if !initialized || managerInstance == nil {
		panic("钱包管理器未初始化，请先调用 wallet.Initialize()")
	}
	return managerInstance
}

// IsInitialized 检查钱包管理器是否已初始化
func IsInitialized() bool {
	return initialized && managerInstance != nil
}

// SetManager 设置钱包管理器实例（主要用于测试）
func SetManager(manager IWalletManager) {
	managerInstance = manager
	initialized = true
}
