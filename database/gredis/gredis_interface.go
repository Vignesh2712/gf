package gredis

import (
	"context"
	"time"
)

type DB interface {
	Exists(ctx context.Context, keys ...string) (int64, error)
	Type(ctx context.Context, key string) (string, error)
	Rename(ctx context.Context, key, newKey string) (string, error)
	RenameNX(ctx context.Context, key, newKey string) (bool, error)
	Move(ctx context.Context, key, db string) (bool, error)
	Del(ctx context.Context, keys ...string) (int64, error)
	RandomKey(ctx context.Context) (string, error)
	DBSize(ctx context.Context) (int64, error)
	Keys(ctx context.Context, pattern string) ([]string, error)
	FlushDB(ctx context.Context) (string, error)
	FlushAll(ctx context.Context) (string, error)
	Select(ctx context.Context, index int64) (string, error)
}

type String interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error)
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)
	SetEx(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error)
	PSetEx(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error)
	Get(ctx context.Context, key string) (string, error)
	GetSet(ctx context.Context, key string, value interface{}) (string, error)
	StrLen(ctx context.Context, key string) (int64, error)
	Append(ctx context.Context, key string, value string) (int64, error)
	SetRange(ctx context.Context, key string, offset int64, value string) (int64, error)
	GetRange(ctx context.Context, key string, start, end int64) (string, error)
	Incr(ctx context.Context, key string) (int64, error)
	IncrBy(ctx context.Context, key string, value int64) (int64, error)
	IncrByFloat(ctx context.Context, key string, value float64) (float64, error)
	Decr(ctx context.Context, key string) (int64, error)
	DecrBy(ctx context.Context, key string, value int64) (int64, error)
	MSet(ctx context.Context, pairs ...interface{}) (string, error)
	MSetNX(ctx context.Context, pairs ...interface{}) (bool, error)
	MGet(ctx context.Context, keys ...string) ([]string, error)
}

type Hash interface {
	HSet(ctx context.Context, key, field, value string) (bool, error)
	HSetNX(ctx context.Context, key, field, value string) (bool, error)
	HGet(ctx context.Context, key, field string) (string, error)
	HExists(ctx context.Context, key, field string) (bool, error)
	HDel(ctx context.Context, key string, fields ...string) (int64, error)
	HLen(ctx context.Context, key string) (int64, error)
	HStrLen(ctx context.Context, key, field string) (int64, error)
	HIncrBy(ctx context.Context, key, field string, value int64) (int64, error)
	HIncrByFloat(ctx context.Context, key, field string, value float64) (float64, error)
	HMSet(ctx context.Context, key string, fields map[string]string) (string, error)
	HMGet(ctx context.Context, key string, fields ...string) ([]string, error)
	HKeys(ctx context.Context, key string) ([]string, error)
	HVal(ctx context.Context, key string) ([]string, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
}

type List interface {
	LPush(ctx context.Context, key string, values ...string) (int64, error)
	LPushX(ctx context.Context, key string, value string) (int64, error)
	RPush(ctx context.Context, key string, values ...string) (int64, error)
	RPushX(ctx context.Context, key string, value string) (int64, error)
	LPop(ctx context.Context, key string) (string, error)
	RPop(ctx context.Context, key string) (string, error)
	RPopLPush(ctx context.Context, source, destination string) (string, error)
	LRem(ctx context.Context, key string, count int64, value string) (int64, error)
	LLen(ctx context.Context, key string) (int64, error)
	LIndex(ctx context.Context, key string, index int64) (string, error)
	LInsert(ctx context.Context, key, op string, pivot, value string) (int64, error)
	LSet(ctx context.Context, key string, index int64, value string) (string, error)
	LRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	LTrim(ctx context.Context, key string, start, stop int64) (string, error)
	BLPop(ctx context.Context, keys ...string) ([]string, error)
	BRPop(ctx context.Context, keys ...string) ([]string, error)
	BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) (string, error)
}

type Set interface {
	SAdd(ctx context.Context, key string, members ...string) (int64, error)
	SIsMember(ctx context.Context, key string, member string) (bool, error)
	SPop(ctx context.Context, key string) (string, error)
	SRandMember(ctx context.Context, key string) (string, error)
	SRem(ctx context.Context, key string, members ...string) (int64, error)
	SMove(ctx context.Context, source, destination, member string) (bool, error)
	SCard(ctx context.Context, key string) (int64, error)
	SMembers(ctx context.Context, key string) ([]string, error)
	SInter(ctx context.Context, keys ...string) ([]string, error)
	SInterStore(ctx context.Context, destination string, keys ...string) (int64, error)
	SUnion(ctx context.Context, keys ...string) ([]string, error)
	SUnionStore(ctx context.Context, destination string, keys ...string) (int64, error)
	SDiff(ctx context.Context, keys ...string) ([]string, error)
	SDiffStore(ctx context.Context, destination string, keys ...string) (int64, error)
}

type SortedSet interface {
	ZAdd(ctx context.Context, key string, members ...interface{}) (int64, error)
	ZScore(ctx context.Context, key string, member string) (float64, error)
	ZIncrBy(ctx context.Context, key string, value float64, member string) (float64, error)
	ZCard(ctx context.Context, key string) (int64, error)
	ZCount(ctx context.Context, key string, min, max float64) (int64, error)
	ZRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	ZRangeByScore(ctx context.Context, key string, min, max string, withScores bool) ([]string, error)
	ZRevRangeByScore(ctx context.Context, key string, min, max string, withScores bool) ([]string, error)
	ZRank(ctx context.Context, key, member string) (int64, error)
	ZRevRank(ctx context.Context, key, member string) (int64, error)
	ZRem(ctx context.Context, key string, members ...string) (int64, error)
	ZRemRangeByRank(ctx context.Context, key string, start, stop int64) (int64, error)
	ZRemRangeByScore(ctx context.Context, key string, min, max string) (int64, error)
	ZRemRangeByLex(ctx context.Context, key string, min, max string) (int64, error)
	ZLexCount(ctx context.Context, key, min, max string) (int64, error)
	ZRangeByLex(ctx context.Context, key string, min, max string, offset, count int64) ([]string, error)
}

type HyperLogLog interface {
	PFAdd(ctx context.Context, key string, values ...string) (bool, error)
	PFCount(ctx context.Context, keys ...string) (int64, error)
	PFMerge(ctx context.Context, dest string, keys ...string) (string, error)
}

type Geo interface {
	GeoAdd(ctx context.Context, key string, members ...interface{}) (int64, error)
	GeoPos(ctx context.Context, key string, members ...string) ([]interface{}, error)
	GeoDist(ctx context.Context, key, member1, member2, unit string) (float64, error)
	GeoHash(ctx context.Context, key string, members ...string) ([]string, error)
	GeoRadius(ctx context.Context, key string, longitude, latitude float64, query interface{}) (interface{}, error)
	GeoRadiusByMember(ctx context.Context, key, member string, query interface{}) (interface{}, error)
}

type Bit interface {
	SetBit(ctx context.Context, key string, offset int64, value int) (int64, error)
	GetBit(ctx context.Context, key string, offset int64) (int64, error)
	BitCount(ctx context.Context, key string) (int64, error)
	BitPos(ctx context.Context, key string, bit int64) (int64, error)
	BitOp(ctx context.Context, op, dest string, keys ...string) (int64, error)
	BitField(ctx context.Context, key string, args ...interface{}) ([]int64, error)
}

type Expire interface {
	Expire(ctx context.Context, key string, seconds int64) (bool, error)
	ExpireAt(ctx context.Context, key string, timestamp int64) (bool, error)
	TTL(ctx context.Context, key string) (int64, error)
	PErsist(ctx context.Context, key string) (bool, error)
	PExpire(ctx context.Context, key string, milliseconds int64) (bool, error)
	PExpireAt(ctx context.Context, key string, timestamp int64) (bool, error)
	PTTL(ctx context.Context, key string) (int64, error)
}

type Lua interface {
	Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error)
	EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) (interface{}, error)
	ScriptLoad(ctx context.Context, script string) (string, error)
	ScriptExists(ctx context.Context, sha1s ...string) ([]bool, error)
	ScriptFlush(ctx context.Context) error
	ScriptKill(ctx context.Context) error
}

type PubSub interface {
	Publish(ctx context.Context, channel string, message interface{}) (int64, error)
	Subscribe(ctx context.Context, channels ...string) (PubSub, error)
	PSubscribe(ctx context.Context, channels ...string) (PubSub, error)
	UnSubscribe(ctx context.Context, channels ...string) (int64, error)
	PUnSubscribe(ctx context.Context, channels ...string) (int64, error)
	PubSubChannels(ctx context.Context) ([]string, error)
	PubSubNumSub(ctx context.Context, channels ...string) ([]int64, error)
	PubSubNumPat(ctx context.Context) (int64, error)
}

type Summary interface {
	DB
	String
	Hash
	List
	Set
	SortedSet
	HyperLogLog
	Geo
	Bit
	Expire
	Lua
	PubSub
}
