# Redis企业级封装包

一个功能完整、高性能的Go语言Redis客户端封装包，支持单机、集群、哨兵三种部署模式，提供统一的接口设计和企业级特性。

## 🚀 功能特性

### 多种部署模式支持
- **单机模式 (Single)**: 适用于开发环境和小规模应用
- **集群模式 (Cluster)**: 支持Redis Cluster，提供水平扩展能力
- **哨兵模式 (Sentinel)**: 提供高可用性和自动故障转移

### 统一接口设计
- 采用工厂模式创建客户端，对上层应用提供统一接口
- 底层实现细节对应用透明，便于切换不同部署模式
- 支持依赖注入和接口测试

### 键管理功能
- 支持全局键前缀配置，便于多个应用共享同一Redis实例
- 自动键名管理，避免键冲突

### 缓存管理
- 支持设置默认TTL(过期时间)，有效防止缓存无限增长
- 灵活的过期时间配置

### 完整的数据操作支持
- **字符串操作**: Get, Set, Incr, Decr等
- **哈希表操作**: HGet, HSet, HGetAll, HDel等
- **列表操作**: LPush, RPush, LPop, RPop, LRange等
- **集合操作**: SAdd, SRem, SMembers, SInter等
- **有序集合操作**: ZAdd, ZRem, ZRange, ZRank等
- **计数器操作**: 原子性递增递减
- **Lua脚本执行**: 支持复杂的原子操作
- **管道操作**: 批量命令执行，提高性能

## 📦 安装

```bash
go get github.com/redis/go-redis/v9
```

## 🏗️ 架构设计

```
┌─────────────────────────────────────────────────────────────┐
│                        应用层                                │
├─────────────────────────────────────────────────────────────┤
│                     统一客户端接口                            │
│                    (Client Interface)                      │
├─────────────────────────────────────────────────────────────┤
│                      工厂模式层                              │
│                    (Factory Pattern)                       │
├─────────────────────────────────────────────────────────────┤
│  单机客户端    │    集群客户端    │    哨兵客户端              │
│ SingleClient  │  ClusterClient  │  SentinelClient          │
├─────────────────────────────────────────────────────────────┤
│                    go-redis底层库                           │
│              (github.com/redis/go-redis/v9)               │
├─────────────────────────────────────────────────────────────┤
│              Redis服务器 (Single/Cluster/Sentinel)         │
└─────────────────────────────────────────────────────────────┘
```

## 🔧 快速开始

### 单机模式

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "cache"
)

func main() {
    // 创建单机模式配置
    config := &cache.Config{
        Mode: cache.ModeSingle,
        Single: &cache.SingleConfig{
            Addr: "localhost:6379",
            DB:   0,
        },
        Common: cache.CommonConfig{
            Password:     "",
            KeyPrefix:    "myapp:",
            DefaultTTL:   30 * time.Minute,
            PoolSize:     10,
            MinIdleConns: 5,
        },
    }
    
    // 创建客户端
    factory, err := cache.NewFactory(config)
    if err != nil {
        log.Fatal(err)
    }
    
    client, err := factory.CreateClient()
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
    
    ctx := context.Background()
    
    // 基本操作
    err = client.Set(ctx, "key", "value", time.Hour)
    if err != nil {
        log.Fatal(err)
    }
    
    value, err := client.Get(ctx, "key")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Value: %s\n", value)
}
```

### 集群模式

```go
config := &cache.Config{
    Mode: cache.ModeCluster,
    Cluster: &cache.ClusterConfig{
        Addrs: []string{
            "localhost:7000",
            "localhost:7001", 
            "localhost:7002",
        },
    },
    Common: cache.CommonConfig{
        KeyPrefix:  "cluster_app:",
        DefaultTTL: 1 * time.Hour,
    },
}
```

### 哨兵模式

```go
config := &cache.Config{
    Mode: cache.ModeSentinel,
    Sentinel: &cache.SentinelConfig{
        Addrs: []string{
            "localhost:26379",
            "localhost:26380",
            "localhost:26381",
        },
        MasterName: "mymaster",
        DB:         0,
    },
    Common: cache.CommonConfig{
        KeyPrefix:  "sentinel_app:",
        DefaultTTL: 30 * time.Minute,
    },
}
```

## 📚 详细使用示例

### 字符串操作

```go
// 设置值
err := client.Set(ctx, "user:1:name", "张三", time.Hour)

// 获取值
name, err := client.Get(ctx, "user:1:name")

// 递增计数器
count, err := client.Incr(ctx, "page:views")

// 递减计数器
count, err := client.Decr(ctx, "inventory:item:1")
```

### 哈希表操作

```go
// 设置哈希字段
err := client.HSet(ctx, "user:1", "name", "张三")
err = client.HSet(ctx, "user:1", "email", "zhangsan@example.com")

// 获取哈希字段
name, err := client.HGet(ctx, "user:1", "name")

// 获取所有哈希字段
user, err := client.HGetAll(ctx, "user:1")

// 删除哈希字段
err = client.HDel(ctx, "user:1", "email")
```

### 列表操作

```go
// 左侧推入
count, err := client.LPush(ctx, "tasks", "task1", "task2")

// 右侧推入
count, err := client.RPush(ctx, "logs", "log entry")

// 左侧弹出
task, err := client.LPop(ctx, "tasks")

// 获取范围
tasks, err := client.LRange(ctx, "tasks", 0, 10)
```

### 集合操作

```go
// 添加成员
count, err := client.SAdd(ctx, "tags", "go", "redis", "cache")

// 获取所有成员
members, err := client.SMembers(ctx, "tags")

// 交集
intersection, err := client.SInter(ctx, "tags:go", "tags:database")

// 删除成员
count, err := client.SRem(ctx, "tags", "cache")
```

### 有序集合操作

```go
// 添加成员
count, err := client.ZAdd(ctx, "leaderboard", cache.ZMember{
    Score:  100.5,
    Member: "player1",
})

// 获取排名范围
players, err := client.ZRange(ctx, "leaderboard", 0, 10)

// 获取成员排名
rank, err := client.ZRank(ctx, "leaderboard", "player1")

// 删除成员
count, err := client.ZRem(ctx, "leaderboard", "player1")
```

### 管道操作

```go
// 创建管道
pipe := client.Pipeline()

// 批量添加命令
pipe.Set(ctx, "key1", "value1", time.Hour)
pipe.Set(ctx, "key2", "value2", time.Hour)
pipe.Incr(ctx, "counter")

// 执行管道
results, err := pipe.Exec(ctx)
if err != nil {
    log.Fatal(err)
}

// 处理结果
for _, result := range results {
    fmt.Printf("Result: %v, Error: %v\n", result.Result, result.Error)
}
```

### Lua脚本执行

```go
// 原子性递增并设置过期时间的脚本
script := `
local key = KEYS[1]
local increment = ARGV[1]
local ttl = ARGV[2]

local current = redis.call('INCR', key)
redis.call('EXPIRE', key, ttl)
return current
`

result, err := client.Eval(ctx, script, []string{"counter:api"}, 1, 3600)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Counter value: %v\n", result)
```

## ⚙️ 配置选项

### 通用配置 (CommonConfig)

```go
type CommonConfig struct {
    Password     string        // Redis密码
    KeyPrefix    string        // 键前缀
    DefaultTTL   time.Duration // 默认过期时间
    PoolSize     int          // 连接池大小
    MinIdleConns int          // 最小空闲连接数
    DialTimeout  time.Duration // 连接超时
    ReadTimeout  time.Duration // 读取超时
    WriteTimeout time.Duration // 写入超时
    MaxRetries   int          // 最大重试次数
    TLS          *TLSConfig   // TLS配置
}
```

### TLS配置

```go
config := &cache.Config{
    // ... 其他配置
    Common: cache.CommonConfig{
        TLS: &cache.TLSConfig{
            Enabled:            true,
            CertFile:           "/path/to/cert.pem",
            KeyFile:            "/path/to/key.pem",
            CAFile:             "/path/to/ca.pem",
            InsecureSkipVerify: false,
        },
    },
}
```

## 🧪 运行示例

项目提供了完整的使用示例：

```bash
# 单机模式示例
cd examples/single
go run main.go

# 集群模式示例  
cd examples/cluster
go run main.go

# 哨兵模式示例
cd examples/sentinel
go run main.go
```

## 🔍 错误处理

包提供了完整的错误处理机制：

```go
// 检查是否为Redis相关错误
if cache.IsRedisError(err) {
    log.Printf("Redis错误: %v", err)
}

// 检查是否为连接错误
if cache.IsConnectionError(err) {
    log.Printf("连接错误: %v", err)
}

// 检查是否为配置错误
if errors.Is(err, cache.ErrInvalidMode) {
    log.Printf("配置模式无效")
}
```

## 🏆 最佳实践

### 1. 连接池配置

```go
// 根据应用负载调整连接池大小
Common: cache.CommonConfig{
    PoolSize:     20,  // 高并发应用可适当增大
    MinIdleConns: 10,  // 保持一定数量的空闲连接
    DialTimeout:  5 * time.Second,
    ReadTimeout:  3 * time.Second,
    WriteTimeout: 3 * time.Second,
}
```

### 2. 键命名规范

```go
// 使用有意义的键前缀
KeyPrefix: "myapp:v1:"

// 推荐的键命名模式：
// - 用户数据: "user:123:profile"
// - 会话数据: "session:abc123"
// - 缓存数据: "cache:product:456"
// - 计数器: "counter:api:daily"
```

### 3. TTL管理

```go
// 设置合理的默认TTL
DefaultTTL: 30 * time.Minute,

// 根据数据特性设置不同的TTL
client.Set(ctx, "user:session", token, 2*time.Hour)     // 会话数据
client.Set(ctx, "cache:product", data, 10*time.Minute) // 缓存数据
client.Set(ctx, "config:app", config, 24*time.Hour)    // 配置数据
```

### 4. 错误处理

```go
// 优雅处理Redis不可用的情况
value, err := client.Get(ctx, "key")
if err != nil {
    if cache.IsConnectionError(err) {
        // 使用降级策略，如从数据库读取
        return getFromDatabase(key)
    }
    return "", err
}
```

### 5. 性能优化

```go
// 使用管道批量操作
pipe := client.Pipeline()
for _, item := range items {
    pipe.Set(ctx, item.Key, item.Value, time.Hour)
}
results, err := pipe.Exec(ctx)

// 使用Lua脚本实现原子操作
script := `
-- 原子性的获取并递增计数器
local current = redis.call('GET', KEYS[1]) or 0
redis.call('INCR', KEYS[1])
return tonumber(current)
`
result, err := client.Eval(ctx, script, []string{"counter"}, )
```

## 📈 性能特性

- **连接池管理**: 自动管理连接池，支持连接复用
- **管道操作**: 批量命令执行，减少网络往返
- **Lua脚本**: 服务端原子操作，提高性能
- **智能重试**: 自动重试机制，提高可靠性
- **内存优化**: 合理的内存使用和垃圾回收

## 🛡️ 安全特性

- **TLS支持**: 支持加密传输
- **认证机制**: 支持密码认证
- **键前缀隔离**: 多应用安全共享
- **连接安全**: 连接超时和重试保护

## 🔧 故障排除

### 常见问题

1. **连接超时**
   ```go
   // 增加连接超时时间
   DialTimeout: 10 * time.Second,
   ```

2. **内存使用过高**
   ```go
   // 设置合理的TTL
   DefaultTTL: 30 * time.Minute,
   ```

3. **连接池耗尽**
   ```go
   // 增加连接池大小
   PoolSize: 50,
   ```

### 监控指标

建议监控以下指标：
- 连接池使用率
- 命令执行延迟
- 错误率
- 内存使用情况
- 网络流量

## 📁 项目结构

```
cache/
├── README.md              # 项目文档
├── go.mod                 # Go模块定义
├── client.go              # 核心接口定义
├── config.go              # 配置结构体
├── errors.go              # 错误定义
├── factory.go             # 工厂模式实现
├── single_client.go       # 单机模式客户端
├── cluster_client.go      # 集群模式客户端
├── sentinel_client.go     # 哨兵模式客户端
├── pipeliner.go           # 管道操作实现
└── examples/              # 使用示例
    ├── single/            # 单机模式示例
    ├── cluster/           # 集群模式示例
    └── sentinel/          # 哨兵模式示例
```

## 📄 许可证

MIT License

## 🤝 贡献

欢迎提交Issue和Pull Request来改进这个项目。

## 📞 支持

如有问题，请提交Issue或联系维护者。