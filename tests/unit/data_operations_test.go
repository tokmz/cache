package unit

import (
	"context"
	"testing"
	"time"

	"cache"

	"github.com/stretchr/testify/assert"
)

// TestStringOperations 测试字符串操作
func TestStringOperations(t *testing.T) {
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

	// 测试Set和Get操作
	t.Run("Set和Get操作", func(t *testing.T) {
		key := "test:string:basic"
		value := "hello world"

		// Set操作
		err := client.Set(ctx, key, value, 0)
		assert.NoError(t, err)

		// Get操作
		result, err := client.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, value, result)

		// 清理
		client.Del(ctx, key)
	})

	// 测试SetEX操作
	t.Run("SetEX操作", func(t *testing.T) {
		key := "test:string:setex"
		value := "expire value"
		expiration := 1 * time.Second

		// SetEX操作
		err := client.Set(ctx, key, value, expiration)
		assert.NoError(t, err)

		// 立即获取
		result, err := client.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, value, result)

		// 等待过期
		time.Sleep(2 * time.Second)
		result, err = client.Get(ctx, key)
		assert.Error(t, err) // 应该返回错误，因为键已过期
	})

	// 测试Incr和Decr操作
	t.Run("Incr和Decr操作", func(t *testing.T) {
		key := "test:string:counter"

		// 初始化计数器
		client.Set(ctx, key, "10", 0)

		// Incr操作
		result, err := client.Incr(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, int64(11), result)

		// IncrBy操作
		result, err = client.IncrBy(ctx, key, 5)
		assert.NoError(t, err)
		assert.Equal(t, int64(16), result)

		// Decr操作
		result, err = client.Decr(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, int64(15), result)

		// DecrBy操作
		result, err = client.DecrBy(ctx, key, 3)
		assert.NoError(t, err)
		assert.Equal(t, int64(12), result)

		// 清理
		client.Del(ctx, key)
	})
}

// TestHashOperations 测试哈希操作
func TestHashOperations(t *testing.T) {
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

	// 测试HSet和HGet操作
	t.Run("HSet和HGet操作", func(t *testing.T) {
		key := "test:hash:basic"
		field := "name"
		value := "John Doe"

		// HSet操作
		err := client.HSet(ctx, key, field, value)
		assert.NoError(t, err)

		// HGet操作
		hvalue, err := client.HGet(ctx, key, field)
		assert.NoError(t, err)
		assert.Equal(t, value, hvalue)

		// 清理
		client.Del(ctx, key)
	})

	// 测试HMSet和HMGet操作
	t.Run("HMSet和HMGet操作", func(t *testing.T) {
		key := "test:hash:multi"
		fields := map[string]interface{}{
			"name":  "Jane Doe",
			"age":   "30",
			"email": "jane@example.com",
		}

		// HMSet操作
		err := client.HMSet(ctx, key, fields)
		assert.NoError(t, err)

		// HMGet操作
		results, err := client.HMGet(ctx, key, "name", "age", "email")
		assert.NoError(t, err)
		assert.Len(t, results, 3)
		assert.Equal(t, "Jane Doe", results[0])
		assert.Equal(t, "30", results[1])
		assert.Equal(t, "jane@example.com", results[2])

		// 清理
		client.Del(ctx, key)
	})

	// 测试HGetAll操作
	t.Run("HGetAll操作", func(t *testing.T) {
		key := "test:hash:getall"
		fields := map[string]interface{}{
			"field1": "value1",
			"field2": "value2",
			"field3": "value3",
		}

		// 设置多个字段
		client.HMSet(ctx, key, fields)

		// HGetAll操作
		result, err := client.HGetAll(ctx, key)
		assert.NoError(t, err)
		assert.Len(t, result, 6) // 每个字段对应两个元素：字段名和值

		// 清理
		client.Del(ctx, key)
	})

	// 测试HDel操作
	t.Run("HDel操作", func(t *testing.T) {
		key := "test:hash:del"

		// 设置字段
		client.HSet(ctx, key, "field1", "value1")
		client.HSet(ctx, key, "field2", "value2")

		// HDel操作
		result, err := client.HDel(ctx, key, "field1")
		assert.NoError(t, err)
		assert.Equal(t, int64(1), result)

		// 验证字段已删除
		_, err = client.HGet(ctx, key, "field1")
		assert.Error(t, err)

		// 验证其他字段仍存在
		value, err := client.HGet(ctx, key, "field2")
		assert.NoError(t, err)
		assert.Equal(t, "value2", value)

		// 清理
		client.Del(ctx, key)
	})
}

// TestListOperations 测试列表操作
func TestListOperations(t *testing.T) {
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

	// 测试LPush和LPop操作
	t.Run("LPush和LPop操作", func(t *testing.T) {
		key := "test:list:lpush"

		// LPush操作
		result, err := client.LPush(ctx, key, "item1", "item2", "item3")
		assert.NoError(t, err)
		assert.Equal(t, int64(3), result)

		// LPop操作
		value, err := client.LPop(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, "item3", value) // LIFO顺序

		// 清理
		client.Del(ctx, key)
	})

	// 测试RPush和RPop操作
	t.Run("RPush和RPop操作", func(t *testing.T) {
		key := "test:list:rpush"

		// RPush操作
		result, err := client.RPush(ctx, key, "item1", "item2", "item3")
		assert.NoError(t, err)
		assert.Equal(t, int64(3), result)

		// RPop操作
		value, err := client.RPop(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, "item3", value) // LIFO顺序

		// 清理
		client.Del(ctx, key)
	})

	// 测试LRange操作
	t.Run("LRange操作", func(t *testing.T) {
		key := "test:list:range"

		// 添加元素
		client.RPush(ctx, key, "item1", "item2", "item3", "item4", "item5")

		// LRange操作
		result, err := client.LRange(ctx, key, 0, 2)
		assert.NoError(t, err)
		assert.Len(t, result, 3)
		assert.Equal(t, "item1", result[0])
		assert.Equal(t, "item2", result[1])
		assert.Equal(t, "item3", result[2])

		// 获取所有元素
		allItems, err := client.LRange(ctx, key, 0, -1)
		assert.NoError(t, err)
		assert.Len(t, allItems, 5)

		// 清理
		client.Del(ctx, key)
	})

	// 测试LLen操作
	t.Run("LLen操作", func(t *testing.T) {
		key := "test:list:len"

		// 添加元素
		client.RPush(ctx, key, "item1", "item2", "item3")

		// LLen操作
		length, err := client.LLen(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), length)

		// 清理
		client.Del(ctx, key)
	})
}

// TestSetOperations 测试集合操作
func TestSetOperations(t *testing.T) {
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

	// 测试SAdd和SMembers操作
	t.Run("SAdd和SMembers操作", func(t *testing.T) {
		key := "test:set:basic"

		// SAdd操作
		result, err := client.SAdd(ctx, key, "member1", "member2", "member3")
		assert.NoError(t, err)
		assert.Equal(t, int64(3), result)

		// SMembers操作
		members, err := client.SMembers(ctx, key)
		assert.NoError(t, err)
		assert.Len(t, members, 3)
		assert.Contains(t, members, "member1")
		assert.Contains(t, members, "member2")
		assert.Contains(t, members, "member3")

		// 清理
		client.Del(ctx, key)
	})

	// 测试SIsMember操作
	t.Run("SIsMember操作", func(t *testing.T) {
		key := "test:set:ismember"

		// 添加成员
		client.SAdd(ctx, key, "member1", "member2")

		// SIsMember操作
		isMember, err := client.SIsMember(ctx, key, "member1")
		assert.NoError(t, err)
		assert.True(t, isMember)

		isMember, err = client.SIsMember(ctx, key, "member3")
		assert.NoError(t, err)
		assert.False(t, isMember)

		// 清理
		client.Del(ctx, key)
	})

	// 测试SRem操作
	t.Run("SRem操作", func(t *testing.T) {
		key := "test:set:rem"

		// 添加成员
		client.SAdd(ctx, key, "member1", "member2", "member3")

		// SRem操作
		result, err := client.SRem(ctx, key, "member2")
		assert.NoError(t, err)
		assert.Equal(t, int64(1), result)

		// 验证成员已删除
		isMember, err := client.SIsMember(ctx, key, "member2")
		assert.NoError(t, err)
		assert.False(t, isMember)

		// 清理
		client.Del(ctx, key)
	})

	// 测试SCard操作
	t.Run("SCard操作", func(t *testing.T) {
		key := "test:set:card"

		// 添加成员
		client.SAdd(ctx, key, "member1", "member2", "member3")

		// SCard操作
		count, err := client.SCard(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), count)

		// 清理
		client.Del(ctx, key)
	})
}

// TestKeyOperations 测试键操作
func TestKeyOperations(t *testing.T) {
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

	// 测试Exists操作
	t.Run("Exists操作", func(t *testing.T) {
		key := "test:key:exists"

		// 键不存在
		count, err := client.Exists(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)

		// 设置键
		client.Set(ctx, key, "value", 0)

		// 键存在
		count, err = client.Exists(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)

		// 清理
		client.Del(ctx, key)
	})

	// 测试Del操作
	t.Run("Del操作", func(t *testing.T) {
		key1 := "test:key:del1"
		key2 := "test:key:del2"

		// 设置键
		client.Set(ctx, key1, "value1", 0)
		client.Set(ctx, key2, "value2", 0)

		// Del操作
		count, err := client.Del(ctx, key1, key2)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), count)

		// 验证键已删除
		count, err = client.Exists(ctx, key1)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	// 测试Expire操作
	t.Run("Expire操作", func(t *testing.T) {
		key := "test:key:expire"

		// 设置键
		client.Set(ctx, key, "value", 0)

		// Expire操作
		result, err := client.Expire(ctx, key, 1*time.Second)
		assert.NoError(t, err)
		assert.True(t, result)

		// 立即检查键存在
		count, err := client.Exists(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)

		// 等待过期
		time.Sleep(2 * time.Second)
		count, err = client.Exists(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	// 测试TTL操作
	t.Run("TTL操作", func(t *testing.T) {
		key := "test:key:ttl"

		// 设置带过期时间的键
		client.Set(ctx, key, "value", 10*time.Second)

		// TTL操作
		ttl, err := client.TTL(ctx, key)
		assert.NoError(t, err)
		assert.True(t, ttl > 0 && ttl <= 10*time.Second)

		// 清理
		client.Del(ctx, key)
	})
}

// BenchmarkStringOperations 字符串操作基准测试
func BenchmarkStringOperations(b *testing.B) {
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

	b.Run("Set操作", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			client.Set(ctx, "benchmark:set", "value", 0)
		}
	})

	b.Run("Get操作", func(b *testing.B) {
		client.Set(ctx, "benchmark:get", "value", 0)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			client.Get(ctx, "benchmark:get")
		}
	})

	b.Run("Incr操作", func(b *testing.B) {
		client.Set(ctx, "benchmark:incr", "0", 0)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			client.Incr(ctx, "benchmark:incr")
		}
	})
}

// BenchmarkHashOperations 哈希操作基准测试
func BenchmarkHashOperations(b *testing.B) {
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

	b.Run("HSet操作", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			client.HSet(ctx, "benchmark:hset", "field", "value")
		}
	})

	b.Run("HGet操作", func(b *testing.B) {
		client.HSet(ctx, "benchmark:hget", "field", "value")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			client.HGet(ctx, "benchmark:hget", "field")
		}
	})
}
