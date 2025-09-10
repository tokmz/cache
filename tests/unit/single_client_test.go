package unit

import (
	"testing"
	"time"

	"cache"
	"github.com/stretchr/testify/assert"
)

// TestSingleClientCreation 测试单机客户端创建
func TestSingleClientCreation(t *testing.T) {
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
	assert.NoError(t, err)
	assert.NotNil(t, factory)

	client, err := factory.CreateClient()
	assert.NoError(t, err)
	assert.NotNil(t, client)

	// 验证客户端类型
	singleClient, ok := client.(*cache.SingleClient)
	assert.True(t, ok)
	assert.NotNil(t, singleClient)
}

// TestSingleClientInterface 测试单机客户端接口实现
func TestSingleClientInterface(t *testing.T) {
	config := &cache.Config{
		Mode: cache.ModeSingle,
		Single: &cache.SingleConfig{
			Addr: "localhost:6379",
			DB:   0,
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

// TestSingleClientConfiguration 测试单机客户端配置
func TestSingleClientConfiguration(t *testing.T) {
	tests := []struct {
		name    string
		config  *cache.Config
		wantErr bool
	}{
		{
			name: "有效配置",
			config: &cache.Config{
				Mode: cache.ModeSingle,
				Single: &cache.SingleConfig{
					Addr: "localhost:6379",
					DB:   0,
				},
				Common: cache.CommonConfig{
					PoolSize: 10,
				},
			},
			wantErr: false,
		},
		{
			name: "缺少地址配置",
			config: &cache.Config{
				Mode: cache.ModeSingle,
				Single: &cache.SingleConfig{
					Addr: "",
					DB:   0,
				},
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

// TestSingleClientMethods 测试单机客户端方法存在性
func TestSingleClientMethods(t *testing.T) {
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

	// 验证客户端有必要的方法（通过接口检查）
	assert.Implements(t, (*cache.Client)(nil), client)
}

// TestSingleClientErrorHandling 测试错误处理
func TestSingleClientErrorHandling(t *testing.T) {
	tests := []struct {
		name   string
		config *cache.Config
	}{
		{
			name: "无效地址配置",
			config: &cache.Config{
				Mode: cache.ModeSingle,
				Single: &cache.SingleConfig{
					Addr: "", // 空地址
					DB:   0,
				},
			},
		},
		{
			name: "无效数据库索引",
			config: &cache.Config{
				Mode: cache.ModeSingle,
				Single: &cache.SingleConfig{
					Addr: "localhost:6379",
					DB:   -1, // 无效DB索引
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			assert.Error(t, err)
		})
	}
}

// TestSingleClientConcurrency 测试并发安全性
func TestSingleClientConcurrency(t *testing.T) {
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

	// 并发测试将在集成测试中进行
	// 这里只验证客户端创建成功
}

// BenchmarkSingleClientCreation 单机客户端创建基准测试
func BenchmarkSingleClientCreation(b *testing.B) {
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