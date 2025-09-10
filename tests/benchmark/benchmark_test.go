package benchmark

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"cache"
)

// 基准测试配置
var (
	benchmarkConfig = &cache.Config{
		Mode: cache.ModeSingle,
		Single: &cache.SingleConfig{
			Addr: "localhost:6379",
			DB:   0,
		},
		Common: cache.CommonConfig{
			Password:    "",
			PoolSize:    100,
			DialTimeout: 5 * time.Second,
		},
	}
	benchmarkClient cache.Client
	benchmarkCtx    = context.Background()
)

// 初始化基准测试客户端
func init() {
	factory, err := cache.NewFactory(benchmarkConfig)
	if err != nil {
		panic(fmt.Sprintf("Failed to create factory: %v", err))
	}

	client, err := factory.CreateClient()
	if err != nil {
		panic(fmt.Sprintf("Failed to create client: %v", err))
	}

	benchmarkClient = client

	// 测试连接
	if err := benchmarkClient.Ping(benchmarkCtx); err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}
}

// BenchmarkStringOperations 字符串操作基准测试
func BenchmarkStringOperations(b *testing.B) {
	b.Run("Set", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("benchmark:string:set:%d", i)
			value := fmt.Sprintf("value_%d", i)
			err := benchmarkClient.Set(benchmarkCtx, key, value, time.Hour)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Get", func(b *testing.B) {
		// 预设数据
		for i := 0; i < 1000; i++ {
			key := fmt.Sprintf("benchmark:string:get:%d", i)
			value := fmt.Sprintf("value_%d", i)
			benchmarkClient.Set(benchmarkCtx, key, value, time.Hour)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("benchmark:string:get:%d", i%1000)
			_, err := benchmarkClient.Get(benchmarkCtx, key)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Incr", func(b *testing.B) {
		key := "benchmark:string:incr"
		benchmarkClient.Set(benchmarkCtx, key, "0", time.Hour)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := benchmarkClient.Incr(benchmarkCtx, key)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Del", func(b *testing.B) {
		// 预设数据
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("benchmark:string:del:%d", i)
			value := fmt.Sprintf("value_%d", i)
			benchmarkClient.Set(benchmarkCtx, key, value, time.Hour)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("benchmark:string:del:%d", i)
			_, err := benchmarkClient.Del(benchmarkCtx, key)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkHashOperations 哈希操作基准测试
func BenchmarkHashOperations(b *testing.B) {
	b.Run("HSet", func(b *testing.B) {
		hashKey := "benchmark:hash:hset"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			field := fmt.Sprintf("field_%d", i)
			value := fmt.Sprintf("value_%d", i)
			err := benchmarkClient.HSet(benchmarkCtx, hashKey, field, value)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("HGet", func(b *testing.B) {
		hashKey := "benchmark:hash:hget"
		// 预设数据
		for i := 0; i < 1000; i++ {
			field := fmt.Sprintf("field_%d", i)
			value := fmt.Sprintf("value_%d", i)
			benchmarkClient.HSet(benchmarkCtx, hashKey, field, value)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			field := fmt.Sprintf("field_%d", i%1000)
			_, err := benchmarkClient.HGet(benchmarkCtx, hashKey, field)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("HGetAll", func(b *testing.B) {
		hashKey := "benchmark:hash:hgetall"
		// 预设数据
		for i := 0; i < 100; i++ {
			field := fmt.Sprintf("field_%d", i)
			value := fmt.Sprintf("value_%d", i)
			benchmarkClient.HSet(benchmarkCtx, hashKey, field, value)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := benchmarkClient.HGetAll(benchmarkCtx, hashKey)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("HDel", func(b *testing.B) {
		hashKey := "benchmark:hash:hdel"
		// 预设数据
		for i := 0; i < b.N; i++ {
			field := fmt.Sprintf("field_%d", i)
			value := fmt.Sprintf("value_%d", i)
			benchmarkClient.HSet(benchmarkCtx, hashKey, field, value)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			field := fmt.Sprintf("field_%d", i)
			_, err := benchmarkClient.HDel(benchmarkCtx, hashKey, field)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkListOperations 列表操作基准测试
func BenchmarkListOperations(b *testing.B) {
	b.Run("LPush", func(b *testing.B) {
		listKey := "benchmark:list:lpush"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			value := fmt.Sprintf("item_%d", i)
			_, err := benchmarkClient.LPush(benchmarkCtx, listKey, value)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("RPush", func(b *testing.B) {
		listKey := "benchmark:list:rpush"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			value := fmt.Sprintf("item_%d", i)
			_, err := benchmarkClient.RPush(benchmarkCtx, listKey, value)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("LPop", func(b *testing.B) {
		listKey := "benchmark:list:lpop"
		// 预设数据
		for i := 0; i < b.N; i++ {
			value := fmt.Sprintf("item_%d", i)
			benchmarkClient.LPush(benchmarkCtx, listKey, value)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := benchmarkClient.LPop(benchmarkCtx, listKey)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("RPop", func(b *testing.B) {
		listKey := "benchmark:list:rpop"
		// 预设数据
		for i := 0; i < b.N; i++ {
			value := fmt.Sprintf("item_%d", i)
			benchmarkClient.RPush(benchmarkCtx, listKey, value)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := benchmarkClient.RPop(benchmarkCtx, listKey)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("LLen", func(b *testing.B) {
		listKey := "benchmark:list:llen"
		// 预设数据
		for i := 0; i < 1000; i++ {
			value := fmt.Sprintf("item_%d", i)
			benchmarkClient.LPush(benchmarkCtx, listKey, value)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := benchmarkClient.LLen(benchmarkCtx, listKey)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("LRange", func(b *testing.B) {
		listKey := "benchmark:list:lrange"
		// 预设数据
		for i := 0; i < 1000; i++ {
			value := fmt.Sprintf("item_%d", i)
			benchmarkClient.LPush(benchmarkCtx, listKey, value)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := benchmarkClient.LRange(benchmarkCtx, listKey, 0, 99)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkSetOperations 集合操作基准测试
func BenchmarkSetOperations(b *testing.B) {
	b.Run("SAdd", func(b *testing.B) {
		setKey := "benchmark:set:sadd"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			member := fmt.Sprintf("member_%d", i)
			_, err := benchmarkClient.SAdd(benchmarkCtx, setKey, member)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("SIsMember", func(b *testing.B) {
		setKey := "benchmark:set:sismember"
		// 预设数据
		for i := 0; i < 1000; i++ {
			member := fmt.Sprintf("member_%d", i)
			benchmarkClient.SAdd(benchmarkCtx, setKey, member)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			member := fmt.Sprintf("member_%d", i%1000)
			_, err := benchmarkClient.SIsMember(benchmarkCtx, setKey, member)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("SCard", func(b *testing.B) {
		setKey := "benchmark:set:scard"
		// 预设数据
		for i := 0; i < 1000; i++ {
			member := fmt.Sprintf("member_%d", i)
			benchmarkClient.SAdd(benchmarkCtx, setKey, member)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := benchmarkClient.SCard(benchmarkCtx, setKey)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("SMembers", func(b *testing.B) {
		setKey := "benchmark:set:smembers"
		// 预设数据
		for i := 0; i < 100; i++ {
			member := fmt.Sprintf("member_%d", i)
			benchmarkClient.SAdd(benchmarkCtx, setKey, member)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := benchmarkClient.SMembers(benchmarkCtx, setKey)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("SRem", func(b *testing.B) {
		setKey := "benchmark:set:srem"
		// 预设数据
		for i := 0; i < b.N; i++ {
			member := fmt.Sprintf("member_%d", i)
			benchmarkClient.SAdd(benchmarkCtx, setKey, member)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			member := fmt.Sprintf("member_%d", i)
			_, err := benchmarkClient.SRem(benchmarkCtx, setKey, member)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkPipelineOperations 管道操作基准测试
func BenchmarkPipelineOperations(b *testing.B) {
	b.Run("Pipeline_10_Commands", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			pipe := benchmarkClient.Pipeline()

			// 添加10个命令
			for j := 0; j < 10; j++ {
				key := fmt.Sprintf("benchmark:pipe:10:%d:%d", i, j)
				value := fmt.Sprintf("value_%d_%d", i, j)
				pipe.Set(benchmarkCtx, key, value, time.Hour)
			}

			_, err := pipe.Exec(benchmarkCtx)
			if err != nil {
				b.Fatal(err)
			}
			pipe.Close()
		}
	})

	b.Run("Pipeline_100_Commands", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			pipe := benchmarkClient.Pipeline()

			// 添加100个命令
			for j := 0; j < 100; j++ {
				key := fmt.Sprintf("benchmark:pipe:100:%d:%d", i, j)
				value := fmt.Sprintf("value_%d_%d", i, j)
				pipe.Set(benchmarkCtx, key, value, time.Hour)
			}

			_, err := pipe.Exec(benchmarkCtx)
			if err != nil {
				b.Fatal(err)
			}
			pipe.Close()
		}
	})

	b.Run("TxPipeline_10_Commands", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			txPipe := benchmarkClient.TxPipeline()

			// 添加10个命令
			for j := 0; j < 10; j++ {
				key := fmt.Sprintf("benchmark:txpipe:10:%d:%d", i, j)
				value := fmt.Sprintf("value_%d_%d", i, j)
				txPipe.Set(benchmarkCtx, key, value, time.Hour)
			}

			_, err := txPipe.Exec(benchmarkCtx)
			if err != nil {
				b.Fatal(err)
			}
			txPipe.Close()
		}
	})
}

// BenchmarkConcurrentOperations 并发操作基准测试
func BenchmarkConcurrentOperations(b *testing.B) {
	b.Run("Concurrent_Set", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				key := fmt.Sprintf("benchmark:concurrent:set:%d", i)
				value := fmt.Sprintf("value_%d", i)
				err := benchmarkClient.Set(benchmarkCtx, key, value, time.Hour)
				if err != nil {
					b.Fatal(err)
				}
				i++
			}
		})
	})

	b.Run("Concurrent_Get", func(b *testing.B) {
		// 预设数据
		for i := 0; i < 10000; i++ {
			key := fmt.Sprintf("benchmark:concurrent:get:%d", i)
			value := fmt.Sprintf("value_%d", i)
			benchmarkClient.Set(benchmarkCtx, key, value, time.Hour)
		}

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				key := fmt.Sprintf("benchmark:concurrent:get:%d", i%10000)
				_, err := benchmarkClient.Get(benchmarkCtx, key)
				if err != nil {
					b.Fatal(err)
				}
				i++
			}
		})
	})

	b.Run("Concurrent_Incr", func(b *testing.B) {
		key := "benchmark:concurrent:incr"
		benchmarkClient.Set(benchmarkCtx, key, "0", time.Hour)

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := benchmarkClient.Incr(benchmarkCtx, key)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

// BenchmarkMemoryUsage 内存使用基准测试
func BenchmarkMemoryUsage(b *testing.B) {
	b.Run("Large_String_Values", func(b *testing.B) {
		// 创建1KB的值
		largeValue := string(make([]byte, 1024))
		for i := range largeValue {
			largeValue = largeValue[:i] + "a" + largeValue[i+1:]
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("benchmark:memory:large:%d", i)
			err := benchmarkClient.Set(benchmarkCtx, key, largeValue, time.Hour)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Many_Small_Keys", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("benchmark:memory:small:%d", i)
			value := strconv.Itoa(i)
			err := benchmarkClient.Set(benchmarkCtx, key, value, time.Hour)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Hash_Many_Fields", func(b *testing.B) {
		hashKey := "benchmark:memory:hash"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			field := fmt.Sprintf("field_%d", i)
			value := fmt.Sprintf("value_%d", i)
			err := benchmarkClient.HSet(benchmarkCtx, hashKey, field, value)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkConnectionPool 连接池基准测试
func BenchmarkConnectionPool(b *testing.B) {
	b.Run("Pool_Size_10", func(b *testing.B) {
		config := &cache.Config{
			Mode: cache.ModeSingle,
			Single: &cache.SingleConfig{
				Addr: "localhost:6379",
				DB:   0,
			},
			Common: cache.CommonConfig{
				Password:    "",
				PoolSize:    10,
				DialTimeout: 5 * time.Second,
			},
		}

		factory, err := cache.NewFactory(config)
		if err != nil {
			b.Fatal(err)
		}

		client, err := factory.CreateClient()
		if err != nil {
			b.Fatal(err)
		}
		defer client.Close()

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				key := fmt.Sprintf("benchmark:pool:10:%d", i)
				value := fmt.Sprintf("value_%d", i)
				err := client.Set(benchmarkCtx, key, value, time.Hour)
				if err != nil {
					b.Fatal(err)
				}
				i++
			}
		})
	})

	b.Run("Pool_Size_100", func(b *testing.B) {
		config := &cache.Config{
			Mode: cache.ModeSingle,
			Single: &cache.SingleConfig{
				Addr: "localhost:6379",
				DB:   0,
			},
			Common: cache.CommonConfig{
				Password:    "",
				PoolSize:    100,
				DialTimeout: 5 * time.Second,
			},
		}

		factory, err := cache.NewFactory(config)
		if err != nil {
			b.Fatal(err)
		}

		client, err := factory.CreateClient()
		if err != nil {
			b.Fatal(err)
		}
		defer client.Close()

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				key := fmt.Sprintf("benchmark:pool:100:%d", i)
				value := fmt.Sprintf("value_%d", i)
				err := client.Set(benchmarkCtx, key, value, time.Hour)
				if err != nil {
					b.Fatal(err)
				}
				i++
			}
		})
	})
}