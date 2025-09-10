package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// SingleClient 单机模式Redis客户端
type SingleClient struct {
	client *redis.Client
	config *Config
}

// Close 关闭客户端连接
func (c *SingleClient) Close() error {
	return c.client.Close()
}

// Ping 测试连接
func (c *SingleClient) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

// 字符串操作

// Get 获取字符串值
func (c *SingleClient) Get(ctx context.Context, key string) (string, error) {
	key = c.config.GetKeyWithPrefix(key)
	result := c.client.Get(ctx, key)
	if result.Err() == redis.Nil {
		return "", ErrKeyNotFound
	}
	return result.Result()
}

// Set 设置字符串值
func (c *SingleClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	key = c.config.GetKeyWithPrefix(key)
	expiration = c.config.GetTTL(expiration)
	return c.client.Set(ctx, key, value, expiration).Err()
}

// SetNX 仅当键不存在时设置值
func (c *SingleClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	key = c.config.GetKeyWithPrefix(key)
	expiration = c.config.GetTTL(expiration)
	return c.client.SetNX(ctx, key, value, expiration).Result()
}

// GetSet 设置新值并返回旧值
func (c *SingleClient) GetSet(ctx context.Context, key string, value interface{}) (string, error) {
	key = c.config.GetKeyWithPrefix(key)
	result := c.client.GetSet(ctx, key, value)
	if result.Err() == redis.Nil {
		return "", ErrKeyNotFound
	}
	return result.Result()
}

// MGet 批量获取多个键的值
func (c *SingleClient) MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	prefixedKeys := make([]string, len(keys))
	for i, key := range keys {
		prefixedKeys[i] = c.config.GetKeyWithPrefix(key)
	}
	return c.client.MGet(ctx, prefixedKeys...).Result()
}

// MSet 批量设置多个键值对
func (c *SingleClient) MSet(ctx context.Context, pairs ...interface{}) error {
	// 为键添加前缀
	prefixedPairs := make([]interface{}, len(pairs))
	for i := 0; i < len(pairs); i += 2 {
		if i+1 < len(pairs) {
			key := pairs[i].(string)
			prefixedPairs[i] = c.config.GetKeyWithPrefix(key)
			prefixedPairs[i+1] = pairs[i+1]
		}
	}
	return c.client.MSet(ctx, prefixedPairs...).Err()
}

// Incr 递增计数器
func (c *SingleClient) Incr(ctx context.Context, key string) (int64, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.Incr(ctx, key).Result()
}

// IncrBy 按指定值递增计数器
func (c *SingleClient) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.IncrBy(ctx, key, value).Result()
}

// Decr 递减计数器
func (c *SingleClient) Decr(ctx context.Context, key string) (int64, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.Decr(ctx, key).Result()
}

// DecrBy 按指定值递减计数器
func (c *SingleClient) DecrBy(ctx context.Context, key string, value int64) (int64, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.DecrBy(ctx, key, value).Result()
}

// 哈希表操作

// HGet 获取哈希表字段值
func (c *SingleClient) HGet(ctx context.Context, key, field string) (string, error) {
	key = c.config.GetKeyWithPrefix(key)
	result := c.client.HGet(ctx, key, field)
	if result.Err() == redis.Nil {
		return "", ErrKeyNotFound
	}
	return result.Result()
}

// HSet 设置哈希表字段值
func (c *SingleClient) HSet(ctx context.Context, key, field string, value interface{}) error {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.HSet(ctx, key, field, value).Err()
}

// HSetNX 仅当字段不存在时设置哈希表字段值
func (c *SingleClient) HSetNX(ctx context.Context, key, field string, value interface{}) (bool, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.HSetNX(ctx, key, field, value).Result()
}

// HDel 删除哈希表字段
func (c *SingleClient) HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.HDel(ctx, key, fields...).Result()
}

// HExists 检查哈希表字段是否存在
func (c *SingleClient) HExists(ctx context.Context, key, field string) (bool, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.HExists(ctx, key, field).Result()
}

// HGetAll 获取哈希表所有字段和值
func (c *SingleClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.HGetAll(ctx, key).Result()
}

// HKeys 获取哈希表所有字段
func (c *SingleClient) HKeys(ctx context.Context, key string) ([]string, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.HKeys(ctx, key).Result()
}

// HVals 获取哈希表所有值
func (c *SingleClient) HVals(ctx context.Context, key string) ([]string, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.HVals(ctx, key).Result()
}

// HLen 获取哈希表字段数量
func (c *SingleClient) HLen(ctx context.Context, key string) (int64, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.HLen(ctx, key).Result()
}

// HMGet 批量获取哈希表字段值
func (c *SingleClient) HMGet(ctx context.Context, key string, fields ...string) ([]interface{}, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.HMGet(ctx, key, fields...).Result()
}

// HMSet 批量设置哈希表字段值
func (c *SingleClient) HMSet(ctx context.Context, key string, pairs ...interface{}) error {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.HMSet(ctx, key, pairs...).Err()
}

// HIncrBy 递增哈希表字段值
func (c *SingleClient) HIncrBy(ctx context.Context, key, field string, incr int64) (int64, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.HIncrBy(ctx, key, field, incr).Result()
}

// 列表操作

// LPush 从列表左侧推入元素
func (c *SingleClient) LPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.LPush(ctx, key, values...).Result()
}

// RPush 从列表右侧推入元素
func (c *SingleClient) RPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.RPush(ctx, key, values...).Result()
}

// LPop 从列表左侧弹出元素
func (c *SingleClient) LPop(ctx context.Context, key string) (string, error) {
	key = c.config.GetKeyWithPrefix(key)
	result := c.client.LPop(ctx, key)
	if result.Err() == redis.Nil {
		return "", ErrKeyNotFound
	}
	return result.Result()
}

// RPop 从列表右侧弹出元素
func (c *SingleClient) RPop(ctx context.Context, key string) (string, error) {
	key = c.config.GetKeyWithPrefix(key)
	result := c.client.RPop(ctx, key)
	if result.Err() == redis.Nil {
		return "", ErrKeyNotFound
	}
	return result.Result()
}

// LLen 获取列表长度
func (c *SingleClient) LLen(ctx context.Context, key string) (int64, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.LLen(ctx, key).Result()
}

// LRange 获取列表指定范围的元素
func (c *SingleClient) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.LRange(ctx, key, start, stop).Result()
}

// LIndex 获取列表指定索引的元素
func (c *SingleClient) LIndex(ctx context.Context, key string, index int64) (string, error) {
	key = c.config.GetKeyWithPrefix(key)
	result := c.client.LIndex(ctx, key, index)
	if result.Err() == redis.Nil {
		return "", ErrKeyNotFound
	}
	return result.Result()
}

// LSet 设置列表指定索引的元素值
func (c *SingleClient) LSet(ctx context.Context, key string, index int64, value interface{}) error {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.LSet(ctx, key, index, value).Err()
}

// LRem 从列表中移除元素
func (c *SingleClient) LRem(ctx context.Context, key string, count int64, value interface{}) (int64, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.LRem(ctx, key, count, value).Result()
}

// LTrim 修剪列表，只保留指定范围的元素
func (c *SingleClient) LTrim(ctx context.Context, key string, start, stop int64) error {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.LTrim(ctx, key, start, stop).Err()
}

// 集合操作

// SAdd 向集合添加成员
func (c *SingleClient) SAdd(ctx context.Context, key string, members ...interface{}) (int64, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.SAdd(ctx, key, members...).Result()
}

// SRem 从集合移除成员
func (c *SingleClient) SRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.SRem(ctx, key, members...).Result()
}

// SMembers 获取集合所有成员
func (c *SingleClient) SMembers(ctx context.Context, key string) ([]string, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.SMembers(ctx, key).Result()
}

// SIsMember 检查成员是否在集合中
func (c *SingleClient) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.SIsMember(ctx, key, member).Result()
}

// SCard 获取集合成员数量
func (c *SingleClient) SCard(ctx context.Context, key string) (int64, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.SCard(ctx, key).Result()
}

// SPop 随机移除并返回集合中的一个成员
func (c *SingleClient) SPop(ctx context.Context, key string) (string, error) {
	key = c.config.GetKeyWithPrefix(key)
	result := c.client.SPop(ctx, key)
	if result.Err() == redis.Nil {
		return "", ErrKeyNotFound
	}
	return result.Result()
}

// SRandMember 随机返回集合中的一个成员
func (c *SingleClient) SRandMember(ctx context.Context, key string) (string, error) {
	key = c.config.GetKeyWithPrefix(key)
	result := c.client.SRandMember(ctx, key)
	if result.Err() == redis.Nil {
		return "", ErrKeyNotFound
	}
	return result.Result()
}

// SInter 计算多个集合的交集
func (c *SingleClient) SInter(ctx context.Context, keys ...string) ([]string, error) {
	prefixedKeys := make([]string, len(keys))
	for i, key := range keys {
		prefixedKeys[i] = c.config.GetKeyWithPrefix(key)
	}
	return c.client.SInter(ctx, prefixedKeys...).Result()
}

// SUnion 计算多个集合的并集
func (c *SingleClient) SUnion(ctx context.Context, keys ...string) ([]string, error) {
	prefixedKeys := make([]string, len(keys))
	for i, key := range keys {
		prefixedKeys[i] = c.config.GetKeyWithPrefix(key)
	}
	return c.client.SUnion(ctx, prefixedKeys...).Result()
}

// SDiff 计算多个集合的差集
func (c *SingleClient) SDiff(ctx context.Context, keys ...string) ([]string, error) {
	prefixedKeys := make([]string, len(keys))
	for i, key := range keys {
		prefixedKeys[i] = c.config.GetKeyWithPrefix(key)
	}
	return c.client.SDiff(ctx, prefixedKeys...).Result()
}

// 有序集合操作

// ZAdd 向有序集合添加成员
func (c *SingleClient) ZAdd(ctx context.Context, key string, members ...ZMember) (int64, error) {
	key = c.config.GetKeyWithPrefix(key)
	redisMembers := make([]redis.Z, len(members))
	for i, member := range members {
		redisMembers[i] = redis.Z{
			Score:  member.Score,
			Member: member.Member,
		}
	}
	return c.client.ZAdd(ctx, key, redisMembers...).Result()
}

// ZRem 从有序集合移除成员
func (c *SingleClient) ZRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.ZRem(ctx, key, members...).Result()
}

// ZScore 获取有序集合成员的分数
func (c *SingleClient) ZScore(ctx context.Context, key, member string) (float64, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.ZScore(ctx, key, member).Result()
}

// ZRank 获取有序集合成员的排名（从小到大）
func (c *SingleClient) ZRank(ctx context.Context, key, member string) (int64, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.ZRank(ctx, key, member).Result()
}

// ZRevRank 获取有序集合成员的排名（从大到小）
func (c *SingleClient) ZRevRank(ctx context.Context, key, member string) (int64, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.ZRevRank(ctx, key, member).Result()
}

// ZRange 获取有序集合指定范围的成员（从小到大）
func (c *SingleClient) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.ZRange(ctx, key, start, stop).Result()
}

// ZRevRange 获取有序集合指定范围的成员（从大到小）
func (c *SingleClient) ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.ZRevRange(ctx, key, start, stop).Result()
}

// ZRangeWithScores 获取有序集合指定范围的成员和分数（从小到大）
func (c *SingleClient) ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]ZMember, error) {
	key = c.config.GetKeyWithPrefix(key)
	result, err := c.client.ZRangeWithScores(ctx, key, start, stop).Result()
	if err != nil {
		return nil, err
	}

	members := make([]ZMember, len(result))
	for i, z := range result {
		members[i] = ZMember{
			Score:  z.Score,
			Member: z.Member,
		}
	}
	return members, nil
}

// ZRevRangeWithScores 获取有序集合指定范围的成员和分数（从大到小）
func (c *SingleClient) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) ([]ZMember, error) {
	key = c.config.GetKeyWithPrefix(key)
	result, err := c.client.ZRevRangeWithScores(ctx, key, start, stop).Result()
	if err != nil {
		return nil, err
	}

	members := make([]ZMember, len(result))
	for i, z := range result {
		members[i] = ZMember{
			Score:  z.Score,
			Member: z.Member,
		}
	}
	return members, nil
}

// ZRangeByScore 根据分数范围获取有序集合成员
func (c *SingleClient) ZRangeByScore(ctx context.Context, key string, min, max string) ([]string, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: min,
		Max: max,
	}).Result()
}

// ZRevRangeByScore 根据分数范围获取有序集合成员（逆序）
func (c *SingleClient) ZRevRangeByScore(ctx context.Context, key string, max, min string) ([]string, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.ZRevRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: min,
		Max: max,
	}).Result()
}

// ZCard 获取有序集合成员数量
func (c *SingleClient) ZCard(ctx context.Context, key string) (int64, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.ZCard(ctx, key).Result()
}

// ZCount 计算指定分数范围内的成员数量
func (c *SingleClient) ZCount(ctx context.Context, key, min, max string) (int64, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.ZCount(ctx, key, min, max).Result()
}

// ZIncrBy 增加有序集合成员的分数
func (c *SingleClient) ZIncrBy(ctx context.Context, key string, increment float64, member string) (float64, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.ZIncrBy(ctx, key, increment, member).Result()
}

// 通用键操作

// Del 删除键
func (c *SingleClient) Del(ctx context.Context, keys ...string) (int64, error) {
	prefixedKeys := make([]string, len(keys))
	for i, key := range keys {
		prefixedKeys[i] = c.config.GetKeyWithPrefix(key)
	}
	return c.client.Del(ctx, prefixedKeys...).Result()
}

// Exists 检查键是否存在
func (c *SingleClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	prefixedKeys := make([]string, len(keys))
	for i, key := range keys {
		prefixedKeys[i] = c.config.GetKeyWithPrefix(key)
	}
	return c.client.Exists(ctx, prefixedKeys...).Result()
}

// Expire 设置键的过期时间
func (c *SingleClient) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.Expire(ctx, key, expiration).Result()
}

// ExpireAt 设置键在指定时间过期
func (c *SingleClient) ExpireAt(ctx context.Context, key string, tm time.Time) (bool, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.ExpireAt(ctx, key, tm).Result()
}

// TTL 获取键的剩余生存时间
func (c *SingleClient) TTL(ctx context.Context, key string) (time.Duration, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.TTL(ctx, key).Result()
}

// Type 获取键的数据类型
func (c *SingleClient) Type(ctx context.Context, key string) (string, error) {
	key = c.config.GetKeyWithPrefix(key)
	return c.client.Type(ctx, key).Result()
}

// Keys 查找匹配模式的键
func (c *SingleClient) Keys(ctx context.Context, pattern string) ([]string, error) {
	pattern = c.config.GetKeyWithPrefix(pattern)
	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	// 移除前缀
	if c.config.Common.KeyPrefix != "" {
		prefixLen := len(c.config.Common.KeyPrefix)
		for i, key := range keys {
			if len(key) > prefixLen {
				keys[i] = key[prefixLen:]
			}
		}
	}

	return keys, nil
}

// Scan 迭代数据库中的键
func (c *SingleClient) Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error) {
	match = c.config.GetKeyWithPrefix(match)
	keys, cursor, err := c.client.Scan(ctx, cursor, match, count).Result()
	if err != nil {
		return nil, cursor, err
	}

	// 移除前缀
	if c.config.Common.KeyPrefix != "" {
		prefixLen := len(c.config.Common.KeyPrefix)
		for i, key := range keys {
			if len(key) > prefixLen {
				keys[i] = key[prefixLen:]
			}
		}
	}

	return keys, cursor, nil
}

// Lua脚本操作

// Eval 执行Lua脚本
func (c *SingleClient) Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error) {
	prefixedKeys := make([]string, len(keys))
	for i, key := range keys {
		prefixedKeys[i] = c.config.GetKeyWithPrefix(key)
	}
	return c.client.Eval(ctx, script, prefixedKeys, args...).Result()
}

// EvalSha 通过SHA1执行Lua脚本
func (c *SingleClient) EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) (interface{}, error) {
	prefixedKeys := make([]string, len(keys))
	for i, key := range keys {
		prefixedKeys[i] = c.config.GetKeyWithPrefix(key)
	}
	return c.client.EvalSha(ctx, sha1, prefixedKeys, args...).Result()
}

// ScriptExists 检查脚本是否存在
func (c *SingleClient) ScriptExists(ctx context.Context, hashes ...string) ([]bool, error) {
	return c.client.ScriptExists(ctx, hashes...).Result()
}

// ScriptFlush 清空脚本缓存
func (c *SingleClient) ScriptFlush(ctx context.Context) error {
	return c.client.ScriptFlush(ctx).Err()
}

// ScriptKill 终止正在执行的脚本
func (c *SingleClient) ScriptKill(ctx context.Context) error {
	return c.client.ScriptKill(ctx).Err()
}

// ScriptLoad 加载脚本到缓存
func (c *SingleClient) ScriptLoad(ctx context.Context, script string) (string, error) {
	return c.client.ScriptLoad(ctx, script).Result()
}

// 管道操作

// Pipeline 创建管道
func (c *SingleClient) Pipeline() Pipeliner {
	return &SinglePipeliner{
		pipe:   c.client.Pipeline(),
		config: c.config,
	}
}

// TxPipeline 创建事务管道
func (c *SingleClient) TxPipeline() Pipeliner {
	return &SinglePipeliner{
		pipe:   c.client.TxPipeline(),
		config: c.config,
	}
}