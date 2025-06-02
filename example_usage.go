package wallet

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shopspring/decimal"
)

// ExampleUsage 展示如何使用独立的钱包模块
func ExampleUsage() {
	ctx := context.Background()

	// 获取钱包管理器单例实例
	// 注意：钱包模块必须在 boot.go 中完成初始化
	// 然后在系统任何地方都可以直接使用，无需依赖注入
	walletMgr := Manager()

	// 示例1: 查询用户余额
	fmt.Println("=== 示例1: 查询用户余额 ===")
	userID := uint64(12345)
	tokenSymbol := "USDT"

	balanceInfo, err := walletMgr.GetBalance(ctx, userID, tokenSymbol)
	if err != nil {
		g.Log().Errorf(ctx, "查询余额失败: %v", err)
	} else {
		fmt.Printf("用户 %d 的 %s 余额:\n", userID, tokenSymbol)
		fmt.Printf("  可用余额: %s\n", balanceInfo.AvailableBalance.String())
		fmt.Printf("  冻结余额: %s\n", balanceInfo.FrozenBalance.String())
		fmt.Printf("  总余额: %s\n", balanceInfo.TotalBalance.String())
		fmt.Printf("  最后同步时间: %s\n", balanceInfo.LastSyncTime)
	}

	// 示例2: 在事务中增加资金
	fmt.Println("\n=== 示例2: 增加资金 ===")
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		req := &FundOperationRequest{
			UserID:      userID,
			TokenSymbol: tokenSymbol,
			Amount:      decimal.NewFromFloat(100.50),
			BusinessID:  "deposit_20240527_001",
			Description: "用户充值",
			Metadata: map[string]string{
				"source": "bank_transfer",
				"ref_id": "TXN123456",
			},
		}

		result, err := walletMgr.CreditFundsInTx(ctx, tx, req)
		if err != nil {
			return err
		}

		fmt.Printf("资金增加成功:\n")
		fmt.Printf("  交易ID: %s\n", result.TransactionID)
		fmt.Printf("  操作前余额: %s\n", result.BalanceBefore.String())
		fmt.Printf("  操作后余额: %s\n", result.BalanceAfter.String())
		fmt.Printf("  原始金额: %s\n", result.RawAmount.String())

		return nil
	})

	if err != nil {
		g.Log().Errorf(ctx, "增加资金失败: %v", err)
	}

	// 示例3: 在事务中减少资金
	fmt.Println("\n=== 示例3: 减少资金 ===")
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		req := &FundOperationRequest{
			UserID:      userID,
			TokenSymbol: tokenSymbol,
			Amount:      decimal.NewFromFloat(50.25),
			BusinessID:  "withdraw_20240527_001",
			Description: "用户提现",
			Metadata: map[string]string{
				"destination": "user_wallet",
				"address":     "TXXXxxxXXXxxxXXX",
			},
		}

		result, err := walletMgr.DebitFundsInTx(ctx, tx, req)
		if err != nil {
			return err
		}

		fmt.Printf("资金减少成功:\n")
		fmt.Printf("  交易ID: %s\n", result.TransactionID)
		fmt.Printf("  操作前余额: %s\n", result.BalanceBefore.String())
		fmt.Printf("  操作后余额: %s\n", result.BalanceAfter.String())
		fmt.Printf("  原始金额: %s\n", result.RawAmount.String())

		return nil
	})

	if err != nil {
		g.Log().Errorf(ctx, "减少资金失败: %v", err)
	}

	fmt.Println("\n=== 钱包模块使用示例完成 ===")
}

// 展示钱包模块的独立性和单例模式
func ShowIndependence() {
	fmt.Println("=== 钱包模块独立性和单例模式说明 ===")
	fmt.Println("1. 单例模式：在 boot.go 中初始化一次，全系统共享同一实例")
	fmt.Println("2. 无依赖注入：系统任何地方都可以直接调用 wallet.Manager()")
	fmt.Println("3. 线程安全：使用 sync.Once 确保初始化线程安全")
	fmt.Println("4. 不依赖 internal/service 包中的其他服务")
	fmt.Println("5. 拥有独立的 DAO 层进行数据访问")
	fmt.Println("6. 拥有独立的 Logic 层处理业务逻辑")
	fmt.Println("7. 支持用户管理、代币管理、余额管理、财务操作等完整功能")
	fmt.Println("8. 可以独立部署和测试")
	fmt.Println("9. 通过统一的 IWalletManager 接口对外提供服务")
	fmt.Println()
	fmt.Println("=== 使用方式 ===")
	fmt.Println("// 在系统任何地方直接使用")
	fmt.Println("import \"github.com/yalks/wallet\"")
	fmt.Println()
	fmt.Println("func someFunction(ctx context.Context) {")
	fmt.Println("    // 直接获取单例实例，无需依赖注入")
	fmt.Println("    walletMgr := wallet.Manager()")
	fmt.Println("    balance, err := walletMgr.GetBalance(ctx, userID, \"USDT\")")
	fmt.Println("}")
}
