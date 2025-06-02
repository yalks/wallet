package wallet

import (
	"context"
	"fmt"

	"github.com/yalks/wallet/constants"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
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

	// 示例2: 使用 TransactionBuilder 创建充值交易（自动生成Reference）
	fmt.Println("\n=== 示例2: 使用Builder模式增加资金（自动生成Reference） ===")
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 使用 TransactionBuilder 构建交易请求
		txReq, err := constants.NewTransactionBuilder().
			WithUser(int64(userID)).
			WithWallet(1). // 假设钱包ID为1，实际应用中需要获取真实钱包ID
			WithAmount("100.50").
			WithToken(1). // 假设USDT的token ID为1，实际应用中需要获取真实token ID
			WithFundType(constants.FundTypeDeposit).
			// WithReference 是可选的，这里不设置，将自动生成格式: {FundType}_{UserID}_{Timestamp}
			// 不调用 WithReference，让系统自动生成
			WithDescription("用户充值").
			WithRequestSource("api").
			WithMetadata("source", "bank_transfer").
			WithMetadata("ref_id", "TXN123456").
			Build()
		if err != nil {
			return err
		}

		transactionID, err := walletMgr.CreateTransactionWithBuilder(ctx, txReq)
		if err != nil {
			return err
		}

		fmt.Printf("资金增加成功:\n")
		fmt.Printf("  交易ID: %d\n", transactionID)
		fmt.Printf("  金额: %s\n", txReq.Amount)
		fmt.Printf("  资金类型: %s\n", txReq.FundType)
		fmt.Printf("  业务引用（自动生成）: %s\n", txReq.Reference)

		return nil
	})

	if err != nil {
		g.Log().Errorf(ctx, "增加资金失败: %v", err)
	}

	// 示例3: 使用 TransactionBuilder 创建提现交易（用户自定义Reference）
	fmt.Println("\n=== 示例3: 使用Builder模式减少资金（用户自定义Reference） ===")
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 使用 TransactionBuilder 构建提现交易请求
		txReq, err := constants.NewTransactionBuilder().
			WithUser(int64(userID)).
			WithWallet(1). // 假设钱包ID为1，实际应用中需要获取真实钱包ID
			WithAmount("50.25").
			WithToken(1). // 假设USDT的token ID为1，实际应用中需要获取真实token ID
			WithFundType(constants.FundTypeWithdraw).
			// 用户提供自定义引用ID（如银行流水号、第三方支付订单号等）
			WithReference("BANK_TXN_withdraw_20240527_001").
			WithDescription("用户提现").
			WithRequestSource("api").
			WithMetadata("destination", "user_wallet").
			WithMetadata("address", "TXXXxxxXXXxxxXXX").
			Build()
		if err != nil {
			return err
		}

		transactionID, err := walletMgr.CreateTransactionWithBuilder(ctx, txReq)
		if err != nil {
			return err
		}

		fmt.Printf("资金减少成功:\n")
		fmt.Printf("  交易ID: %d\n", transactionID)
		fmt.Printf("  金额: %s\n", txReq.Amount)
		fmt.Printf("  资金类型: %s\n", txReq.FundType)
		fmt.Printf("  业务引用（用户自定义）: %s\n", txReq.Reference)

		return nil
	})

	if err != nil {
		g.Log().Errorf(ctx, "减少资金失败: %v", err)
	}

	// 示例4: 使用预定义的Builder函数创建转账交易
	fmt.Println("\n=== 示例4: 使用预定义Builder创建转账 ===")
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 使用预定义的Builder函数创建转账交易
		txReq, err := constants.BuildTransferOutTransaction(
			int64(userID), // 发送方用户ID
			1,             // 钱包ID
			"25.75",       // 转账金额
			1,             // token ID
			12345,         // 转账ID
			67890,         // 接收方用户ID
			"recipient_user", // 接收方用户名
		)
		if err != nil {
			return err
		}

		transactionID, err := walletMgr.CreateTransactionWithBuilder(ctx, txReq)
		if err != nil {
			return err
		}

		fmt.Printf("转账交易创建成功:\n")
		fmt.Printf("  交易ID: %d\n", transactionID)
		fmt.Printf("  金额: %s\n", txReq.Amount)
		fmt.Printf("  接收方: %s (ID: %d)\n", txReq.TargetUsername, txReq.TargetUserID)
		fmt.Printf("  转账引用ID: %d\n", txReq.RelatedID)

		return nil
	})

	if err != nil {
		g.Log().Errorf(ctx, "创建转账交易失败: %v", err)
	}

	// 示例5: 展示Reference的不同用法对比
	fmt.Println("\n=== 示例5: Reference用法对比 ===")
	
	// 5.1 完全省略WithReference，让系统自动生成
	fmt.Println("5.1 自动生成Reference:")
	txReq1, _ := constants.NewTransactionBuilder().
		WithUser(int64(userID)).
		WithWallet(1).
		WithAmount("10.00").
		WithToken(1).
		WithFundType(constants.FundTypeDeposit).
		Build()
	fmt.Printf("  自动生成的Reference: %s\n", txReq1.Reference)
	
	// 5.2 用户自定义Reference
	fmt.Println("5.2 用户自定义Reference:")
	txReq2, _ := constants.NewTransactionBuilder().
		WithUser(int64(userID)).
		WithWallet(1).
		WithAmount("10.00").
		WithToken(1).
		WithFundType(constants.FundTypeDeposit).
		WithReference("CUSTOM_REF_12345").
		Build()
	fmt.Printf("  用户自定义的Reference: %s\n", txReq2.Reference)

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
