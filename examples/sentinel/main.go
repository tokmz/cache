package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"cache"
)

func main() {
	// 创建哨兵模式配置
	config := &cache.Config{
		Mode: cache.ModeSentinel,
		Sentinel: &cache.SentinelConfig{
			Addrs: []string{
				"localhost:26379", // 哨兵1
				"localhost:26380", // 哨兵2
				"localhost:26381", // 哨兵3
			},
			MasterName:       "mymaster", // 主服务器名称
			DB:               0,
			SentinelPassword: "", // 哨兵密码（如果有）
		},
		Common: cache.CommonConfig{
			Password:     "", // Redis密码（如果有）
			KeyPrefix:    "sentinel_app:",
			DefaultTTL:   30 * time.Minute,
			PoolSize:     15,
			MinIdleConns: 8,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
			MaxRetries:   3,
		},
	}

	// 验证配置
	if err := config.Validate(); err != nil {
		log.Fatalf("配置验证失败: %v", err)
	}

	// 创建客户端工厂
	factory, err := cache.NewFactory(config)
	if err != nil {
		log.Fatalf("创建工厂失败: %v", err)
	}

	// 创建Redis哨兵客户端
	client, err := factory.CreateClient()
	if err != nil {
		log.Fatalf("创建Redis哨兵客户端失败: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// 测试连接
	if err := client.Ping(ctx); err != nil {
		log.Fatalf("Redis哨兵连接测试失败: %v", err)
	}
	fmt.Println("✅ Redis哨兵模式连接成功")

	// 哨兵模式特性演示
	fmt.Println("\n=== Redis哨兵模式特性演示 ===")
	sentinelFeatures(ctx, client)

	// 高可用性测试
	fmt.Println("\n=== 高可用性和故障转移演示 ===")
	highAvailabilityTest(ctx, client)

	// 读写分离演示
	fmt.Println("\n=== 读写操作演示 ===")
	readWriteOperations(ctx, client)

	// 持续监控演示
	fmt.Println("\n=== 连接监控演示 ===")
	connectionMonitoring(ctx, client)

	fmt.Println("\n🎉 Redis哨兵模式演示完成")
}

// 哨兵模式特性演示
func sentinelFeatures(ctx context.Context, client cache.Client) {
	fmt.Println("演示哨兵模式的基本特性...")

	// 设置一些基础数据
	baseData := map[string]string{
		"app:version":    "1.0.0",
		"app:name":       "哨兵模式演示应用",
		"app:start_time": time.Now().Format(time.RFC3339),
		"app:mode":       "sentinel",
	}

	fmt.Println("设置应用基础信息...")
	for key, value := range baseData {
		err := client.Set(ctx, key, value, 1*time.Hour)
		if err != nil {
			log.Printf("设置 %s 失败: %v", key, err)
			continue
		}
		fmt.Printf("✓ 设置 %s = %s\n", key, value)
	}

	// 验证数据读取
	fmt.Println("\n验证数据读取...")
	for key := range baseData {
		value, err := client.Get(ctx, key)
		if err != nil {
			log.Printf("获取 %s 失败: %v", key, err)
			continue
		}
		fmt.Printf("✓ 读取 %s = %s\n", key, value)
	}

	// 哈希表操作
	fmt.Println("\n哈希表操作演示...")
	serverInfo := map[string]string{
		"hostname": "sentinel-server-01",
		"ip":       "192.168.1.100",
		"port":     "6379",
		"role":     "master",
		"status":   "online",
	}

	for field, value := range serverInfo {
		err := client.HSet(ctx, "server:info", field, value)
		if err != nil {
			log.Printf("设置服务器信息失败: %v", err)
			continue
		}
	}
	fmt.Println("✓ 设置服务器信息")

	info, err := client.HGetAll(ctx, "server:info")
	if err != nil {
		log.Printf("获取服务器信息失败: %v", err)
	} else {
		fmt.Printf("✓ 服务器信息: %v\n", info)
	}
}

// 高可用性测试
func highAvailabilityTest(ctx context.Context, client cache.Client) {
	fmt.Println("演示哨兵模式的高可用性特性...")

	// 创建一些关键业务数据
	businessData := map[string]interface{}{
		"business:total_users":    1000,
		"business:active_sessions": 150,
		"business:daily_revenue":   25000.50,
		"business:last_backup":     time.Now().Format(time.RFC3339),
	}

	fmt.Println("设置关键业务数据...")
	for key, value := range businessData {
		err := client.Set(ctx, key, fmt.Sprintf("%v", value), 2*time.Hour)
		if err != nil {
			log.Printf("设置业务数据失败: %v", err)
			continue
		}
		fmt.Printf("✓ 设置 %s = %v\n", key, value)
	}

	// 模拟连续的读写操作（在主从切换时仍能正常工作）
	fmt.Println("\n执行连续的读写操作（模拟故障转移场景）...")
	for round := 1; round <= 5; round++ {
		fmt.Printf("第 %d 轮操作:\n", round)

		// 写操作：更新活跃会话数
		newSessionCount, err := client.Incr(ctx, "business:active_sessions")
		if err != nil {
			log.Printf("  ✗ 更新会话数失败: %v", err)
		} else {
			fmt.Printf("  ✓ 活跃会话数更新为: %d\n", newSessionCount)
		}

		// 读操作：获取用户总数
		userCount, err := client.Get(ctx, "business:total_users")
		if err != nil {
			log.Printf("  ✗ 读取用户总数失败: %v", err)
		} else {
			fmt.Printf("  ✓ 用户总数: %s\n", userCount)
		}

		// 哈希操作：更新服务器状态
		timestamp := time.Now().Format("15:04:05")
		err = client.HSet(ctx, "server:status", "last_check", timestamp)
		if err != nil {
			log.Printf("  ✗ 更新服务器状态失败: %v", err)
		} else {
			fmt.Printf("  ✓ 服务器状态更新时间: %s\n", timestamp)
		}

		time.Sleep(2 * time.Second)
	}

	fmt.Println("\n💡 哨兵模式优势:")
	fmt.Println("   - 自动故障检测和主从切换")
	fmt.Println("   - 无需手动干预的高可用性")
	fmt.Println("   - 客户端自动重连到新的主服务器")
	fmt.Println("   - 保证数据一致性和服务连续性")
}

// 读写操作演示
func readWriteOperations(ctx context.Context, client cache.Client) {
	fmt.Println("演示读写操作的可靠性...")

	// 批量写入用户数据
	fmt.Println("批量写入用户数据...")
	for i := 1; i <= 10; i++ {
		userKey := fmt.Sprintf("user:%d", i)
		userData := map[string]string{
			"id":         fmt.Sprintf("%d", i),
			"name":       fmt.Sprintf("用户%d", i),
			"email":      fmt.Sprintf("user%d@example.com", i),
			"created_at": time.Now().Format(time.RFC3339),
			"status":     "active",
		}

		for field, value := range userData {
			err := client.HSet(ctx, userKey, field, value)
			if err != nil {
				log.Printf("设置用户数据失败: %v", err)
				break
			}
		}
		fmt.Printf("✓ 创建用户: %s\n", userKey)
	}

	// 批量读取验证
	fmt.Println("\n批量读取验证...")
	for i := 1; i <= 10; i++ {
		userKey := fmt.Sprintf("user:%d", i)
		userData, err := client.HGetAll(ctx, userKey)
		if err != nil {
			log.Printf("读取用户数据失败: %v", err)
			continue
		}
		fmt.Printf("✓ 用户 %s: %s (%s)\n", userData["id"], userData["name"], userData["email"])
	}

	// 列表操作
	fmt.Println("\n列表操作演示...")
	logEntries := []string{
		"[INFO] 应用启动",
		"[INFO] 连接到哨兵模式Redis",
		"[INFO] 开始处理用户请求",
		"[WARN] 检测到高负载",
		"[INFO] 负载恢复正常",
	}

	for _, entry := range logEntries {
		timestampedEntry := fmt.Sprintf("%s %s", time.Now().Format("15:04:05"), entry)
		_, err := client.LPush(ctx, "app:logs", timestampedEntry)
		if err != nil {
			log.Printf("添加日志失败: %v", err)
			continue
		}
	}
	fmt.Println("✓ 添加应用日志")

	// 获取最近的日志
	recentLogs, err := client.LRange(ctx, "app:logs", 0, 4)
	if err != nil {
		log.Printf("获取日志失败: %v", err)
	} else {
		fmt.Println("最近的日志:")
		for _, logEntry := range recentLogs {
			fmt.Printf("  %s\n", logEntry)
		}
	}
}

// 连接监控演示
func connectionMonitoring(ctx context.Context, client cache.Client) {
	fmt.Println("演示连接监控和健康检查...")

	// 持续监控连接状态
	for i := 1; i <= 5; i++ {
		fmt.Printf("第 %d 次健康检查:\n", i)

		// Ping测试
		start := time.Now()
		err := client.Ping(ctx)
		latency := time.Since(start)

		if err != nil {
			fmt.Printf("  ✗ Ping失败: %v\n", err)
		} else {
			fmt.Printf("  ✓ Ping成功，延迟: %v\n", latency)
		}

		// 记录监控数据
		monitorKey := "monitor:health_check"
		monitorData := map[string]string{
			"timestamp": time.Now().Format(time.RFC3339),
			"latency":   latency.String(),
			"status":    "healthy",
			"check_id":  fmt.Sprintf("%d", i),
		}

		for field, value := range monitorData {
			err := client.HSet(ctx, monitorKey, field, value)
			if err != nil {
				log.Printf("  记录监控数据失败: %v", err)
				break
			}
		}
		fmt.Printf("  ✓ 记录监控数据\n")

		// 更新统计计数器
		_, err = client.Incr(ctx, "monitor:total_checks")
		if err != nil {
			log.Printf("  更新检查计数失败: %v", err)
		} else {
			fmt.Printf("  ✓ 更新检查计数\n")
		}

		if i < 5 {
			time.Sleep(3 * time.Second)
		}
	}

	// 获取监控统计
	totalChecks, err := client.Get(ctx, "monitor:total_checks")
	if err != nil {
		log.Printf("获取检查总数失败: %v", err)
	} else {
		fmt.Printf("\n📊 总检查次数: %s\n", totalChecks)
	}

	lastCheck, err := client.HGetAll(ctx, "monitor:health_check")
	if err != nil {
		log.Printf("获取最后检查数据失败: %v", err)
	} else {
		fmt.Printf("📊 最后检查: %v\n", lastCheck)
	}

	fmt.Println("\n💡 监控建议:")
	fmt.Println("   - 定期执行健康检查")
	fmt.Println("   - 监控连接延迟和错误率")
	fmt.Println("   - 设置告警阈值")
	fmt.Println("   - 记录关键指标用于分析")
}