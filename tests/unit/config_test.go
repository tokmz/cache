package unit

import (
	"testing"
	"time"

	"cache"
)

// TestConfigValidation 测试配置验证功能
func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  *cache.Config
		wantErr bool
	}{
		{
			name: "有效的单机模式配置",
			config: &cache.Config{
				Mode: cache.ModeSingle,
				Single: &cache.SingleConfig{
					Addr: "localhost:6379",
					DB:   0,
				},
				Common: cache.CommonConfig{
					Password: "",
				},
			},
			wantErr: false,
		},
		{
			name: "有效的集群模式配置",
			config: &cache.Config{
				Mode: cache.ModeCluster,
				Cluster: &cache.ClusterConfig{
					Addrs: []string{"localhost:7000", "localhost:7001", "localhost:7002"},
				},
				Common: cache.CommonConfig{
					Password: "",
				},
			},
			wantErr: false,
		},
		{
			name: "有效的哨兵模式配置",
			config: &cache.Config{
				Mode: cache.ModeSentinel,
				Sentinel: &cache.SentinelConfig{
					MasterName: "mymaster",
					Addrs: []string{"localhost:26379", "localhost:26380", "localhost:26381"},
					DB: 0,
				},
				Common: cache.CommonConfig{
					Password: "",
				},
			},
			wantErr: false,
		},
		{
			name: "无效模式",
			config: &cache.Config{
				Mode: "invalid",
			},
			wantErr: true,
		},
		{
			name: "单机模式缺少配置",
			config: &cache.Config{
				Mode: cache.ModeSingle,
			},
			wantErr: true,
		},
		{
			name: "集群模式缺少地址",
			config: &cache.Config{
				Mode: cache.ModeCluster,
				Cluster: &cache.ClusterConfig{
					Addrs: []string{},
				},
			},
			wantErr: true,
		},
		{
			name: "哨兵模式缺少主节点名称",
			config: &cache.Config{
				Mode: cache.ModeSentinel,
				Sentinel: &cache.SentinelConfig{
					Addrs: []string{"localhost:26379"},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestFactoryCreation 测试工厂模式创建
func TestFactoryCreation(t *testing.T) {
	tests := []struct {
		name    string
		config  *cache.Config
		wantErr bool
	}{
		{
			name: "创建单机模式工厂",
			config: &cache.Config{
				Mode: cache.ModeSingle,
				Single: &cache.SingleConfig{
					Addr: "localhost:6379",
					DB:   0,
				},
				Common: cache.CommonConfig{
					Password: "",
				},
			},
			wantErr: false,
		},
		{
			name: "创建集群模式工厂",
			config: &cache.Config{
				Mode: cache.ModeCluster,
				Cluster: &cache.ClusterConfig{
					Addrs: []string{"localhost:7000", "localhost:7001"},
				},
				Common: cache.CommonConfig{
					Password: "",
				},
			},
			wantErr: false,
		},
		{
			name: "创建哨兵模式工厂",
			config: &cache.Config{
				Mode: cache.ModeSentinel,
				Sentinel: &cache.SentinelConfig{
					MasterName: "mymaster",
					Addrs: []string{"localhost:26379"},
					DB: 0,
				},
				Common: cache.CommonConfig{
					Password: "",
				},
			},
			wantErr: false,
		},
		{
			name: "无效配置创建工厂",
			config: &cache.Config{
				Mode: "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory, err := cache.NewFactory(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFactory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && factory == nil {
				t.Error("NewFactory() returned nil factory")
			}
		})
	}
}

// TestConfigDefaults 测试配置默认值
func TestConfigDefaults(t *testing.T) {
	config := &cache.Config{
		Mode: cache.ModeSingle,
		Single: &cache.SingleConfig{
			Addr: "localhost:6379",
		},
		Common: cache.CommonConfig{},
	}

	// 验证默认值设置
	if config.Common.PoolSize == 0 {
		config.Common.PoolSize = 10
	}
	if config.Common.MinIdleConns == 0 {
		config.Common.MinIdleConns = 5
	}
	if config.Common.DialTimeout == 0 {
		config.Common.DialTimeout = 5 * time.Second
	}
	if config.Common.ReadTimeout == 0 {
		config.Common.ReadTimeout = 3 * time.Second
	}
	if config.Common.WriteTimeout == 0 {
		config.Common.WriteTimeout = 3 * time.Second
	}

	// 验证默认值是否正确设置
	if config.Common.PoolSize != 10 {
		t.Errorf("Expected PoolSize to be 10, got %d", config.Common.PoolSize)
	}
	if config.Common.MinIdleConns != 5 {
		t.Errorf("Expected MinIdleConns to be 5, got %d", config.Common.MinIdleConns)
	}
	if config.Common.DialTimeout != 5*time.Second {
		t.Errorf("Expected DialTimeout to be 5s, got %v", config.Common.DialTimeout)
	}
}

// BenchmarkConfigValidation 配置验证性能基准测试
func BenchmarkConfigValidation(b *testing.B) {
	config := &cache.Config{
		Mode: cache.ModeSingle,
		Single: &cache.SingleConfig{
			Addr: "localhost:6379",
			DB:   0,
		},
		Common: cache.CommonConfig{
			Password: "",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = config.Validate()
	}
}

// BenchmarkFactoryCreation 工厂创建性能基准测试
func BenchmarkFactoryCreation(b *testing.B) {
	config := &cache.Config{
		Mode: cache.ModeSingle,
		Single: &cache.SingleConfig{
			Addr: "localhost:6379",
			DB:   0,
		},
		Common: cache.CommonConfig{
			Password: "",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		factory, _ := cache.NewFactory(config)
		if factory != nil {
			// 避免编译器优化
			_ = factory
		}
	}
}