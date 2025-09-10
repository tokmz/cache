package cache

import (
	"context"
	"time"
)

// Client Redis客户端统一接口
// 提供对Redis各种数据类型的操作方法
type Client interface {
	// 基础操作
	Close() error
	Ping(ctx context.Context) error

	// 字符串操作
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)
	GetSet(ctx context.Context, key string, value interface{}) (string, error)
	MGet(ctx context.Context, keys ...string) ([]interface{}, error)
	MSet(ctx context.Context, pairs ...interface{}) error
	Incr(ctx context.Context, key string) (int64, error)
	IncrBy(ctx context.Context, key string, value int64) (int64, error)
	Decr(ctx context.Context, key string) (int64, error)
	DecrBy(ctx context.Context, key string, value int64) (int64, error)

	// 哈希表操作
	HGet(ctx context.Context, key, field string) (string, error)
	HSet(ctx context.Context, key, field string, value interface{}) error
	HSetNX(ctx context.Context, key, field string, value interface{}) (bool, error)
	HDel(ctx context.Context, key string, fields ...string) (int64, error)
	HExists(ctx context.Context, key, field string) (bool, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HKeys(ctx context.Context, key string) ([]string, error)
	HVals(ctx context.Context, key string) ([]string, error)
	HLen(ctx context.Context, key string) (int64, error)
	HMGet(ctx context.Context, key string, fields ...string) ([]interface{}, error)
	HMSet(ctx context.Context, key string, pairs ...interface{}) error
	HIncrBy(ctx context.Context, key, field string, incr int64) (int64, error)

	// 列表操作
	LPush(ctx context.Context, key string, values ...interface{}) (int64, error)
	RPush(ctx context.Context, key string, values ...interface{}) (int64, error)
	LPop(ctx context.Context, key string) (string, error)
	RPop(ctx context.Context, key string) (string, error)
	LLen(ctx context.Context, key string) (int64, error)
	LRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	LIndex(ctx context.Context, key string, index int64) (string, error)
	LSet(ctx context.Context, key string, index int64, value interface{}) error
	LRem(ctx context.Context, key string, count int64, value interface{}) (int64, error)
	LTrim(ctx context.Context, key string, start, stop int64) error

	// 集合操作
	SAdd(ctx context.Context, key string, members ...interface{}) (int64, error)
	SRem(ctx context.Context, key string, members ...interface{}) (int64, error)
	SMembers(ctx context.Context, key string) ([]string, error)
	SIsMember(ctx context.Context, key string, member interface{}) (bool, error)
	SCard(ctx context.Context, key string) (int64, error)
	SPop(ctx context.Context, key string) (string, error)
	SRandMember(ctx context.Context, key string) (string, error)
	SInter(ctx context.Context, keys ...string) ([]string, error)
	SUnion(ctx context.Context, keys ...string) ([]string, error)
	SDiff(ctx context.Context, keys ...string) ([]string, error)

	// 有序集合操作
	ZAdd(ctx context.Context, key string, members ...ZMember) (int64, error)
	ZRem(ctx context.Context, key string, members ...interface{}) (int64, error)
	ZScore(ctx context.Context, key, member string) (float64, error)
	ZRank(ctx context.Context, key, member string) (int64, error)
	ZRevRank(ctx context.Context, key, member string) (int64, error)
	ZRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]ZMember, error)
	ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) ([]ZMember, error)
	ZRangeByScore(ctx context.Context, key string, min, max string) ([]string, error)
	ZRevRangeByScore(ctx context.Context, key string, max, min string) ([]string, error)
	ZCard(ctx context.Context, key string) (int64, error)
	ZCount(ctx context.Context, key, min, max string) (int64, error)
	ZIncrBy(ctx context.Context, key string, increment float64, member string) (float64, error)

	// 通用键操作
	Del(ctx context.Context, keys ...string) (int64, error)
	Exists(ctx context.Context, keys ...string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) (bool, error)
	ExpireAt(ctx context.Context, key string, tm time.Time) (bool, error)
	TTL(ctx context.Context, key string) (time.Duration, error)
	Type(ctx context.Context, key string) (string, error)
	Keys(ctx context.Context, pattern string) ([]string, error)
	Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error)

	// Lua脚本操作
	Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error)
	EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) (interface{}, error)
	ScriptExists(ctx context.Context, hashes ...string) ([]bool, error)
	ScriptFlush(ctx context.Context) error
	ScriptKill(ctx context.Context) error
	ScriptLoad(ctx context.Context, script string) (string, error)

	// 管道操作
	Pipeline() Pipeliner
	TxPipeline() Pipeliner
}

// Pipeliner 管道操作接口
type Pipeliner interface {
	// 执行管道中的所有命令
	Exec(ctx context.Context) ([]interface{}, error)
	// 丢弃管道中的所有命令
	Discard() error
	// 关闭管道
	Close() error

	// 字符串操作
	Get(ctx context.Context, key string) *StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *StatusCmd
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *BoolCmd
	Incr(ctx context.Context, key string) *IntCmd
	Decr(ctx context.Context, key string) *IntCmd

	// 哈希表操作
	HGet(ctx context.Context, key, field string) *StringCmd
	HSet(ctx context.Context, key, field string, value interface{}) *IntCmd
	HDel(ctx context.Context, key string, fields ...string) *IntCmd

	// 列表操作
	LPush(ctx context.Context, key string, values ...interface{}) *IntCmd
	RPush(ctx context.Context, key string, values ...interface{}) *IntCmd
	LPop(ctx context.Context, key string) *StringCmd
	RPop(ctx context.Context, key string) *StringCmd

	// 集合操作
	SAdd(ctx context.Context, key string, members ...interface{}) *IntCmd
	SRem(ctx context.Context, key string, members ...interface{}) *IntCmd
	SMembers(ctx context.Context, key string) *StringSliceCmd

	// 有序集合操作
	ZAdd(ctx context.Context, key string, members ...ZMember) *IntCmd
	ZRem(ctx context.Context, key string, members ...interface{}) *IntCmd
	ZRange(ctx context.Context, key string, start, stop int64) *StringSliceCmd

	// 通用操作
	Del(ctx context.Context, keys ...string) *IntCmd
	Exists(ctx context.Context, keys ...string) *IntCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *BoolCmd
}

// ZMember 有序集合成员
type ZMember struct {
	Score  float64
	Member interface{}
}

// 命令结果类型
type StringCmd struct {
	val string
	err error
}

func (cmd *StringCmd) Result() (string, error) {
	return cmd.val, cmd.err
}

func (cmd *StringCmd) Val() string {
	return cmd.val
}

func (cmd *StringCmd) Err() error {
	return cmd.err
}

type StatusCmd struct {
	val string
	err error
}

func (cmd *StatusCmd) Result() (string, error) {
	return cmd.val, cmd.err
}

func (cmd *StatusCmd) Val() string {
	return cmd.val
}

func (cmd *StatusCmd) Err() error {
	return cmd.err
}

type IntCmd struct {
	val int64
	err error
}

func (cmd *IntCmd) Result() (int64, error) {
	return cmd.val, cmd.err
}

func (cmd *IntCmd) Val() int64 {
	return cmd.val
}

func (cmd *IntCmd) Err() error {
	return cmd.err
}

type BoolCmd struct {
	val bool
	err error
}

func (cmd *BoolCmd) Result() (bool, error) {
	return cmd.val, cmd.err
}

func (cmd *BoolCmd) Val() bool {
	return cmd.val
}

func (cmd *BoolCmd) Err() error {
	return cmd.err
}

type StringSliceCmd struct {
	val []string
	err error
}

func (cmd *StringSliceCmd) Result() ([]string, error) {
	return cmd.val, cmd.err
}

func (cmd *StringSliceCmd) Val() []string {
	return cmd.val
}

func (cmd *StringSliceCmd) Err() error {
	return cmd.err
}