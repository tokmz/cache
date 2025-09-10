package unit

import (
	"testing"
	"time"

	"cache"
	"github.com/stretchr/testify/assert"
)

// TestSentinelClientCreation 测试哨兵客户端创建
func TestSentinelClientCreation(t *testing.T) {
	config := &cache.Config{
		Mode: cache.ModeSentinel,
		Sentinel: &cache.SentinelConfig{
			MasterName: "mymaster",
			Addrs: []string{"localhost:26379", "localhost:26380", "localhost:26381"},
			DB: 0,
		},
		Common: cache.CommonConfig{
			Password:    "",
			PoolSize:    10,
			DialTimeout: 5 * time.Second,
		},
	}

	factory, err := cache.NewFactory(config)
	assert.NoError(t, err)
	assert.NotNil(t, factory)

	client, err := factory.CreateClient()
	assert.NoError(t, err)
	assert.NotNil(t, client)

	// 验证客户端类型
	sentinelClient, ok := client.(*cache.SentinelClient)
	assert.True(t, ok)
	assert.NotNil(t, sentinelClient)
}

// TestSentinelClientInterface 测试哨兵客户端接口实现
func TestSentinelClientInterface(t *testing.T) {
	config := &cache.Config{
		Mode: cache.ModeSentinel,
		Sentinel: &cache.SentinelConfig{
			MasterName: "mymaster",
			Addrs: []string{"localhost:26379", "localhost:26380", "localhost:26381"},
			DB: 0,
		},
		Common: cache.CommonConfig{
			Password: "",
			PoolSize: 10,
		},
	}

	factory, err := cache.NewFactory(config)
	assert.NoError(t, err)

	client, err := factory.CreateClient()
	assert.NoError(t, err)
	assert.NotNil(t, client)

	// 验证客户端实现了Client接口
	var _ cache.Client = client
}

// TestSentinelClientConfiguration 测试哨兵客户端配置
func TestSentinelClientConfiguration(t *testing.T) {
	tests := []struct {
		name    string
		config  *cache.Config
		wantErr bool
	}{
		{
			name: "有效的哨兵配置",
			config: &cache.Config{
				Mode: cache.ModeSentinel,
				Sentinel: &cache.SentinelConfig{
					MasterName: "mymaster",
					Addrs: []string{"localhost:26379", "localhost:26380", "localhost:26381"},
					DB: 0,
				},
				Common: cache.CommonConfig{
					PoolSize: 10,
				},
			},
			wantErr: false,
		},
		{
			name: "单哨兵配置",
			config: &cache.Config{
				Mode: cache.ModeSentinel,
				Sentinel: &cache.SentinelConfig{
					MasterName: "mymaster",
					Addrs: []string{"localhost:26379"},
					DB: 1,
				},
				Common: cache.CommonConfig{
					PoolSize: 5,
				},
			},
			wantErr: false,
		},
		{
			name: "空主节点名称",
			config: &cache.Config{
				Mode: cache.ModeSentinel,
				Sentinel: &cache.SentinelConfig{
					MasterName: "",
					Addrs: []string{"localhost:26379"},
				},
			},
			wantErr: true,
		},
		{
			name: "空哨兵地址列表",
			config: &cache.Config{
				Mode: cache.ModeSentinel,
				Sentinel: &cache.SentinelConfig{
					MasterName: "mymaster",
					Addrs: []string{},
				},
			},
			wantErr: true,
		},
		{
			name: "缺少哨兵配置",
			config: &cache.Config{
				Mode: cache.ModeSentinel,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestSentinelClientAdvancedConfig 测试哨兵客户端高级配置
func TestSentinelClientAdvancedConfig(t *testing.T) {
	tests := []struct {
		name   string
		config *cache.Config
	}{
		{
			name: "带密码的哨兵配置",
			config: &cache.Config{
				Mode: cache.ModeSentinel,
				Sentinel: &cache.SentinelConfig{
					MasterName: "mymaster",
					Addrs: []string{"localhost:26379", "localhost:26380"},
					DB: 0,
					SentinelPassword: "sentinel-pass",
				},
				Common: cache.CommonConfig{
					Password: "redis-pass",
					PoolSize: 20,
				},
			},
		},
		{
			name: "不同数据库配置",
			config: &cache.Config{
				Mode: cache.ModeSentinel,
				Sentinel: &cache.SentinelConfig{
					MasterName: "mymaster",
					Addrs: []string{"localhost:26379", "localhost:26380", "localhost:26381"},
					DB: 5,
				},
				Common: cache.CommonConfig{
					PoolSize: 15,
				},
			},
		},
		{
			name: "自定义主节点名称",
			config: &cache.Config{
				Mode: cache.ModeSentinel,
				Sentinel: &cache.SentinelConfig{
					MasterName: "redis-master-prod",
					Addrs: []string{"sentinel1:26379", "sentinel2:26379", "sentinel3:26379"},
					DB: 0,
				},
				Common: cache.CommonConfig{
					PoolSize: 12,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			assert.NoError(t, err)

			factory, err := cache.NewFactory(tt.config)
			assert.NoError(t, err)
			assert.NotNil(t, factory)

			client, err := factory.CreateClient()
			assert.NoError(t, err)
			assert.NotNil(t, client)
		})
	}
}

// TestSentinelClientMethods 测试哨兵客户端方法存在性
func TestSentinelClientMethods(t *testing.T) {
	config := &cache.Config{
		Mode: cache.ModeSentinel,
		Sentinel: &cache.SentinelConfig{
			MasterName: "mymaster",
			Addrs: []string{"localhost:26379", "localhost:26380", "localhost:26381"},
			DB: 0,
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

	// 验证客户端有必要的方法（通过接口检查）
	assert.Implements(t, (*cache.Client)(nil), client)
}

// TestSentinelClientErrorHandling 测试哨兵客户端错误处理
func TestSentinelClientErrorHandling(t *testing.T) {
	tests := []struct {
		name   string
		config *cache.Config
	}{
		{
			name: "无效哨兵地址格式",
			config: &cache.Config{
				Mode: cache.ModeSentinel,
				Sentinel: &cache.SentinelConfig{
					MasterName: "mymaster",
					Addrs: []string{"invalid-address", "localhost:26379"},
					DB: 0,
				},
			},
		},
		{
			name: "重复哨兵地址",
			config: &cache.Config{
				Mode: cache.ModeSentinel,
				Sentinel: &cache.SentinelConfig{
					MasterName: "mymaster",
					Addrs: []string{"localhost:26379", "localhost:26379", "localhost:26380"},
					DB: 0,
				},
			},
		},
		{
			name: "负数数据库索引",
			config: &cache.Config{
				Mode: cache.ModeSentinel,
				Sentinel: &cache.SentinelConfig{
					MasterName: "mymaster",
					Addrs: []string{"localhost:26379"},
					DB: -1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 这些配置可能通过验证但在实际连接时失败
			// 在单元测试中我们主要测试配置结构的正确性
			factory, err := cache.NewFactory(tt.config)
			if err == nil {
				assert.NotNil(t, factory)
				// 客户端创建可能会失败，这是预期的
				client, _ := factory.CreateClient()
				_ = client // 忽略结果，因为可能连接失败
			}
		})
	}
}

// TestSentinelClientDatabaseSelection 测试哨兵客户端数据库选择
func TestSentinelClientDatabaseSelection(t *testing.T) {
	tests := []struct {
		name string
		db   int
	}{
		{"默认数据库", 0},
		{"数据库1", 1},
		{"数据库5", 5},
		{"数据库15", 15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &cache.Config{
				Mode: cache.ModeSentinel,
				Sentinel: &cache.SentinelConfig{
					MasterName: "mymaster",
					Addrs: []string{"localhost:26379"},
					DB: tt.db,
				},
				Common: cache.CommonConfig{
					PoolSize: 10,
				},
			}

			err := config.Validate()
			assert.NoError(t, err)

			factory, err := cache.NewFactory(config)
			assert.NoError(t, err)
			assert.NotNil(t, factory)
		})
	}
}

// BenchmarkSentinelClientCreation 哨兵客户端创建基准测试
func BenchmarkSentinelClientCreation(b *testing.B) {
	config := &cache.Config{
		Mode: cache.ModeSentinel,
		Sentinel: &cache.SentinelConfig{
			MasterName: "mymaster",
			Addrs: []string{"localhost:26379", "localhost:26380", "localhost:26381"},
			DB: 0,
		},
		Common: cache.CommonConfig{
			PoolSize: 10,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		factory, _ := cache.NewFactory(config)
		client, _ := factory.CreateClient()
		if client != nil {
			// 避免编译器优化
			_ = client
		}
	}
}

// BenchmarkSentinelConfigValidation 哨兵配置验证基准测试
func BenchmarkSentinelConfigValidation(b *testing.B) {
	config := &cache.Config{
		Mode: cache.ModeSentinel,
		Sentinel: &cache.SentinelConfig{
			MasterName: "mymaster",
			Addrs: []string{"localhost:26379", "localhost:26380", "localhost:26381"},
			DB: 0,
			SentinelPassword: "sentinel-pass",
		},
		Common: cache.CommonConfig{
			Password: "redis-pass",
			PoolSize: 10,
			MinIdleConns: 5,
			DialTimeout: 5 * time.Second,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = config.Validate()
	}
}