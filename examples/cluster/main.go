package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"cache"
)

func main() {
	// 创建集群模式配置
	config := &cache.Config{
		Mode: cache.ModeCluster,
		Cluster: &cache.ClusterConfig{
			Addrs: []string{
				"localhost:7000",
				"localhost:7001",
				"localhost:7002",
				"localhost:7003",
				"localhost:7004",
				"localhost:7005",
			},
			MaxRedirects: 3,
			ReadOnly:     false,
		},
		Common: cache.CommonConfig{
			Password:     "", // 如果有密码请填写
			KeyPrefix:    "cluster_app:",
			DefaultTTL:   30 * time.Minute,
			PoolSize:     20, // 集群模式建议更大的连接池
			MinIdleConns: 10,
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

	// 创建Redis集群客户端
	client, err := factory.CreateClient()
	if err != nil {
		log.Fatalf("创建Redis集群客户端失败: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// 测试连接
	if err := client.Ping(ctx); err != nil {
		log.Fatalf("Redis集群连接测试失败: %v", err)
	}
	fmt.Println("✅ Redis集群连接成功")

	// 集群模式特性演示
	fmt.Println("\n=== Redis集群模式特性演示 ===")
	clusterFeatures(ctx, client)

	// 分布式数据操作
	fmt.Println("\n=== 分布式数据操作 ===")
	distributedOperations(ctx, client)

	// 批量操作演示
	fmt.Println("\n=== 批量操作演示 ===")
	batchOperations(ctx, client)

	// 高可用性测试
	fmt.Println("\n=== 高可用性特性 ===")
	highAvailabilityDemo(ctx, client)

	fmt.Println("\n🎉 Redis集群操作演示完成")
}

// 集群模式特性演示
func clusterFeatures(ctx context.Context, client cache.Client) {
	// 在不同的槽位存储数据，展示数据分布
	keys := []string{
		"user:1001", // 这些键会被分布到不同的节点
		"user:2002",
		"user:3003",
		"product:A001",
		"product:B002",
		"order:O001",
	}

	fmt.Println("在集群中分布存储数据...")
	for i, key := range keys {
		value := fmt.Sprintf("数据_%d", i+1)
		err := client.Set(ctx, key, value, 10*time.Minute)
		if err != nil {
			log.Printf("设置 %s 失败: %v", key, err)
			continue
		}
		fmt.Printf("✓ 设置 %s = %s\n", key, value)
	}

	// 验证数据读取
	fmt.Println("\n验证分布式数据读取...")
	for _, key := range keys {
		value, err := client.Get(ctx, key)
		if err != nil {
			log.Printf("获取 %s 失败: %v", key, err)
			continue
		}
		fmt.Printf("✓ 读取 %s = %s\n", key, value)
	}
}

// 分布式数据操作
func distributedOperations(ctx context.Context, client cache.Client) {
	// 分布式计数器
	counters := []string{"counter:page_view", "counter:api_call", "counter:user_login"}

	fmt.Println("分布式计数器操作...")
	for _, counter := range counters {
		// 递增计数器
		for i := 0; i < 5; i++ {
			count, err := client.Incr(ctx, counter)
			if err != nil {
				log.Printf("递增计数器 %s 失败: %v", counter, err)
				break
			}
			if i == 4 { // 只打印最后一次结果
				fmt.Printf("✓ %s 当前值: %d\n", counter, count)
			}
		}
	}

	// 分布式哈希表
	fmt.Println("\n分布式哈希表操作...")
	userSessions := map[string]map[string]string{
		"session:user1": {
			"user_id":    "1001",
			"username":   "张三",
			"login_time": time.Now().Format(time.RFC3339),
		},
		"session:user2": {
			"user_id":    "1002",
			"username":   "李四",
			"login_time": time.Now().Format(time.RFC3339),
		},
	}

	for sessionKey, sessionData := range userSessions {
		for field, value := range sessionData {
			err := client.HSet(ctx, sessionKey, field, value)
			if err != nil {
				log.Printf("设置会话数据失败: %v", err)
				continue
			}
		}
		fmt.Printf("✓ 创建用户会话: %s\n", sessionKey)
	}

	// 读取会话数据
	for sessionKey := range userSessions {
		sessionData, err := client.HGetAll(ctx, sessionKey)
		if err != nil {
			log.Printf("获取会话数据失败: %v", err)
			continue
		}
		fmt.Printf("✓ %s 数据: %v\n", sessionKey, sessionData)
	}
}

// 批量操作演示
func batchOperations(ctx context.Context, client cache.Client) {
	// 使用管道进行批量操作
	pipe := client.Pipeline()

	// 批量设置用户信息
	userData := map[string]string{
		"batch:user:1": "用户1信息",
		"batch:user:2": "用户2信息",
		"batch:user:3": "用户3信息",
		"batch:user:4": "用户4信息",
		"batch:user:5": "用户5信息",
	}

	fmt.Println("使用管道批量设置数据...")
	for key, value := range userData {
		pipe.Set(ctx, key, value, 5*time.Minute)
	}

	// 批量递增计数器
	for i := 1; i <= 3; i++ {
		counterKey := fmt.Sprintf("batch:counter:%d", i)
		pipe.Incr(ctx, counterKey)
	}

	// 执行管道
	results, err := pipe.Exec(ctx)
	if err != nil {
		log.Printf("执行管道失败: %v", err)
		return
	}

	fmt.Printf("✓ 批量操作完成，执行了 %d 个命令\n", len(results))

	// 关闭管道
	pipe.Close()

	// 验证批量设置的数据
	fmt.Println("\n验证批量设置的数据...")
	for key := range userData {
		value, err := client.Get(ctx, key)
		if err != nil {
			log.Printf("获取 %s 失败: %v", key, err)
			continue
		}
		fmt.Printf("✓ %s = %s\n", key, value)
	}
}

// 高可用性演示
func highAvailabilityDemo(ctx context.Context, client cache.Client) {
	fmt.Println("演示集群的高可用性特性...")

	// 设置一些测试数据
	testKeys := []string{"ha:test1", "ha:test2", "ha:test3"}
	for i, key := range testKeys {
		value := fmt.Sprintf("高可用测试数据_%d", i+1)
		err := client.Set(ctx, key, value, 10*time.Minute)
		if err != nil {
			log.Printf("设置测试数据失败: %v", err)
			continue
		}
		fmt.Printf("✓ 设置测试数据: %s\n", key)
	}

	// 模拟连续读取操作（在实际环境中，即使某个节点故障，集群仍能正常服务）
	fmt.Println("\n执行连续读取操作（模拟高可用场景）...")
	for round := 1; round <= 3; round++ {
		fmt.Printf("第 %d 轮读取:\n", round)
		for _, key := range testKeys {
			value, err := client.Get(ctx, key)
			if err != nil {
				log.Printf("  ✗ 读取 %s 失败: %v", key, err)
			} else {
				fmt.Printf("  ✓ 读取 %s = %s\n", key, value)
			}
		}
		time.Sleep(1 * time.Second)
	}

	// 集合操作测试（跨节点）
	fmt.Println("\n跨节点集合操作测试...")
	setKey := "ha:distributed_set"
	members := []interface{}{"member1", "member2", "member3", "member4", "member5"}

	_, err := client.SAdd(ctx, setKey, members...)
	if err != nil {
		log.Printf("添加集合成员失败: %v", err)
		return
	}

	setMembers, err := client.SMembers(ctx, setKey)
	if err != nil {
		log.Printf("获取集合成员失败: %v", err)
		return
	}

	fmt.Printf("✓ 分布式集合成员: %v\n", setMembers)
	fmt.Printf("✓ 集合大小: %d\n", len(setMembers))

	fmt.Println("\n💡 提示: 在生产环境中，Redis集群能够：")
	fmt.Println("   - 自动进行数据分片和负载均衡")
	fmt.Println("   - 在节点故障时自动故障转移")
	fmt.Println("   - 支持水平扩展，可动态添加/移除节点")
	fmt.Println("   - 提供高可用性和数据冗余")
}