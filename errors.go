package cache

import "errors"

// 配置相关错误
var (
	// ErrInvalidMode 无效的部署模式
	ErrInvalidMode = errors.New("invalid redis mode")
	// ErrMissingSingleConfig 缺少单机模式配置
	ErrMissingSingleConfig = errors.New("missing single mode config")
	// ErrMissingClusterConfig 缺少集群模式配置
	ErrMissingClusterConfig = errors.New("missing cluster mode config")
	// ErrMissingSentinelConfig 缺少哨兵模式配置
	ErrMissingSentinelConfig = errors.New("missing sentinel mode config")
	// ErrMissingAddr 缺少服务器地址
	ErrMissingAddr = errors.New("missing server address")
	// ErrMissingAddrs 缺少服务器地址列表
	ErrMissingAddrs = errors.New("missing server addresses")
	// ErrMissingMasterName 缺少主节点名称
	ErrMissingMasterName = errors.New("missing master name")
)

// 客户端操作相关错误
var (
	// ErrClientClosed 客户端已关闭
	ErrClientClosed = errors.New("redis client is closed")
	// ErrNilResult 空结果
	ErrNilResult = errors.New("redis: nil result")
	// ErrKeyNotFound 键不存在
	ErrKeyNotFound = errors.New("redis: key not found")
	// ErrInvalidType 无效的数据类型
	ErrInvalidType = errors.New("redis: invalid data type")
	// ErrScriptNotFound Lua脚本不存在
	ErrScriptNotFound = errors.New("redis: script not found")
)

// 连接相关错误
var (
	// ErrConnectionFailed 连接失败
	ErrConnectionFailed = errors.New("redis: connection failed")
	// ErrConnectionTimeout 连接超时
	ErrConnectionTimeout = errors.New("redis: connection timeout")
	// ErrPoolExhausted 连接池耗尽
	ErrPoolExhausted = errors.New("redis: connection pool exhausted")
	// ErrAuthFailed 认证失败
	ErrAuthFailed = errors.New("redis: authentication failed")
)

// 集群相关错误
var (
	// ErrClusterDown 集群不可用
	ErrClusterDown = errors.New("redis: cluster is down")
	// ErrNoReachableNode 没有可达的节点
	ErrNoReachableNode = errors.New("redis: no reachable cluster node")
	// ErrTooManyRedirects 重定向次数过多
	ErrTooManyRedirects = errors.New("redis: too many cluster redirects")
)

// 哨兵相关错误
var (
	// ErrSentinelNoMaster 哨兵找不到主节点
	ErrSentinelNoMaster = errors.New("redis: sentinel no master found")
	// ErrSentinelMasterDown 主节点不可用
	ErrSentinelMasterDown = errors.New("redis: sentinel master is down")
	// ErrNoSentinelAvailable 没有可用的哨兵
	ErrNoSentinelAvailable = errors.New("redis: no sentinel available")
)

// 管道相关错误
var (
	// ErrPipelineEmpty 管道为空
	ErrPipelineEmpty = errors.New("redis: pipeline is empty")
	// ErrPipelineClosed 管道已关闭
	ErrPipelineClosed = errors.New("redis: pipeline is closed")
)

// IsRedisError 判断是否为Redis相关错误
func IsRedisError(err error) bool {
	if err == nil {
		return false
	}

	// 检查是否为已定义的Redis错误
	redisErrors := []error{
		ErrInvalidMode, ErrMissingSingleConfig, ErrMissingClusterConfig,
		ErrMissingSentinelConfig, ErrMissingAddr, ErrMissingAddrs,
		ErrMissingMasterName, ErrClientClosed, ErrNilResult,
		ErrKeyNotFound, ErrInvalidType, ErrScriptNotFound,
		ErrConnectionFailed, ErrConnectionTimeout, ErrPoolExhausted,
		ErrAuthFailed, ErrClusterDown, ErrNoReachableNode,
		ErrTooManyRedirects, ErrSentinelNoMaster, ErrSentinelMasterDown,
		ErrNoSentinelAvailable, ErrPipelineEmpty, ErrPipelineClosed,
	}

	for _, redisErr := range redisErrors {
		if errors.Is(err, redisErr) {
			return true
		}
	}

	return false
}

// IsConnectionError 判断是否为连接相关错误
func IsConnectionError(err error) bool {
	if err == nil {
		return false
	}

	connectionErrors := []error{
		ErrConnectionFailed, ErrConnectionTimeout,
		ErrPoolExhausted, ErrAuthFailed,
	}

	for _, connErr := range connectionErrors {
		if errors.Is(err, connErr) {
			return true
		}
	}

	return false
}

// IsClusterError 判断是否为集群相关错误
func IsClusterError(err error) bool {
	if err == nil {
		return false
	}

	clusterErrors := []error{
		ErrClusterDown, ErrNoReachableNode, ErrTooManyRedirects,
	}

	for _, clusterErr := range clusterErrors {
		if errors.Is(err, clusterErr) {
			return true
		}
	}

	return false
}

// IsSentinelError 判断是否为哨兵相关错误
func IsSentinelError(err error) bool {
	if err == nil {
		return false
	}

	sentinelErrors := []error{
		ErrSentinelNoMaster, ErrSentinelMasterDown, ErrNoSentinelAvailable,
	}

	for _, sentinelErr := range sentinelErrors {
		if errors.Is(err, sentinelErr) {
			return true
		}
	}

	return false
}