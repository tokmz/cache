package cache

import (
	"time"
)

// Config Redis客户端配置
type Config struct {
	// 部署模式
	Mode Mode `json:"mode" yaml:"mode"`

	// 单机模式配置
	Single *SingleConfig `json:"single,omitempty" yaml:"single,omitempty"`

	// 集群模式配置
	Cluster *ClusterConfig `json:"cluster,omitempty" yaml:"cluster,omitempty"`

	// 哨兵模式配置
	Sentinel *SentinelConfig `json:"sentinel,omitempty" yaml:"sentinel,omitempty"`

	// 通用配置
	Common CommonConfig `json:"common" yaml:"common"`
}

// Mode Redis部署模式
type Mode string

const (
	// ModeSingle 单机模式
	ModeSingle Mode = "single"
	// ModeCluster 集群模式
	ModeCluster Mode = "cluster"
	// ModeSentinel 哨兵模式
	ModeSentinel Mode = "sentinel"
)

// SingleConfig 单机模式配置
type SingleConfig struct {
	// Redis服务器地址
	Addr string `json:"addr" yaml:"addr"`
	// 数据库索引
	DB int `json:"db" yaml:"db"`
}

// ClusterConfig 集群模式配置
type ClusterConfig struct {
	// 集群节点地址列表
	Addrs []string `json:"addrs" yaml:"addrs"`
	// 最大重定向次数
	MaxRedirects int `json:"max_redirects" yaml:"max_redirects"`
	// 只读模式
	ReadOnly bool `json:"read_only" yaml:"read_only"`
	// 路由模式
	RouteByLatency bool `json:"route_by_latency" yaml:"route_by_latency"`
	RouteRandomly  bool `json:"route_randomly" yaml:"route_randomly"`
}

// SentinelConfig 哨兵模式配置
type SentinelConfig struct {
	// 哨兵节点地址列表
	Addrs []string `json:"addrs" yaml:"addrs"`
	// 主节点名称
	MasterName string `json:"master_name" yaml:"master_name"`
	// 数据库索引
	DB int `json:"db" yaml:"db"`
	// 哨兵密码
	SentinelPassword string `json:"sentinel_password,omitempty" yaml:"sentinel_password,omitempty"`
}

// CommonConfig 通用配置
type CommonConfig struct {
	// 认证密码
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
	// 用户名（Redis 6.0+）
	Username string `json:"username,omitempty" yaml:"username,omitempty"`

	// 连接池配置
	PoolSize     int           `json:"pool_size" yaml:"pool_size"`         // 连接池大小
	MinIdleConns int           `json:"min_idle_conns" yaml:"min_idle_conns"` // 最小空闲连接数
	MaxIdleConns int           `json:"max_idle_conns" yaml:"max_idle_conns"` // 最大空闲连接数
	ConnMaxAge   time.Duration `json:"conn_max_age" yaml:"conn_max_age"`     // 连接最大存活时间
	PoolTimeout  time.Duration `json:"pool_timeout" yaml:"pool_timeout"`     // 获取连接超时时间
	IdleTimeout  time.Duration `json:"idle_timeout" yaml:"idle_timeout"`     // 空闲连接超时时间

	// 网络配置
	DialTimeout  time.Duration `json:"dial_timeout" yaml:"dial_timeout"`   // 连接超时时间
	ReadTimeout  time.Duration `json:"read_timeout" yaml:"read_timeout"`   // 读取超时时间
	WriteTimeout time.Duration `json:"write_timeout" yaml:"write_timeout"` // 写入超时时间

	// 重试配置
	MaxRetries      int           `json:"max_retries" yaml:"max_retries"`           // 最大重试次数
	MinRetryBackoff time.Duration `json:"min_retry_backoff" yaml:"min_retry_backoff"` // 最小重试间隔
	MaxRetryBackoff time.Duration `json:"max_retry_backoff" yaml:"max_retry_backoff"` // 最大重试间隔

	// 键前缀配置
	KeyPrefix string `json:"key_prefix,omitempty" yaml:"key_prefix,omitempty"`

	// 默认TTL配置
	DefaultTTL time.Duration `json:"default_ttl" yaml:"default_ttl"`

	// TLS配置
	TLSConfig *TLSConfig `json:"tls,omitempty" yaml:"tls,omitempty"`

	// 监控配置
	EnableMetrics bool `json:"enable_metrics" yaml:"enable_metrics"`
	EnableTracing bool `json:"enable_tracing" yaml:"enable_tracing"`
}

// TLSConfig TLS配置
type TLSConfig struct {
	// 启用TLS
	Enabled bool `json:"enabled" yaml:"enabled"`
	// 跳过证书验证
	InsecureSkipVerify bool `json:"insecure_skip_verify" yaml:"insecure_skip_verify"`
	// 证书文件路径
	CertFile string `json:"cert_file,omitempty" yaml:"cert_file,omitempty"`
	// 私钥文件路径
	KeyFile string `json:"key_file,omitempty" yaml:"key_file,omitempty"`
	// CA证书文件路径
	CAFile string `json:"ca_file,omitempty" yaml:"ca_file,omitempty"`
	// 服务器名称
	ServerName string `json:"server_name,omitempty" yaml:"server_name,omitempty"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Mode: ModeSingle,
		Single: &SingleConfig{
			Addr: "localhost:6379",
			DB:   0,
		},
		Common: CommonConfig{
			PoolSize:        10,
			MinIdleConns:    2,
			MaxIdleConns:    5,
			ConnMaxAge:      time.Hour,
			PoolTimeout:     time.Second * 4,
			IdleTimeout:     time.Minute * 5,
			DialTimeout:     time.Second * 5,
			ReadTimeout:     time.Second * 3,
			WriteTimeout:    time.Second * 3,
			MaxRetries:      3,
			MinRetryBackoff: time.Millisecond * 8,
			MaxRetryBackoff: time.Millisecond * 512,
			DefaultTTL:      time.Hour * 24,
			EnableMetrics:   false,
			EnableTracing:   false,
		},
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.Mode == "" {
		return ErrInvalidMode
	}

	switch c.Mode {
	case ModeSingle:
		if c.Single == nil {
			return ErrMissingSingleConfig
		}
		if c.Single.Addr == "" {
			return ErrMissingAddr
		}
	case ModeCluster:
		if c.Cluster == nil {
			return ErrMissingClusterConfig
		}
		if len(c.Cluster.Addrs) == 0 {
			return ErrMissingAddrs
		}
	case ModeSentinel:
		if c.Sentinel == nil {
			return ErrMissingSentinelConfig
		}
		if len(c.Sentinel.Addrs) == 0 {
			return ErrMissingAddrs
		}
		if c.Sentinel.MasterName == "" {
			return ErrMissingMasterName
		}
	default:
		return ErrInvalidMode
	}

	return nil
}

// GetKeyWithPrefix 获取带前缀的键名
func (c *Config) GetKeyWithPrefix(key string) string {
	if c.Common.KeyPrefix == "" {
		return key
	}
	return c.Common.KeyPrefix + key
}

// GetTTL 获取TTL，如果指定了TTL则使用指定值，否则使用默认值
func (c *Config) GetTTL(ttl time.Duration) time.Duration {
	if ttl > 0 {
		return ttl
	}
	return c.Common.DefaultTTL
}