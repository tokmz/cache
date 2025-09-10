package cache

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Factory Redis客户端工厂
type Factory struct {
	config *Config
}

// NewFactory 创建新的工厂实例
func NewFactory(config *Config) (*Factory, error) {
	if config == nil {
		config = DefaultConfig()
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &Factory{
		config: config,
	}, nil
}

// CreateClient 根据配置创建Redis客户端
func (f *Factory) CreateClient() (Client, error) {
	switch f.config.Mode {
	case ModeSingle:
		return f.createSingleClient()
	case ModeCluster:
		return f.createClusterClient()
	case ModeSentinel:
		return f.createSentinelClient()
	default:
		return nil, fmt.Errorf("unsupported mode: %s", f.config.Mode)
	}
}

// createSingleClient 创建单机模式客户端
func (f *Factory) createSingleClient() (Client, error) {
	opts := &redis.Options{
		Addr:     f.config.Single.Addr,
		DB:       f.config.Single.DB,
		Username: f.config.Common.Username,
		Password: f.config.Common.Password,

		// 连接池配置
		PoolSize:     f.config.Common.PoolSize,
		MinIdleConns: f.config.Common.MinIdleConns,
		MaxIdleConns: f.config.Common.MaxIdleConns,
		PoolTimeout:  f.config.Common.PoolTimeout,

		// 网络配置
		DialTimeout:  f.config.Common.DialTimeout,
		ReadTimeout:  f.config.Common.ReadTimeout,
		WriteTimeout: f.config.Common.WriteTimeout,

		// 重试配置
		MaxRetries:      f.config.Common.MaxRetries,
		MinRetryBackoff: f.config.Common.MinRetryBackoff,
		MaxRetryBackoff: f.config.Common.MaxRetryBackoff,
	}

	// 配置TLS
	if f.config.Common.TLSConfig != nil && f.config.Common.TLSConfig.Enabled {
		tlsConfig, err := f.buildTLSConfig(f.config.Common.TLSConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to build TLS config: %w", err)
		}
		opts.TLSConfig = tlsConfig
	}

	rdb := redis.NewClient(opts)

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		rdb.Close()
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &SingleClient{
		client: rdb,
		config: f.config,
	}, nil
}

// createClusterClient 创建集群模式客户端
func (f *Factory) createClusterClient() (Client, error) {
	opts := &redis.ClusterOptions{
		Addrs:    f.config.Cluster.Addrs,
		Username: f.config.Common.Username,
		Password: f.config.Common.Password,

		// 集群特定配置
		MaxRedirects:   f.config.Cluster.MaxRedirects,
		ReadOnly:       f.config.Cluster.ReadOnly,
		RouteByLatency: f.config.Cluster.RouteByLatency,
		RouteRandomly:  f.config.Cluster.RouteRandomly,

		// 连接池配置
		PoolSize:     f.config.Common.PoolSize,
		MinIdleConns: f.config.Common.MinIdleConns,
		MaxIdleConns: f.config.Common.MaxIdleConns,
		PoolTimeout:  f.config.Common.PoolTimeout,

		// 网络配置
		DialTimeout:  f.config.Common.DialTimeout,
		ReadTimeout:  f.config.Common.ReadTimeout,
		WriteTimeout: f.config.Common.WriteTimeout,

		// 重试配置
		MaxRetries:      f.config.Common.MaxRetries,
		MinRetryBackoff: f.config.Common.MinRetryBackoff,
		MaxRetryBackoff: f.config.Common.MaxRetryBackoff,
	}

	// 配置TLS
	if f.config.Common.TLSConfig != nil && f.config.Common.TLSConfig.Enabled {
		tlsConfig, err := f.buildTLSConfig(f.config.Common.TLSConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to build TLS config: %w", err)
		}
		opts.TLSConfig = tlsConfig
	}

	// 设置默认值
	if opts.MaxRedirects == 0 {
		opts.MaxRedirects = 3
	}

	rdb := redis.NewClusterClient(opts)

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		rdb.Close()
		return nil, fmt.Errorf("failed to connect to redis cluster: %w", err)
	}

	return &ClusterClient{
		client: rdb,
		config: f.config,
	}, nil
}

// createSentinelClient 创建哨兵模式客户端
func (f *Factory) createSentinelClient() (Client, error) {
	opts := &redis.FailoverOptions{
		MasterName:       f.config.Sentinel.MasterName,
		SentinelAddrs:    f.config.Sentinel.Addrs,
		SentinelPassword: f.config.Sentinel.SentinelPassword,
		DB:               f.config.Sentinel.DB,
		Username:         f.config.Common.Username,
		Password:         f.config.Common.Password,

		// 连接池配置
		PoolSize:     f.config.Common.PoolSize,
		MinIdleConns: f.config.Common.MinIdleConns,
		MaxIdleConns: f.config.Common.MaxIdleConns,
		PoolTimeout:  f.config.Common.PoolTimeout,

		// 网络配置
		DialTimeout:  f.config.Common.DialTimeout,
		ReadTimeout:  f.config.Common.ReadTimeout,
		WriteTimeout: f.config.Common.WriteTimeout,

		// 重试配置
		MaxRetries:      f.config.Common.MaxRetries,
		MinRetryBackoff: f.config.Common.MinRetryBackoff,
		MaxRetryBackoff: f.config.Common.MaxRetryBackoff,
	}

	// 配置TLS
	if f.config.Common.TLSConfig != nil && f.config.Common.TLSConfig.Enabled {
		tlsConfig, err := f.buildTLSConfig(f.config.Common.TLSConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to build TLS config: %w", err)
		}
		opts.TLSConfig = tlsConfig
	}

	rdb := redis.NewFailoverClient(opts)

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		rdb.Close()
		return nil, fmt.Errorf("failed to connect to redis sentinel: %w", err)
	}

	return &SentinelClient{
		client: rdb,
		config: f.config,
	}, nil
}

// buildTLSConfig 构建TLS配置
func (f *Factory) buildTLSConfig(tlsConfig *TLSConfig) (*tls.Config, error) {
	config := &tls.Config{
		InsecureSkipVerify: tlsConfig.InsecureSkipVerify,
		ServerName:         tlsConfig.ServerName,
	}

	// 加载证书
	if tlsConfig.CertFile != "" && tlsConfig.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(tlsConfig.CertFile, tlsConfig.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load certificate: %w", err)
		}
		config.Certificates = []tls.Certificate{cert}
	}

	// 加载CA证书
	if tlsConfig.CAFile != "" {
		// 这里可以添加CA证书加载逻辑
		// 为了简化，暂时跳过
	}

	return config, nil
}

// GetConfig 获取配置
func (f *Factory) GetConfig() *Config {
	return f.config
}

// UpdateConfig 更新配置
func (f *Factory) UpdateConfig(config *Config) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	f.config = config
	return nil
}

// 便捷函数

// NewSingleClient 创建单机模式客户端的便捷函数
func NewSingleClient(addr string, db int, password string) (Client, error) {
	config := DefaultConfig()
	config.Mode = ModeSingle
	config.Single.Addr = addr
	config.Single.DB = db
	config.Common.Password = password

	factory, err := NewFactory(config)
	if err != nil {
		return nil, err
	}

	return factory.CreateClient()
}

// NewClusterClient 创建集群模式客户端的便捷函数
func NewClusterClient(addrs []string, password string) (Client, error) {
	config := DefaultConfig()
	config.Mode = ModeCluster
	config.Cluster = &ClusterConfig{
		Addrs:        addrs,
		MaxRedirects: 3,
	}
	config.Common.Password = password

	factory, err := NewFactory(config)
	if err != nil {
		return nil, err
	}

	return factory.CreateClient()
}

// NewSentinelClient 创建哨兵模式客户端的便捷函数
func NewSentinelClient(sentinelAddrs []string, masterName string, db int, password string) (Client, error) {
	config := DefaultConfig()
	config.Mode = ModeSentinel
	config.Sentinel = &SentinelConfig{
		Addrs:      sentinelAddrs,
		MasterName: masterName,
		DB:         db,
	}
	config.Common.Password = password

	factory, err := NewFactory(config)
	if err != nil {
		return nil, err
	}

	return factory.CreateClient()
}

// NewClientFromConfig 从配置创建客户端的便捷函数
func NewClientFromConfig(config *Config) (Client, error) {
	factory, err := NewFactory(config)
	if err != nil {
		return nil, err
	}

	return factory.CreateClient()
}