package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// SinglePipeliner 单机模式管道实现
type SinglePipeliner struct {
	pipe   redis.Pipeliner
	config *Config
}

// Exec 执行管道中的所有命令
func (p *SinglePipeliner) Exec(ctx context.Context) ([]interface{}, error) {
	cmds, err := p.pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]interface{}, len(cmds))
	for i, cmd := range cmds {
		results[i] = cmd
	}
	return results, nil
}

// Discard 丢弃管道中的所有命令
func (p *SinglePipeliner) Discard() error {
	p.pipe.Discard()
	return nil
}

// Close 关闭管道
func (p *SinglePipeliner) Close() error {
	// Redis管道不需要显式关闭
	return nil
}

// 字符串操作

// Get 获取字符串值
func (p *SinglePipeliner) Get(ctx context.Context, key string) *StringCmd {
	key = p.config.GetKeyWithPrefix(key)
	cmd := p.pipe.Get(ctx, key)
	return &StringCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// Set 设置字符串值
func (p *SinglePipeliner) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *StatusCmd {
	key = p.config.GetKeyWithPrefix(key)
	expiration = p.config.GetTTL(expiration)
	cmd := p.pipe.Set(ctx, key, value, expiration)
	return &StatusCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// SetNX 仅当键不存在时设置值
func (p *SinglePipeliner) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *BoolCmd {
	key = p.config.GetKeyWithPrefix(key)
	expiration = p.config.GetTTL(expiration)
	cmd := p.pipe.SetNX(ctx, key, value, expiration)
	return &BoolCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// Incr 递增计数器
func (p *SinglePipeliner) Incr(ctx context.Context, key string) *IntCmd {
	key = p.config.GetKeyWithPrefix(key)
	cmd := p.pipe.Incr(ctx, key)
	return &IntCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// Decr 递减计数器
func (p *SinglePipeliner) Decr(ctx context.Context, key string) *IntCmd {
	key = p.config.GetKeyWithPrefix(key)
	cmd := p.pipe.Decr(ctx, key)
	return &IntCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// 哈希表操作

// HGet 获取哈希表字段值
func (p *SinglePipeliner) HGet(ctx context.Context, key, field string) *StringCmd {
	key = p.config.GetKeyWithPrefix(key)
	cmd := p.pipe.HGet(ctx, key, field)
	return &StringCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// HSet 设置哈希表字段值
func (p *SinglePipeliner) HSet(ctx context.Context, key, field string, value interface{}) *IntCmd {
	key = p.config.GetKeyWithPrefix(key)
	cmd := p.pipe.HSet(ctx, key, field, value)
	return &IntCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// HDel 删除哈希表字段
func (p *SinglePipeliner) HDel(ctx context.Context, key string, fields ...string) *IntCmd {
	key = p.config.GetKeyWithPrefix(key)
	cmd := p.pipe.HDel(ctx, key, fields...)
	return &IntCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// 列表操作

// LPush 从列表左侧推入元素
func (p *SinglePipeliner) LPush(ctx context.Context, key string, values ...interface{}) *IntCmd {
	key = p.config.GetKeyWithPrefix(key)
	cmd := p.pipe.LPush(ctx, key, values...)
	return &IntCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// RPush 从列表右侧推入元素
func (p *SinglePipeliner) RPush(ctx context.Context, key string, values ...interface{}) *IntCmd {
	key = p.config.GetKeyWithPrefix(key)
	cmd := p.pipe.RPush(ctx, key, values...)
	return &IntCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// LPop 从列表左侧弹出元素
func (p *SinglePipeliner) LPop(ctx context.Context, key string) *StringCmd {
	key = p.config.GetKeyWithPrefix(key)
	cmd := p.pipe.LPop(ctx, key)
	return &StringCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// RPop 从列表右侧弹出元素
func (p *SinglePipeliner) RPop(ctx context.Context, key string) *StringCmd {
	key = p.config.GetKeyWithPrefix(key)
	cmd := p.pipe.RPop(ctx, key)
	return &StringCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// 集合操作

// SAdd 向集合添加成员
func (p *SinglePipeliner) SAdd(ctx context.Context, key string, members ...interface{}) *IntCmd {
	key = p.config.GetKeyWithPrefix(key)
	cmd := p.pipe.SAdd(ctx, key, members...)
	return &IntCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// SRem 从集合移除成员
func (p *SinglePipeliner) SRem(ctx context.Context, key string, members ...interface{}) *IntCmd {
	key = p.config.GetKeyWithPrefix(key)
	cmd := p.pipe.SRem(ctx, key, members...)
	return &IntCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// SMembers 获取集合所有成员
func (p *SinglePipeliner) SMembers(ctx context.Context, key string) *StringSliceCmd {
	key = p.config.GetKeyWithPrefix(key)
	cmd := p.pipe.SMembers(ctx, key)
	return &StringSliceCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// 有序集合操作

// ZAdd 向有序集合添加成员
func (p *SinglePipeliner) ZAdd(ctx context.Context, key string, members ...ZMember) *IntCmd {
	key = p.config.GetKeyWithPrefix(key)
	redisMembers := make([]redis.Z, len(members))
	for i, member := range members {
		redisMembers[i] = redis.Z{
			Score:  member.Score,
			Member: member.Member,
		}
	}
	cmd := p.pipe.ZAdd(ctx, key, redisMembers...)
	return &IntCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// ZRem 从有序集合移除成员
func (p *SinglePipeliner) ZRem(ctx context.Context, key string, members ...interface{}) *IntCmd {
	key = p.config.GetKeyWithPrefix(key)
	cmd := p.pipe.ZRem(ctx, key, members...)
	return &IntCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// ZRange 获取有序集合指定范围的成员（从小到大）
func (p *SinglePipeliner) ZRange(ctx context.Context, key string, start, stop int64) *StringSliceCmd {
	key = p.config.GetKeyWithPrefix(key)
	cmd := p.pipe.ZRange(ctx, key, start, stop)
	return &StringSliceCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// 通用操作

// Del 删除键
func (p *SinglePipeliner) Del(ctx context.Context, keys ...string) *IntCmd {
	prefixedKeys := make([]string, len(keys))
	for i, key := range keys {
		prefixedKeys[i] = p.config.GetKeyWithPrefix(key)
	}
	cmd := p.pipe.Del(ctx, prefixedKeys...)
	return &IntCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// Exists 检查键是否存在
func (p *SinglePipeliner) Exists(ctx context.Context, keys ...string) *IntCmd {
	prefixedKeys := make([]string, len(keys))
	for i, key := range keys {
		prefixedKeys[i] = p.config.GetKeyWithPrefix(key)
	}
	cmd := p.pipe.Exists(ctx, prefixedKeys...)
	return &IntCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// Expire 设置键的过期时间
func (p *SinglePipeliner) Expire(ctx context.Context, key string, expiration time.Duration) *BoolCmd {
	key = p.config.GetKeyWithPrefix(key)
	cmd := p.pipe.Expire(ctx, key, expiration)
	return &BoolCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// ClusterPipeliner 集群模式管道实现
type ClusterPipeliner struct {
	pipe   redis.Pipeliner
	config *Config
}

// Exec 执行管道中的所有命令
func (p *ClusterPipeliner) Exec(ctx context.Context) ([]interface{}, error) {
	cmds, err := p.pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]interface{}, len(cmds))
	for i, cmd := range cmds {
		results[i] = cmd
	}
	return results, nil
}

// Discard 丢弃管道中的所有命令
func (p *ClusterPipeliner) Discard() error {
	p.pipe.Discard()
	return nil
}

// Close 关闭管道
func (p *ClusterPipeliner) Close() error {
	// Redis管道不需要显式关闭
	return nil
}

// 字符串操作

// Get 获取字符串值
func (p *ClusterPipeliner) Get(ctx context.Context, key string) *StringCmd {
	key = p.config.GetKeyWithPrefix(key)
	cmd := p.pipe.Get(ctx, key)
	return &StringCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// Set 设置字符串值
func (p *ClusterPipeliner) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *StatusCmd {
	key = p.config.GetKeyWithPrefix(key)
	expiration = p.config.GetTTL(expiration)
	cmd := p.pipe.Set(ctx, key, value, expiration)
	return &StatusCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// SetNX 仅当键不存在时设置值
func (p *ClusterPipeliner) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *BoolCmd {
	key = p.config.GetKeyWithPrefix(key)
	expiration = p.config.GetTTL(expiration)
	cmd := p.pipe.SetNX(ctx, key, value, expiration)
	return &BoolCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// Incr 递增计数器
func (p *ClusterPipeliner) Incr(ctx context.Context, key string) *IntCmd {
	key = p.config.GetKeyWithPrefix(key)
	cmd := p.pipe.Incr(ctx, key)
	return &IntCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// Decr 递减计数器
func (p *ClusterPipeliner) Decr(ctx context.Context, key string) *IntCmd {
	key = p.config.GetKeyWithPrefix(key)
	cmd := p.pipe.Decr(ctx, key)
	return &IntCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// 哈希表操作

// HGet 获取哈希表字段值
func (p *ClusterPipeliner) HGet(ctx context.Context, key, field string) *StringCmd {
	key = p.config.GetKeyWithPrefix(key)
	cmd := p.pipe.HGet(ctx, key, field)
	return &StringCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// HSet 设置哈希表字段值
func (p *ClusterPipeliner) HSet(ctx context.Context, key, field string, value interface{}) *IntCmd {
	key = p.config.GetKeyWithPrefix(key)
	cmd := p.pipe.HSet(ctx, key, field, value)
	return &IntCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// HDel 删除哈希表字段
func (p *ClusterPipeliner) HDel(ctx context.Context, key string, fields ...string) *IntCmd {
	key = p.config.GetKeyWithPrefix(key)
	cmd := p.pipe.HDel(ctx, key, fields...)
	return &IntCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// 列表操作

// LPush 从列表左侧推入元素
func (p *ClusterPipeliner) LPush(ctx context.Context, key string, values ...interface{}) *IntCmd {
	key = p.config.GetKeyWithPrefix(key)
	cmd := p.pipe.LPush(ctx, key, values...)
	return &IntCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// RPush 从列表右侧推入元素
func (p *ClusterPipeliner) RPush(ctx context.Context, key string, values ...interface{}) *IntCmd {
	key = p.config.GetKeyWithPrefix(key)
	cmd := p.pipe.RPush(ctx, key, values...)
	return &IntCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// LPop 从列表左侧弹出元素
func (p *ClusterPipeliner) LPop(ctx context.Context, key string) *StringCmd {
	key = p.config.GetKeyWithPrefix(key)
	cmd := p.pipe.LPop(ctx, key)
	return &StringCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// RPop 从列表右侧弹出元素
func (p *ClusterPipeliner) RPop(ctx context.Context, key string) *StringCmd {
	key = p.config.GetKeyWithPrefix(key)
	cmd := p.pipe.RPop(ctx, key)
	return &StringCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// 集合操作

// SAdd 向集合添加成员
func (p *ClusterPipeliner) SAdd(ctx context.Context, key string, members ...interface{}) *IntCmd {
	key = p.config.GetKeyWithPrefix(key)
	cmd := p.pipe.SAdd(ctx, key, members...)
	return &IntCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// SRem 从集合移除成员
func (p *ClusterPipeliner) SRem(ctx context.Context, key string, members ...interface{}) *IntCmd {
	key = p.config.GetKeyWithPrefix(key)
	cmd := p.pipe.SRem(ctx, key, members...)
	return &IntCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// SMembers 获取集合所有成员
func (p *ClusterPipeliner) SMembers(ctx context.Context, key string) *StringSliceCmd {
	key = p.config.GetKeyWithPrefix(key)
	cmd := p.pipe.SMembers(ctx, key)
	return &StringSliceCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// 有序集合操作

// ZAdd 向有序集合添加成员
func (p *ClusterPipeliner) ZAdd(ctx context.Context, key string, members ...ZMember) *IntCmd {
	key = p.config.GetKeyWithPrefix(key)
	redisMembers := make([]redis.Z, len(members))
	for i, member := range members {
		redisMembers[i] = redis.Z{
			Score:  member.Score,
			Member: member.Member,
		}
	}
	cmd := p.pipe.ZAdd(ctx, key, redisMembers...)
	return &IntCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// ZRem 从有序集合移除成员
func (p *ClusterPipeliner) ZRem(ctx context.Context, key string, members ...interface{}) *IntCmd {
	key = p.config.GetKeyWithPrefix(key)
	cmd := p.pipe.ZRem(ctx, key, members...)
	return &IntCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// ZRange 获取有序集合指定范围的成员（从小到大）
func (p *ClusterPipeliner) ZRange(ctx context.Context, key string, start, stop int64) *StringSliceCmd {
	key = p.config.GetKeyWithPrefix(key)
	cmd := p.pipe.ZRange(ctx, key, start, stop)
	return &StringSliceCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// 通用操作

// Del 删除键
func (p *ClusterPipeliner) Del(ctx context.Context, keys ...string) *IntCmd {
	prefixedKeys := make([]string, len(keys))
	for i, key := range keys {
		prefixedKeys[i] = p.config.GetKeyWithPrefix(key)
	}
	cmd := p.pipe.Del(ctx, prefixedKeys...)
	return &IntCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// Exists 检查键是否存在
func (p *ClusterPipeliner) Exists(ctx context.Context, keys ...string) *IntCmd {
	prefixedKeys := make([]string, len(keys))
	for i, key := range keys {
		prefixedKeys[i] = p.config.GetKeyWithPrefix(key)
	}
	cmd := p.pipe.Exists(ctx, prefixedKeys...)
	return &IntCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}

// Expire 设置键的过期时间
func (p *ClusterPipeliner) Expire(ctx context.Context, key string, expiration time.Duration) *BoolCmd {
	key = p.config.GetKeyWithPrefix(key)
	cmd := p.pipe.Expire(ctx, key, expiration)
	return &BoolCmd{
		val: cmd.Val(),
		err: cmd.Err(),
	}
}