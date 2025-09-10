package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"cache"
)

func main() {
	// 创建单机模式配置
	config := &cache.Config{
		Mode: cache.ModeSingle,
		Single: &cache.SingleConfig{
			Addr: "localhost:6379",
			DB:   0,
		},
		Common: cache.CommonConfig{
			Password:     "", // 如果有密码请填写
			KeyPrefix:    "myapp:",
			DefaultTTL:   30 * time.Minute,
			PoolSize:     10,
			MinIdleConns: 5,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
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

	// 创建Redis客户端
	client, err := factory.CreateClient()
	if err != nil {
		log.Fatalf("创建Redis客户端失败: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// 测试连接
	if err := client.Ping(ctx); err != nil {
		log.Fatalf("Redis连接测试失败: %v", err)
	}
	fmt.Println("✅ Redis连接成功")

	// 字符串操作示例
	fmt.Println("\n=== 字符串操作示例 ===")
	stringOperations(ctx, client)

	// 哈希表操作示例
	fmt.Println("\n=== 哈希表操作示例 ===")
	hashOperations(ctx, client)

	// 列表操作示例
	fmt.Println("\n=== 列表操作示例 ===")
	listOperations(ctx, client)

	// 集合操作示例
	fmt.Println("\n=== 集合操作示例 ===")
	setOperations(ctx, client)

	// 有序集合操作示例
	fmt.Println("\n=== 有序集合操作示例 ===")
	zsetOperations(ctx, client)

	// 计数器操作示例
	fmt.Println("\n=== 计数器操作示例 ===")
	counterOperations(ctx, client)

	// 管道操作示例
	fmt.Println("\n=== 管道操作示例 ===")
	pipelineOperations(ctx, client)

	fmt.Println("\n🎉 所有操作完成")
}

// 字符串操作示例
func stringOperations(ctx context.Context, client cache.Client) {
	// 设置字符串值
	err := client.Set(ctx, "user:1:name", "张三", 0)
	if err != nil {
		log.Printf("设置字符串失败: %v", err)
		return
	}
	fmt.Println("设置用户名: 张三")

	// 获取字符串值
	name, err := client.Get(ctx, "user:1:name")
	if err != nil {
		log.Printf("获取字符串失败: %v", err)
		return
	}
	fmt.Printf("获取用户名: %s\n", name)

	// 批量设置
	err = client.MSet(ctx, "user:1:email", "zhangsan@example.com", "user:1:age", "25")
	if err != nil {
		log.Printf("批量设置失败: %v", err)
		return
	}
	fmt.Println("批量设置用户信息")

	// 批量获取
	values, err := client.MGet(ctx, "user:1:name", "user:1:email", "user:1:age")
	if err != nil {
		log.Printf("批量获取失败: %v", err)
		return
	}
	fmt.Printf("批量获取结果: %v\n", values)
}

// 哈希表操作示例
func hashOperations(ctx context.Context, client cache.Client) {
	// 设置哈希表字段
	err := client.HSet(ctx, "user:2", "name", "李四")
	if err != nil {
		log.Printf("设置哈希字段失败: %v", err)
		return
	}

	err = client.HSet(ctx, "user:2", "email", "lisi@example.com")
	if err != nil {
		log.Printf("设置哈希字段失败: %v", err)
		return
	}

	err = client.HSet(ctx, "user:2", "age", "30")
	if err != nil {
		log.Printf("设置哈希字段失败: %v", err)
		return
	}
	fmt.Println("设置用户哈希信息")

	// 获取哈希表所有字段
	userInfo, err := client.HGetAll(ctx, "user:2")
	if err != nil {
		log.Printf("获取哈希表失败: %v", err)
		return
	}
	fmt.Printf("用户信息: %v\n", userInfo)

	// 获取单个字段
	name, err := client.HGet(ctx, "user:2", "name")
	if err != nil {
		log.Printf("获取哈希字段失败: %v", err)
		return
	}
	fmt.Printf("用户姓名: %s\n", name)
}

// 列表操作示例
func listOperations(ctx context.Context, client cache.Client) {
	// 从左侧推入元素
	_, err := client.LPush(ctx, "tasks", "任务1", "任务2", "任务3")
	if err != nil {
		log.Printf("推入列表失败: %v", err)
		return
	}
	fmt.Println("推入任务到列表")

	// 获取列表长度
	length, err := client.LLen(ctx, "tasks")
	if err != nil {
		log.Printf("获取列表长度失败: %v", err)
		return
	}
	fmt.Printf("任务列表长度: %d\n", length)

	// 获取列表范围
	tasks, err := client.LRange(ctx, "tasks", 0, -1)
	if err != nil {
		log.Printf("获取列表范围失败: %v", err)
		return
	}
	fmt.Printf("所有任务: %v\n", tasks)

	// 从右侧弹出元素
	task, err := client.RPop(ctx, "tasks")
	if err != nil {
		log.Printf("弹出任务失败: %v", err)
		return
	}
	fmt.Printf("弹出的任务: %s\n", task)
}

// 集合操作示例
func setOperations(ctx context.Context, client cache.Client) {
	// 向集合添加成员
	_, err := client.SAdd(ctx, "tags", "Go", "Redis", "缓存", "数据库")
	if err != nil {
		log.Printf("添加集合成员失败: %v", err)
		return
	}
	fmt.Println("添加标签到集合")

	// 获取集合所有成员
	tags, err := client.SMembers(ctx, "tags")
	if err != nil {
		log.Printf("获取集合成员失败: %v", err)
		return
	}
	fmt.Printf("所有标签: %v\n", tags)

	// 检查成员是否存在
	exists, err := client.SIsMember(ctx, "tags", "Go")
	if err != nil {
		log.Printf("检查成员失败: %v", err)
		return
	}
	fmt.Printf("Go标签是否存在: %t\n", exists)

	// 获取集合大小
	size, err := client.SCard(ctx, "tags")
	if err != nil {
		log.Printf("获取集合大小失败: %v", err)
		return
	}
	fmt.Printf("标签集合大小: %d\n", size)
}

// 有序集合操作示例
func zsetOperations(ctx context.Context, client cache.Client) {
	// 向有序集合添加成员
	members := []cache.ZMember{
		{Score: 100, Member: "用户A"},
		{Score: 85, Member: "用户B"},
		{Score: 92, Member: "用户C"},
		{Score: 78, Member: "用户D"},
	}

	_, err := client.ZAdd(ctx, "leaderboard", members...)
	if err != nil {
		log.Printf("添加有序集合成员失败: %v", err)
		return
	}
	fmt.Println("添加用户到排行榜")

	// 获取排行榜前3名（分数从高到低）
	top3, err := client.ZRevRange(ctx, "leaderboard", 0, 2)
	if err != nil {
		log.Printf("获取排行榜失败: %v", err)
		return
	}
	fmt.Printf("排行榜前3名: %v\n", top3)

	// 获取用户排名
	rank, err := client.ZRevRank(ctx, "leaderboard", "用户A")
	if err != nil {
		log.Printf("获取用户排名失败: %v", err)
		return
	}
	fmt.Printf("用户A的排名: %d\n", rank+1) // 排名从0开始，所以+1

	// 获取用户分数
	score, err := client.ZScore(ctx, "leaderboard", "用户A")
	if err != nil {
		log.Printf("获取用户分数失败: %v", err)
		return
	}
	fmt.Printf("用户A的分数: %.0f\n", score)
}

// 计数器操作示例
func counterOperations(ctx context.Context, client cache.Client) {
	// 递增计数器
	count, err := client.Incr(ctx, "page_views")
	if err != nil {
		log.Printf("递增计数器失败: %v", err)
		return
	}
	fmt.Printf("页面访问次数: %d\n", count)

	// 按指定值递增
	count, err = client.IncrBy(ctx, "page_views", 5)
	if err != nil {
		log.Printf("按值递增失败: %v", err)
		return
	}
	fmt.Printf("页面访问次数（+5）: %d\n", count)

	// 递减计数器
	count, err = client.Decr(ctx, "page_views")
	if err != nil {
		log.Printf("递减计数器失败: %v", err)
		return
	}
	fmt.Printf("页面访问次数（-1）: %d\n", count)
}

// 管道操作示例
func pipelineOperations(ctx context.Context, client cache.Client) {
	// 创建管道
	pipe := client.Pipeline()

	// 批量添加命令到管道
	pipe.Set(ctx, "batch:1", "值1", 0)
	pipe.Set(ctx, "batch:2", "值2", 0)
	pipe.Set(ctx, "batch:3", "值3", 0)
	pipe.Incr(ctx, "batch:counter")
	pipe.HSet(ctx, "batch:hash", "field1", "value1")

	// 执行管道
	results, err := pipe.Exec(ctx)
	if err != nil {
		log.Printf("执行管道失败: %v", err)
		return
	}

	fmt.Printf("管道执行结果数量: %d\n", len(results))
	fmt.Println("批量操作完成")

	// 关闭管道
	pipe.Close()
}