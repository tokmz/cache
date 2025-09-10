package integration

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"cache"
	"github.com/stretchr/testify/assert"
)

// 获取测试环境的Redis地址
func getRedisAddr() string {
	addr := os.Getenv("TEST_REDIS_ADDR")
	if addr == "" {
		return "localhost:6379" // 默认地址
	}
	return addr
}

// 获取测试环境的Redis集群地址
func getRedisClusterAddrs() []string {
	addrsEnv := os.Getenv("TEST_REDIS_CLUSTER_ADDRS")
	if addrsEnv == "" {
		return []string{"localhost:7000", "localhost:7001", "localhost:7002"} // 默认地址
	}
	// 实际项目中可以解析环境变量中的地址列表
	return []string{addrsEnv}
}

// 获取测试环境的Redis哨兵地址
func getRedisSentinelAddrs() []string {
	addrsEnv := os.Getenv("TEST_REDIS_SENTINEL_ADDRS")
	if addrsEnv == "" {
		return []string{"localhost:26379", "localhost:26380", "localhost:26381"} // 默认地址
	}
	// 实际项目中可以解析环境变量中的地址列表
	return []string{addrsEnv}
}

// 获取测试环境的Redis哨兵主节点名称
func getRedisSentinelMasterName() string {
	masterName := os.Getenv("TEST_REDIS_SENTINEL_MASTER")
	if masterName == "" {
		return "mymaster" // 默认主节点名称
	}
	return masterName
}

// TestSingleClientIntegration 单机模式客户端集成测试
func TestSingleClientIntegration(t *testing.T) {
	// 如果设置了SKIP_INTEGRATION_TESTS环境变量，则跳过集成测试
	if os.Getenv("SKIP_INTEGRATION_TESTS") != "" {
		t.Skip("Skipping integration tests")
	}

	// 创建单机模式配置
	config := &cache.Config{
		Mode: cache.ModeSingle,
		Single: &cache.SingleConfig{
			Addr: getRedisAddr(),
			DB:   0,
		},
		Common: cache.CommonConfig{
			Password:    "", // 根据实际环境设置密码
			PoolSize:    10,
			DialTimeout: 5 * time.Second,
		},
	}

	// 创建工厂和客户端
	factory, err := cache.NewFactory(config)
	assert.NoError(t, err)

	client, err := factory.CreateClient()
	assert.NoError(t, err)
	defer client.Close()

	// 测试连接
	ctx := context.Background()
	err = client.Ping(ctx)
	assert.NoError(t, err, "Failed to connect to Redis server")

	// 测试基本操作
	testKey := "integration:test:single"
	testValue := "integration-test-value"

	// 设置值
	err = client.Set(ctx, testKey, testValue, 10*time.Second)
	assert.NoError(t, err)

	// 获取值
	value, err := client.Get(ctx, testKey)
	assert.NoError(t, err)
	assert.Equal(t, testValue, value)

	// 删除值
	delCount, err := client.Del(ctx, testKey)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), delCount)

	// 验证值已删除
	_, err = client.Get(ctx, testKey)
	assert.Error(t, err, "Key should not exist after deletion")
}

// TestClusterClientIntegration 集群模式客户端集成测试
func TestClusterClientIntegration(t *testing.T) {
	// 如果设置了SKIP_INTEGRATION_TESTS或SKIP_CLUSTER_TESTS环境变量，则跳过测试
	if os.Getenv("SKIP_INTEGRATION_TESTS") != "" || os.Getenv("SKIP_CLUSTER_TESTS") != "" {
		t.Skip("Skipping cluster integration tests")
	}

	// 创建集群模式配置
	config := &cache.Config{
		Mode: cache.ModeCluster,
		Cluster: &cache.ClusterConfig{
			Addrs:       getRedisClusterAddrs(),
			MaxRedirects: 3,
		},
		Common: cache.CommonConfig{
			Password:    "", // 根据实际环境设置密码
			PoolSize:    10,
			DialTimeout: 5 * time.Second,
		},
	}

	// 创建工厂和客户端
	factory, err := cache.NewFactory(config)
	if err != nil {
		t.Skip(fmt.Sprintf("Skipping cluster test due to factory creation error: %v", err))
		return
	}

	client, err := factory.CreateClient()
	if err != nil {
		t.Skip(fmt.Sprintf("Skipping cluster test due to client creation error: %v", err))
		return
	}
	defer client.Close()

	// 测试连接
	ctx := context.Background()
	err = client.Ping(ctx)
	if err != nil {
		t.Skip(fmt.Sprintf("Skipping cluster test due to connection error: %v", err))
		return
	}

	// 测试基本操作
	testKey := "integration:test:cluster"
	testValue := "integration-test-value"

	// 设置值
	err = client.Set(ctx, testKey, testValue, 10*time.Second)
	assert.NoError(t, err)

	// 获取值
	value, err := client.Get(ctx, testKey)
	assert.NoError(t, err)
	assert.Equal(t, testValue, value)

	// 删除值
	delCount, err := client.Del(ctx, testKey)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), delCount)

	// 验证值已删除
	_, err = client.Get(ctx, testKey)
	assert.Error(t, err, "Key should not exist after deletion")
}

// TestSentinelClientIntegration 哨兵模式客户端集成测试
func TestSentinelClientIntegration(t *testing.T) {
	// 如果设置了SKIP_INTEGRATION_TESTS或SKIP_SENTINEL_TESTS环境变量，则跳过测试
	if os.Getenv("SKIP_INTEGRATION_TESTS") != "" || os.Getenv("SKIP_SENTINEL_TESTS") != "" {
		t.Skip("Skipping sentinel integration tests")
	}

	// 创建哨兵模式配置
	config := &cache.Config{
		Mode: cache.ModeSentinel,
		Sentinel: &cache.SentinelConfig{
			MasterName: getRedisSentinelMasterName(),
			Addrs:      getRedisSentinelAddrs(),
			DB:         0,
		},
		Common: cache.CommonConfig{
			Password:    "", // 根据实际环境设置密码
			PoolSize:    10,
			DialTimeout: 5 * time.Second,
		},
	}

	// 创建工厂和客户端
	factory, err := cache.NewFactory(config)
	if err != nil {
		t.Skip(fmt.Sprintf("Skipping sentinel test due to factory creation error: %v", err))
		return
	}

	client, err := factory.CreateClient()
	if err != nil {
		t.Skip(fmt.Sprintf("Skipping sentinel test due to client creation error: %v", err))
		return
	}
	defer client.Close()

	// 测试连接
	ctx := context.Background()
	err = client.Ping(ctx)
	if err != nil {
		t.Skip(fmt.Sprintf("Skipping sentinel test due to connection error: %v", err))
		return
	}

	// 测试基本操作
	testKey := "integration:test:sentinel"
	testValue := "integration-test-value"

	// 设置值
	err = client.Set(ctx, testKey, testValue, 10*time.Second)
	assert.NoError(t, err)

	// 获取值
	value, err := client.Get(ctx, testKey)
	assert.NoError(t, err)
	assert.Equal(t, testValue, value)

	// 删除值
	delCount, err := client.Del(ctx, testKey)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), delCount)

	// 验证值已删除
	_, err = client.Get(ctx, testKey)
	assert.Error(t, err, "Key should not exist after deletion")
}

// TestComplexDataOperationsIntegration 复杂数据操作集成测试
func TestComplexDataOperationsIntegration(t *testing.T) {
	// 如果设置了SKIP_INTEGRATION_TESTS环境变量，则跳过集成测试
	if os.Getenv("SKIP_INTEGRATION_TESTS") != "" {
		t.Skip("Skipping integration tests")
	}

	// 创建单机模式配置
	config := &cache.Config{
		Mode: cache.ModeSingle,
		Single: &cache.SingleConfig{
			Addr: getRedisAddr(),
			DB:   0,
		},
		Common: cache.CommonConfig{
			Password:    "", // 根据实际环境设置密码
			PoolSize:    10,
			DialTimeout: 5 * time.Second,
		},
	}

	// 创建工厂和客户端
	factory, err := cache.NewFactory(config)
	assert.NoError(t, err)

	client, err := factory.CreateClient()
	assert.NoError(t, err)
	defer client.Close()

	ctx := context.Background()

	// 测试哈希表操作
	t.Run("Hash Operations", func(t *testing.T) {
		hashKey := "integration:test:hash"

		// 设置哈希字段
		err := client.HSet(ctx, hashKey, "field1", "value1")
		assert.NoError(t, err)

		err = client.HSet(ctx, hashKey, "field2", "value2")
		assert.NoError(t, err)

		// 获取哈希字段
		value, err := client.HGet(ctx, hashKey, "field1")
		assert.NoError(t, err)
		assert.Equal(t, "value1", value)

		// 获取所有哈希字段
		allFields, err := client.HGetAll(ctx, hashKey)
		assert.NoError(t, err)
		assert.Contains(t, allFields, "field1")
		assert.Contains(t, allFields, "value1")
		assert.Contains(t, allFields, "field2")
		assert.Contains(t, allFields, "value2")

		// 删除哈希字段
		delCount, err := client.HDel(ctx, hashKey, "field1")
		assert.NoError(t, err)
		assert.Equal(t, int64(1), delCount)

		// 清理
		client.Del(ctx, hashKey)
	})

	// 测试列表操作
	t.Run("List Operations", func(t *testing.T) {
		listKey := "integration:test:list"

		// 推入列表元素
		pushCount, err := client.LPush(ctx, listKey, "item1", "item2", "item3")
		assert.NoError(t, err)
		assert.Equal(t, int64(3), pushCount)

		// 获取列表长度
		length, err := client.LLen(ctx, listKey)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), length)

		// 获取列表范围
		items, err := client.LRange(ctx, listKey, 0, -1)
		assert.NoError(t, err)
		assert.Len(t, items, 3)
		assert.Contains(t, items, "item1")
		assert.Contains(t, items, "item2")
		assert.Contains(t, items, "item3")

		// 弹出列表元素
		item, err := client.LPop(ctx, listKey)
		assert.NoError(t, err)
		assert.Contains(t, []string{"item1", "item2", "item3"}, item)

		// 清理
		client.Del(ctx, listKey)
	})

	// 测试集合操作
	t.Run("Set Operations", func(t *testing.T) {
		setKey := "integration:test:set"

		// 添加集合成员
		addCount, err := client.SAdd(ctx, setKey, "member1", "member2", "member3")
		assert.NoError(t, err)
		assert.Equal(t, int64(3), addCount)

		// 获取集合成员数
		count, err := client.SCard(ctx, setKey)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), count)

		// 检查成员是否存在
		isMember, err := client.SIsMember(ctx, setKey, "member1")
		assert.NoError(t, err)
		assert.True(t, isMember)

		// 获取所有成员
		members, err := client.SMembers(ctx, setKey)
		assert.NoError(t, err)
		assert.Len(t, members, 3)
		assert.Contains(t, members, "member1")
		assert.Contains(t, members, "member2")
		assert.Contains(t, members, "member3")

		// 清理
		client.Del(ctx, setKey)
	})
}

// TestPipelineIntegration 管道操作集成测试
func TestPipelineIntegration(t *testing.T) {
	// 如果设置了SKIP_INTEGRATION_TESTS环境变量，则跳过集成测试
	if os.Getenv("SKIP_INTEGRATION_TESTS") != "" {
		t.Skip("Skipping integration tests")
	}

	// 创建单机模式配置
	config := &cache.Config{
		Mode: cache.ModeSingle,
		Single: &cache.SingleConfig{
			Addr: getRedisAddr(),
			DB:   0,
		},
		Common: cache.CommonConfig{
			Password:    "", // 根据实际环境设置密码
			PoolSize:    10,
			DialTimeout: 5 * time.Second,
		},
	}

	// 创建工厂和客户端
	factory, err := cache.NewFactory(config)
	assert.NoError(t, err)

	client, err := factory.CreateClient()
	assert.NoError(t, err)
	defer client.Close()

	ctx := context.Background()

	// 测试管道操作
	t.Run("Pipeline Operations", func(t *testing.T) {
		// 创建管道
		pipe := client.Pipeline()
		defer pipe.Close()

		// 添加多个命令
		setCmd := pipe.Set(ctx, "integration:pipe:key1", "value1", 0)
		incrCmd := pipe.Incr(ctx, "integration:pipe:counter")
		getCmd := pipe.Get(ctx, "integration:pipe:key1")

		// 执行管道
		_, err := pipe.Exec(ctx)
		assert.NoError(t, err)

		// 验证命令结果
		assert.NoError(t, setCmd.Err())
		assert.NoError(t, incrCmd.Err())
		assert.NoError(t, getCmd.Err())
		assert.Equal(t, "value1", getCmd.Val())
		assert.Equal(t, int64(1), incrCmd.Val())

		// 清理
		client.Del(ctx, "integration:pipe:key1", "integration:pipe:counter")
	})

	// 测试事务管道
	t.Run("Transaction Pipeline", func(t *testing.T) {
		// 创建事务管道
		txPipe := client.TxPipeline()
		defer txPipe.Close()

		// 添加多个命令
		setCmd := txPipe.Set(ctx, "integration:tx:key1", "value1", 0)
		incrCmd := txPipe.Incr(ctx, "integration:tx:counter")
		getCmd := txPipe.Get(ctx, "integration:tx:key1")

		// 执行事务
		_, err := txPipe.Exec(ctx)
		assert.NoError(t, err)

		// 验证命令结果
		assert.NoError(t, setCmd.Err())
		assert.NoError(t, incrCmd.Err())
		assert.NoError(t, getCmd.Err())
		assert.Equal(t, "value1", getCmd.Val())
		assert.Equal(t, int64(1), incrCmd.Val())

		// 清理
		client.Del(ctx, "integration:tx:key1", "integration:tx:counter")
	})
}

// TestErrorHandlingIntegration 错误处理集成测试
func TestErrorHandlingIntegration(t *testing.T) {
	// 如果设置了SKIP_INTEGRATION_TESTS环境变量，则跳过集成测试
	if os.Getenv("SKIP_INTEGRATION_TESTS") != "" {
		t.Skip("Skipping integration tests")
	}

	// 创建单机模式配置
	config := &cache.Config{
		Mode: cache.ModeSingle,
		Single: &cache.SingleConfig{
			Addr: getRedisAddr(),
			DB:   0,
		},
		Common: cache.CommonConfig{
			Password:    "", // 根据实际环境设置密码
			PoolSize:    10,
			DialTimeout: 5 * time.Second,
		},
	}

	// 创建工厂和客户端
	factory, err := cache.NewFactory(config)
	assert.NoError(t, err)

	client, err := factory.CreateClient()
	assert.NoError(t, err)
	defer client.Close()

	ctx := context.Background()

	// 测试不存在的键
	t.Run("Non-existent Key", func(t *testing.T) {
		_, err := client.Get(ctx, "integration:nonexistent:key")
		assert.Error(t, err, "Should return error for non-existent key")
	})

	// 测试类型不匹配错误
	t.Run("Type Mismatch", func(t *testing.T) {
		// 设置字符串
		err := client.Set(ctx, "integration:error:type", "string value", 0)
		assert.NoError(t, err)

		// 尝试对字符串键执行列表操作
		_, err = client.LPush(ctx, "integration:error:type", "item")
		assert.Error(t, err, "Should return error for type mismatch")

		// 清理
		client.Del(ctx, "integration:error:type")
	})

	// 测试管道错误处理
	t.Run("Pipeline Error Handling", func(t *testing.T) {
		pipe := client.Pipeline()
		defer pipe.Close()

		// 设置字符串
		setCmd := pipe.Set(ctx, "integration:error:pipe", "string value", 0)

		// 尝试对字符串键执行列表操作
		lpushCmd := pipe.LPush(ctx, "integration:error:pipe", "item")

		// 执行管道
		_, err := pipe.Exec(ctx)
		assert.NoError(t, err, "Pipeline execution itself should not fail")

		// 验证命令结果
		assert.NoError(t, setCmd.Err())
		assert.Error(t, lpushCmd.Err(), "LPush command should fail with type error")

		// 清理
		client.Del(ctx, "integration:error:pipe")
	})
}