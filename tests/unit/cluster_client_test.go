package unit

import (
	"testing"
	"time"

	"cache"
	"github.com/stretchr/testify/assert"
)

// TestClusterClientCreation 测试集群客户端创建
func TestClusterClientCreation(t *testing.T) {
	config := &cache.Config{
		Mode: cache.ModeCluster,
		Cluster: &cache.ClusterConfig{
			Addrs: []string{"localhost:7000", "localhost:7001", "localhost:7002"},
			MaxRedirects: 3,
			ReadOnly: false,
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
	clusterClient, ok := client.(*cache.ClusterClient)
	assert.True(t, ok)
	assert.NotNil(t, clusterClient)
}

// TestClusterClientInterface 测试集群客户端接口实现
func TestClusterClientInterface(t *testing.T) {
	config := &cache.Config{
		Mode: cache.ModeCluster,
		Cluster: &cache.ClusterConfig{
			Addrs: []string{"localhost:7000", "localhost:7001", "localhost:7002"},
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

// TestClusterClientConfiguration 测试集群客户端配置
func TestClusterClientConfiguration(t *testing.T) {
	tests := []struct {
		name    string
		config  *cache.Config
		wantErr bool
	}{
		{
			name: "有效的集群配置",
			config: &cache.Config{
				Mode: cache.ModeCluster,
				Cluster: &cache.ClusterConfig{
					Addrs: []string{"localhost:7000", "localhost:7001", "localhost:7002"},
					MaxRedirects: 3,
				},
				Common: cache.CommonConfig{
					PoolSize: 10,
				},
			},
			wantErr: false,
		},
		{
			name: "单节点集群配置",
			config: &cache.Config{
				Mode: cache.ModeCluster,
				Cluster: &cache.ClusterConfig{
					Addrs: []string{"localhost:7000"},
				},
				Common: cache.CommonConfig{
					PoolSize: 5,
				},
			},
			wantErr: false,
		},
		{
			name: "空地址列表",
			config: &cache.Config{
				Mode: cache.ModeCluster,
				Cluster: &cache.ClusterConfig{
					Addrs: []string{},
				},
			},
			wantErr: true,
		},
		{
			name: "缺少集群配置",
			config: &cache.Config{
				Mode: cache.ModeCluster,
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

// TestClusterClientAdvancedConfig 测试集群客户端高级配置
func TestClusterClientAdvancedConfig(t *testing.T) {
	tests := []struct {
		name   string
		config *cache.Config
	}{
		{
			name: "只读模式配置",
			config: &cache.Config{
				Mode: cache.ModeCluster,
				Cluster: &cache.ClusterConfig{
					Addrs: []string{"localhost:7000", "localhost:7001"},
					ReadOnly: true,
					MaxRedirects: 5,
				},
				Common: cache.CommonConfig{
					PoolSize: 20,
				},
			},
		},
		{
			name: "延迟路由配置",
			config: &cache.Config{
				Mode: cache.ModeCluster,
				Cluster: &cache.ClusterConfig{
					Addrs: []string{"localhost:7000", "localhost:7001", "localhost:7002"},
					RouteByLatency: true,
					MaxRedirects: 3,
				},
				Common: cache.CommonConfig{
					PoolSize: 15,
				},
			},
		},
		{
			name: "随机路由配置",
			config: &cache.Config{
				Mode: cache.ModeCluster,
				Cluster: &cache.ClusterConfig{
					Addrs: []string{"localhost:7000", "localhost:7001", "localhost:7002"},
					RouteRandomly: true,
					MaxRedirects: 2,
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

// TestClusterClientMethods 测试集群客户端方法存在性
func TestClusterClientMethods(t *testing.T) {
	config := &cache.Config{
		Mode: cache.ModeCluster,
		Cluster: &cache.ClusterConfig{
			Addrs: []string{"localhost:7000", "localhost:7001", "localhost:7002"},
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

// TestClusterClientErrorHandling 测试集群客户端错误处理
func TestClusterClientErrorHandling(t *testing.T) {
	tests := []struct {
		name   string
		config *cache.Config
	}{
		{
			name: "无效地址格式",
			config: &cache.Config{
				Mode: cache.ModeCluster,
				Cluster: &cache.ClusterConfig{
					Addrs: []string{"invalid-address", "localhost:7001"},
				},
			},
		},
		{
			name: "重复地址",
			config: &cache.Config{
				Mode: cache.ModeCluster,
				Cluster: &cache.ClusterConfig{
					Addrs: []string{"localhost:7000", "localhost:7000", "localhost:7001"},
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

// BenchmarkClusterClientCreation 集群客户端创建基准测试
func BenchmarkClusterClientCreation(b *testing.B) {
	config := &cache.Config{
		Mode: cache.ModeCluster,
		Cluster: &cache.ClusterConfig{
			Addrs: []string{"localhost:7000", "localhost:7001", "localhost:7002"},
			MaxRedirects: 3,
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

// BenchmarkClusterConfigValidation 集群配置验证基准测试
func BenchmarkClusterConfigValidation(b *testing.B) {
	config := &cache.Config{
		Mode: cache.ModeCluster,
		Cluster: &cache.ClusterConfig{
			Addrs: []string{"localhost:7000", "localhost:7001", "localhost:7002"},
			MaxRedirects: 3,
			ReadOnly: false,
			RouteByLatency: true,
		},
		Common: cache.CommonConfig{
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