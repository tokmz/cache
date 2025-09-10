package unit

import (
	"context"
	"testing"
	"time"

	"cache"
	"github.com/stretchr/testify/assert"
)

// TestPipelineBasicOperations 测试管道基础操作
func TestPipelineBasicOperations(t *testing.T) {
	config := &cache.Config{
		Mode: cache.ModeSingle,
		Single: &cache.SingleConfig{
			Addr: "localhost:6379",
			DB:   0,
		},
		Common: cache.CommonConfig{
			PoolSize: 10,
		},
	}

	factory, err := cache.NewFactory(config)
	assert.NoError(t, err)

	client, err := factory.CreateClient()
	assert.NoError(t, err)
	assert.NotNil(t, client)

	ctx := context.Background()

	// 测试管道创建和基本操作
	t.Run("管道创建和基本操作", func(t *testing.T) {
		pipe := client.Pipeline()
		assert.NotNil(t, pipe)

		// 添加多个命令到管道
		setCmd1 := pipe.Set(ctx, "pipe:key1", "value1", 0)
		setCmd2 := pipe.Set(ctx, "pipe:key2", "value2", 0)
		getCmd1 := pipe.Get(ctx, "pipe:key1")
		getCmd2 := pipe.Get(ctx, "pipe:key2")

		// 执行管道
		results, err := pipe.Exec(ctx)
		assert.NoError(t, err)
		assert.Len(t, results, 4)

		// 验证命令结果
		assert.NoError(t, setCmd1.Err())
		assert.NoError(t, setCmd2.Err())
		assert.NoError(t, getCmd1.Err())
		assert.NoError(t, getCmd2.Err())
		assert.Equal(t, "value1", getCmd1.Val())
		assert.Equal(t, "value2", getCmd2.Val())

		// 关闭管道
		err = pipe.Close()
		assert.NoError(t, err)

		// 清理
		client.Del(ctx, "pipe:key1", "pipe:key2")
	})
}

// TestPipelineStringOperations 测试管道字符串操作
func TestPipelineStringOperations(t *testing.T) {
	config := &cache.Config{
		Mode: cache.ModeSingle,
		Single: &cache.SingleConfig{
			Addr: "localhost:6379",
			DB:   0,
		},
		Common: cache.CommonConfig{
			PoolSize: 10,
		},
	}

	factory, err := cache.NewFactory(config)
	assert.NoError(t, err)

	client, err := factory.CreateClient()
	assert.NoError(t, err)

	ctx := context.Background()

	t.Run("管道字符串操作", func(t *testing.T) {
		pipe := client.Pipeline()
		defer pipe.Close()

		// 批量设置字符串
		keys := []string{"pipe:str1", "pipe:str2", "pipe:str3"}
		values := []string{"value1", "value2", "value3"}

		setCmds := make([]*cache.StatusCmd, len(keys))
		for i, key := range keys {
			setCmds[i] = pipe.Set(ctx, key, values[i], 0)
		}

		// 批量获取字符串
		getCmds := make([]*cache.StringCmd, len(keys))
		for i, key := range keys {
			getCmds[i] = pipe.Get(ctx, key)
		}

		// 执行管道
		_, err := pipe.Exec(ctx)
		assert.NoError(t, err)

		// 验证设置结果
		for _, setCmd := range setCmds {
			assert.NoError(t, setCmd.Err())
		}

		// 验证获取结果
		for i, getCmd := range getCmds {
			assert.NoError(t, getCmd.Err())
			assert.Equal(t, values[i], getCmd.Val())
		}

		// 清理
		client.Del(ctx, keys...)
	})
}

// TestPipelineHashOperations 测试管道哈希操作
func TestPipelineHashOperations(t *testing.T) {
	config := &cache.Config{
		Mode: cache.ModeSingle,
		Single: &cache.SingleConfig{
			Addr: "localhost:6379",
			DB:   0,
		},
		Common: cache.CommonConfig{
			PoolSize: 10,
		},
	}

	factory, err := cache.NewFactory(config)
	assert.NoError(t, err)

	client, err := factory.CreateClient()
	assert.NoError(t, err)

	ctx := context.Background()

	t.Run("管道哈希操作", func(t *testing.T) {
		pipe := client.Pipeline()
		defer pipe.Close()

		hashKey := "pipe:hash"
		fields := []string{"field1", "field2", "field3"}
		values := []string{"value1", "value2", "value3"}

		// 批量设置哈希字段
		hsetCmds := make([]*cache.IntCmd, len(fields))
		for i, field := range fields {
			hsetCmds[i] = pipe.HSet(ctx, hashKey, field, values[i])
		}

		// 批量获取哈希字段
		hgetCmds := make([]*cache.StringCmd, len(fields))
		for i, field := range fields {
			hgetCmds[i] = pipe.HGet(ctx, hashKey, field)
		}

		// 执行管道
		_, err := pipe.Exec(ctx)
		assert.NoError(t, err)

		// 验证设置结果
		for _, hsetCmd := range hsetCmds {
			assert.NoError(t, hsetCmd.Err())
		}

		// 验证获取结果
		for i, hgetCmd := range hgetCmds {
			assert.NoError(t, hgetCmd.Err())
			assert.Equal(t, values[i], hgetCmd.Val())
		}

		// 清理
		client.Del(ctx, hashKey)
	})
}

// TestPipelineListOperations 测试管道列表操作
func TestPipelineListOperations(t *testing.T) {
	config := &cache.Config{
		Mode: cache.ModeSingle,
		Single: &cache.SingleConfig{
			Addr: "localhost:6379",
			DB:   0,
		},
		Common: cache.CommonConfig{
			PoolSize: 10,
		},
	}

	factory, err := cache.NewFactory(config)
	assert.NoError(t, err)

	client, err := factory.CreateClient()
	assert.NoError(t, err)

	ctx := context.Background()

	t.Run("管道列表操作", func(t *testing.T) {
		pipe := client.Pipeline()
		defer pipe.Close()

		listKey := "pipe:list"
		values := []string{"item1", "item2", "item3"}

		// 批量推入列表
		lpushCmds := make([]*cache.IntCmd, len(values))
		for i, value := range values {
			lpushCmds[i] = pipe.LPush(ctx, listKey, value)
		}

		// 弹出列表元素
		lpopCmd1 := pipe.LPop(ctx, listKey)
		lpopCmd2 := pipe.LPop(ctx, listKey)

		// 执行管道
		_, err := pipe.Exec(ctx)
		assert.NoError(t, err)

		// 验证推入结果
		for _, lpushCmd := range lpushCmds {
			assert.NoError(t, lpushCmd.Err())
			assert.True(t, lpushCmd.Val() > 0)
		}

		// 验证弹出结果
		assert.NoError(t, lpopCmd1.Err())
		assert.NoError(t, lpopCmd2.Err())
		assert.Contains(t, values, lpopCmd1.Val())
		assert.Contains(t, values, lpopCmd2.Val())

		// 清理
		client.Del(ctx, listKey)
	})
}

// TestPipelineSetOperations 测试管道集合操作
func TestPipelineSetOperations(t *testing.T) {
	config := &cache.Config{
		Mode: cache.ModeSingle,
		Single: &cache.SingleConfig{
			Addr: "localhost:6379",
			DB:   0,
		},
		Common: cache.CommonConfig{
			PoolSize: 10,
		},
	}

	factory, err := cache.NewFactory(config)
	assert.NoError(t, err)

	client, err := factory.CreateClient()
	assert.NoError(t, err)

	ctx := context.Background()

	t.Run("管道集合操作", func(t *testing.T) {
		pipe := client.Pipeline()
		defer pipe.Close()

		setKey := "pipe:set"
		members := []string{"member1", "member2", "member3"}

		// 批量添加集合成员
		saddCmds := make([]*cache.IntCmd, len(members))
		for i, member := range members {
			saddCmds[i] = pipe.SAdd(ctx, setKey, member)
		}

		// 获取集合所有成员
		smembersCmd := pipe.SMembers(ctx, setKey)

		// 执行管道
		_, err := pipe.Exec(ctx)
		assert.NoError(t, err)

		// 验证添加结果
		for _, saddCmd := range saddCmds {
			assert.NoError(t, saddCmd.Err())
			assert.Equal(t, int64(1), saddCmd.Val())
		}

		// 验证获取结果
		assert.NoError(t, smembersCmd.Err())
		resultMembers := smembersCmd.Val()
		assert.Len(t, resultMembers, len(members))
		for _, member := range members {
			assert.Contains(t, resultMembers, member)
		}

		// 清理
		client.Del(ctx, setKey)
	})
}

// TestPipelineKeyOperations 测试管道键操作
func TestPipelineKeyOperations(t *testing.T) {
	config := &cache.Config{
		Mode: cache.ModeSingle,
		Single: &cache.SingleConfig{
			Addr: "localhost:6379",
			DB:   0,
		},
		Common: cache.CommonConfig{
			PoolSize: 10,
		},
	}

	factory, err := cache.NewFactory(config)
	assert.NoError(t, err)

	client, err := factory.CreateClient()
	assert.NoError(t, err)

	ctx := context.Background()

	t.Run("管道键操作", func(t *testing.T) {
		pipe := client.Pipeline()
		defer pipe.Close()

		keys := []string{"pipe:key1", "pipe:key2", "pipe:key3"}

		// 设置键
		for _, key := range keys {
			pipe.Set(ctx, key, "value", 0)
		}

		// 检查键存在性
		existsCmds := make([]*cache.IntCmd, len(keys))
		for i, key := range keys {
			existsCmds[i] = pipe.Exists(ctx, key)
		}

		// 设置过期时间
		expireCmds := make([]*cache.BoolCmd, len(keys))
		for i, key := range keys {
			expireCmds[i] = pipe.Expire(ctx, key, 10*time.Second)
		}

		// 删除键
		delCmd := pipe.Del(ctx, keys...)

		// 执行管道
		_, err := pipe.Exec(ctx)
		assert.NoError(t, err)

		// 验证存在性检查
		for _, existsCmd := range existsCmds {
			assert.NoError(t, existsCmd.Err())
			assert.Equal(t, int64(1), existsCmd.Val())
		}

		// 验证过期时间设置
		for _, expireCmd := range expireCmds {
			assert.NoError(t, expireCmd.Err())
			assert.True(t, expireCmd.Val())
		}

		// 验证删除结果
		assert.NoError(t, delCmd.Err())
		assert.Equal(t, int64(len(keys)), delCmd.Val())
	})
}

// TestPipelineErrorHandling 测试管道错误处理
func TestPipelineErrorHandling(t *testing.T) {
	config := &cache.Config{
		Mode: cache.ModeSingle,
		Single: &cache.SingleConfig{
			Addr: "localhost:6379",
			DB:   0,
		},
		Common: cache.CommonConfig{
			PoolSize: 10,
		},
	}

	factory, err := cache.NewFactory(config)
	assert.NoError(t, err)

	client, err := factory.CreateClient()
	assert.NoError(t, err)

	ctx := context.Background()

	t.Run("管道错误处理", func(t *testing.T) {
		pipe := client.Pipeline()
		defer pipe.Close()

		// 设置一个字符串键
		setCmd := pipe.Set(ctx, "pipe:error:key", "string_value", 0)

		// 尝试对字符串键执行列表操作（会出错）
		lpushCmd := pipe.LPush(ctx, "pipe:error:key", "item")

		// 正常的获取操作
		getCmd := pipe.Get(ctx, "pipe:error:key")

		// 执行管道
		results, err := pipe.Exec(ctx)
		assert.NoError(t, err) // 管道执行本身不应该出错
		assert.Len(t, results, 3)

		// 验证各命令的结果
		assert.NoError(t, setCmd.Err())
		assert.Error(t, lpushCmd.Err()) // 这个操作应该出错
		assert.NoError(t, getCmd.Err())
		assert.Equal(t, "string_value", getCmd.Val())

		// 清理
		client.Del(ctx, "pipe:error:key")
	})
}

// TestPipelineDiscard 测试管道丢弃操作
func TestPipelineDiscard(t *testing.T) {
	config := &cache.Config{
		Mode: cache.ModeSingle,
		Single: &cache.SingleConfig{
			Addr: "localhost:6379",
			DB:   0,
		},
		Common: cache.CommonConfig{
			PoolSize: 10,
		},
	}

	factory, err := cache.NewFactory(config)
	assert.NoError(t, err)

	client, err := factory.CreateClient()
	assert.NoError(t, err)

	ctx := context.Background()

	t.Run("管道丢弃操作", func(t *testing.T) {
		pipe := client.Pipeline()
		defer pipe.Close()

		// 添加一些命令
		pipe.Set(ctx, "pipe:discard:key1", "value1", 0)
		pipe.Set(ctx, "pipe:discard:key2", "value2", 0)

		// 丢弃管道
		err := pipe.Discard()
		assert.NoError(t, err)

		// 验证键不存在（因为管道被丢弃了）
		count, err := client.Exists(ctx, "pipe:discard:key1", "pipe:discard:key2")
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})
}

// TestTxPipeline 测试事务管道
func TestTxPipeline(t *testing.T) {
	config := &cache.Config{
		Mode: cache.ModeSingle,
		Single: &cache.SingleConfig{
			Addr: "localhost:6379",
			DB:   0,
		},
		Common: cache.CommonConfig{
			PoolSize: 10,
		},
	}

	factory, err := cache.NewFactory(config)
	assert.NoError(t, err)

	client, err := factory.CreateClient()
	assert.NoError(t, err)

	ctx := context.Background()

	t.Run("事务管道操作", func(t *testing.T) {
		txPipe := client.TxPipeline()
		assert.NotNil(t, txPipe)
		defer txPipe.Close()

		// 在事务中添加多个命令
		setCmd1 := txPipe.Set(ctx, "tx:key1", "value1", 0)
		setCmd2 := txPipe.Set(ctx, "tx:key2", "value2", 0)
		incrCmd := txPipe.Incr(ctx, "tx:counter")

		// 执行事务
		results, err := txPipe.Exec(ctx)
		assert.NoError(t, err)
		assert.Len(t, results, 3)

		// 验证所有命令都成功执行
		assert.NoError(t, setCmd1.Err())
		assert.NoError(t, setCmd2.Err())
		assert.NoError(t, incrCmd.Err())
		assert.Equal(t, int64(1), incrCmd.Val())

		// 验证键确实被设置
		value1, err := client.Get(ctx, "tx:key1")
		assert.NoError(t, err)
		assert.Equal(t, "value1", value1)

		value2, err := client.Get(ctx, "tx:key2")
		assert.NoError(t, err)
		assert.Equal(t, "value2", value2)

		// 清理
		client.Del(ctx, "tx:key1", "tx:key2", "tx:counter")
	})
}

// BenchmarkPipelineOperations 管道操作基准测试
func BenchmarkPipelineOperations(b *testing.B) {
	config := &cache.Config{
		Mode: cache.ModeSingle,
		Single: &cache.SingleConfig{
			Addr: "localhost:6379",
			DB:   0,
		},
		Common: cache.CommonConfig{
			PoolSize: 10,
		},
	}

	factory, _ := cache.NewFactory(config)
	client, _ := factory.CreateClient()
	ctx := context.Background()

	b.Run("管道批量Set操作", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			pipe := client.Pipeline()
			for j := 0; j < 100; j++ {
				pipe.Set(ctx, "bench:pipe:set", "value", 0)
			}
			pipe.Exec(ctx)
			pipe.Close()
		}
	})

	b.Run("管道批量Get操作", func(b *testing.B) {
		// 预设一些键
		for i := 0; i < 100; i++ {
			client.Set(ctx, "bench:pipe:get", "value", 0)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			pipe := client.Pipeline()
			for j := 0; j < 100; j++ {
				pipe.Get(ctx, "bench:pipe:get")
			}
			pipe.Exec(ctx)
			pipe.Close()
		}
	})

	b.Run("事务管道操作", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			txPipe := client.TxPipeline()
			for j := 0; j < 10; j++ {
				txPipe.Set(ctx, "bench:tx:set", "value", 0)
				txPipe.Incr(ctx, "bench:tx:counter")
			}
			txPipe.Exec(ctx)
			txPipe.Close()
		}
	})
}

// BenchmarkPipelineVsIndividual 管道操作与单独操作性能对比
func BenchmarkPipelineVsIndividual(b *testing.B) {
	config := &cache.Config{
		Mode: cache.ModeSingle,
		Single: &cache.SingleConfig{
			Addr: "localhost:6379",
			DB:   0,
		},
		Common: cache.CommonConfig{
			PoolSize: 10,
		},
	}

	factory, _ := cache.NewFactory(config)
	client, _ := factory.CreateClient()
	ctx := context.Background()

	b.Run("单独操作", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for j := 0; j < 10; j++ {
				client.Set(ctx, "bench:individual", "value", 0)
				client.Get(ctx, "bench:individual")
			}
		}
	})

	b.Run("管道操作", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			pipe := client.Pipeline()
			for j := 0; j < 10; j++ {
				pipe.Set(ctx, "bench:pipeline", "value", 0)
				pipe.Get(ctx, "bench:pipeline")
			}
			pipe.Exec(ctx)
			pipe.Close()
		}
	})
}