package gredis

import (
	"errors"
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/util/gconv"
)

type GeoLocation struct {
	Name                      string
	Longitude, Latitude string
	GeoHash                   int64
	Dist string
}

func typeInt64(i interface{}, err error) (int64, error) {
	if err != nil {
		return 0, err
	}
	return gconv.Int64(i), nil
}



func typeInt(i interface{}, err error) (int, error) {
	if err != nil {
		return 0, err
	}
	return gconv.Int(i), nil

}
func typeFloat64(i interface{}, err error) (float64, error) {
	if err != nil {
		return 0, err
	}
	return gconv.Float64(i), nil
}

func typeString(i interface{}, err error) (string, error) {
	if err != nil {
		return "", err
	}
	return gconv.String(i), nil
}

func typeStrings(i interface{}, err error) ([]string, error) {
	if err != nil {
		return nil, err
	}
	return gconv.Strings(i), nil
}
func typeStringss(i interface{}, err error) ([][]string, error) {
	if err != nil {
		return nil, err
	}
	ss:=[][]string{}
	is:=gconv.Interfaces(i)
	for _,v:=range is{
		//fmt.Println(gconv.Strings(v))
		ss=append(ss,gconv.Strings(v))
	}

	return ss, nil
}

func typeGeoLocation(i interface{}, err error) ([]GeoLocation, error){
	if err != nil {
		return nil, err
	}
	var loc GeoLocation
	ss:=[]GeoLocation{}
	is:=gconv.Interfaces(i)
	for _,v:=range is{
		s1:=gconv.Strings(v)
		loc.Longitude=s1[0]
		loc.Latitude=s1[1]
		ss=append(ss,loc)
	}

	return ss, nil
}

func typeGeoLocationd(i interface{}, err error) ([]GeoLocation, error){
	if err != nil {
		return nil, err
	}

	var loc GeoLocation
	ss:=[]GeoLocation{}
	is:=gconv.Interfaces(i)
	for _,v:=range is{
		s1:=gconv.Interfaces(v)
		if len(s1)==3{
			loc.Name=gconv.String(s1[0])
			loc.Dist=gconv.String(s1[1])
			s1_3:=gconv.Strings(s1[2])
			loc.Longitude=s1_3[0]
			loc.Latitude=s1_3[1]

		}else{
			loc.Name=gconv.String(s1[0])
			if s1_2,ok:=s1[1].(string);ok==true{
				loc.Dist=s1_2
			}else{
				s1_2s:=gconv.Strings(s1[1])
				loc.Longitude=s1_2s[0]
				loc.Latitude=s1_2s[1]
			}
		}

		ss=append(ss,loc)
	}

	return ss, nil
}




func typeBool(i interface{}, err error) (bool, error) {
	if err != nil {
		return false, err
	}
	return gconv.Bool(i), nil
}

func typeInterfacess(i interface{}, err error) ([]interface{}, error) {
	if err != nil {
		return nil, err
	}
	return gconv.Interfaces(i), nil
}

//==========================================================================key
func (c *Redis) Del(key ...string) (int, error) {
	return typeInt(c.commnddo("DEL", gconv.Interfaces(key)...))
}

func (c *Redis) Exists(key string) (int, error) {
	return typeInt( c.commnddo("EXISTS", key))
}

func (c *Redis) Ttl(key string) (int64, error) {
	return typeInt64(c.commnddo("TTL", key))
}

func (c *Redis) Expire(key string, time int64) (int64, error) {
	return typeInt64(c.commnddo("EXPIRE", key,time))
}

func (c *Redis) Dump(key string) (string, error) {
	return typeString(c.commnddo("DUMP", key))
}

func (c *Redis) Expireat(key string, timestamp int64) (int, error) {
	return typeInt(c.commnddo("EXPIREAT", key, timestamp))
}

// Returns all keys matching pattern, but not for clustering
func (c *Redis) Keys(key string) ([]interface{}, error) {
	return typeInterfacess( c.commnddo("KEYS", key))
}

func (c *Redis) Object(action, key string) (interface{}, error) {
	return c.commnddo("OBJECT", action, key)
}

func (c *Redis) Persist(key string) (int, error) {
	return typeInt(c.commnddo("PERSIST", key))
}
func (c *Redis) Pttl(key string) (int64, error) {
	return typeInt64(c.commnddo("PTTL", key))
}
func (c *Redis) RandomKey() (interface{}, error) {
	return c.commnddo("RANDOMKEY")
}

func (c *Redis) Rename(oldkey, newkey string) (string, error) {
	return typeString(c.commnddo("RENAME", oldkey, newkey))
}

func (c *Redis) Renamenx(oldkey, newkey string) (int, error) {
	return typeInt(c.commnddo("RENAMENX", oldkey, newkey))
}

func (c *Redis) ReStore(key string, ttl int64, serializedvalue string,replace ...string) (string, error) {
	str1:=""
	if len(replace)>0{
		str1=replace[0]
	}
	return typeString(c.commnddo("RESTORE", key, ttl, serializedvalue,str1))
}

func (c *Redis) Sort(key string, params ...interface{}) ([]interface{}, error) {
	return typeInterfacess(c.commnddo("SORT", append([]interface{}{key},params...)...))
}

func (c *Redis) Type(key string) (string, error) {
	return  typeString(c.commnddo("type", key))
}

//============================================================================string
func (c *Redis) Append(key, value string) (int64, error) {
	return typeInt64(c.commnddo("append", key, value))
}

func (c *Redis) Set(key, value string) (interface{}, error) {
	return c.commnddo("set", key, value)
}

func (c *Redis) Get(key string) (string, error) {
	return  typeString(c.commnddo("get", key))
}

func (c *Redis) BitCount(key string) (int, error) {
	return typeInt( c.commnddo("BITCOUNT", key))
}

func (c *Redis) BiTop(params ...string) (int, error) {
	return typeInt( c.commnddo("BITOP", gconv.Interfaces(params)...))
}

func (c *Redis) BitPos(key string, bit int, option ...int) (int, error) {
	return typeInt(c.commnddo("BITPOS", append([]interface{}{key,bit},gconv.Interfaces(option)...)...))
}

func (c *Redis) BitField(option string) ([]interface{}, error) {
	return typeInterfacess(c.commnddo("BITFIELD", option))
}

func (c *Redis) Decr(key string) (int64, error) {
	return typeInt64(c.commnddo("DECR", key))
}

func (c *Redis) Decrby(key string, decrement int64) (int64, error) {
	return typeInt64(c.commnddo("DECRBY", key, decrement))
}

func (c *Redis) Getbit(key string, offset int) (int, error) {
	return typeInt(c.commnddo("GETBIT", key, offset))
}

func (c *Redis) GetRange(key string, start, end int) (string, error) {
	return typeString(c.commnddo("GETRANGE", key, start, end))
}

func (c *Redis) GetSet(key string, value string) (string, error) {
	return typeString(c.commnddo("GETSET", key, value))
}

func (c *Redis) Incr(key string) (int64, error) {
	return typeInt64(c.commnddo("INCR", key))
}

func (c *Redis) IncrBy(key string, increment int64) (int64, error) {
	return typeInt64(c.commnddo("INCRBY", key, increment))
}

func (c *Redis) IncrByFloat(key string, increment float64) (string, error) {
	return typeString(c.commnddo("INCRBYFLOAT", key, increment))
}

func (c *Redis) Mget(key ...string) ([]string, error) {
	if len(key)<1{
		return nil,errors.New("there must be one key's name")
	}
	return typeStrings(c.commnddo("MGET", gconv.Interfaces(key)...))
}

func (c *Redis) Mset(params ...string) (string, error) {
	if len(params)<2{
		return "",errors.New("there must be one k-v ")
	}
	return typeString(c.commnddo("MSET", gconv.Interfaces(params)...))
}

func (c *Redis) Msetnx(params ...string) (int, error) {

	return typeInt(c.commnddo("MSETNX", gconv.Interfaces(params)...))
}

func (c *Redis) Psetex(key string, milliseconds int64, value string) (string, error) {
	return typeString(c.commnddo("PSETEX", key, milliseconds, value))
}

func (c *Redis) Setbit(key string, offset,value int   ) (int, error) {
	return typeInt(c.commnddo("SETBIT", key, offset, value))
}

func (c *Redis) Setex(key string, seconds int64, value string) (string, error) {
	return typeString(c.commnddo("SETEX", key, seconds, value))
}

func (c *Redis) Setnx(key string, value string) (int, error) {
	return typeInt(c.commnddo("SETNX", key, value))
}

func (c *Redis) SetRange(key string, offset int, value string) (int, error) {
	return typeInt( c.commnddo("SETRANGE", key, offset, value))
}

func (c *Redis) Strlen(key string) (int, error) {
	return typeInt( c.commnddo("STRLEN", key))
}

//=======================================================================Hash
func (c *Redis) Hset(key, fieldname string, value interface{}) (int, error) {
	return typeInt(c.commnddo("HSET", key, fieldname, value))
}

func (c *Redis) Hsetnx(key, fieldname string, value interface{}) (int, error) {
	return typeInt(c.commnddo("HSETNX", key, fieldname, value))
}

func (c *Redis) Hget(key, fieldname string) (string, error) {
	return typeString(c.commnddo("HGET", key, fieldname))
}

func (c *Redis) Hexists(key, fieldname string) (int, error) {
	return typeInt(c.commnddo("HEXISTS", key, fieldname))
}

func (c *Redis) Hdel(key string, fields ...string) (int, error) {
	 if len(fields)<1{
	 	return 0,errors.New("must have one field's name")
	 }
	return typeInt(c.commnddo("HDEL", gconv.Interfaces(append([]string{key},fields...))...))
}

func (c *Redis) Hlen(key string) (int, error) {
	return typeInt(c.commnddo("HLEN", key))
}

func (c *Redis) Hstrlen(key, field string) (int, error) {
	return typeInt(c.commnddo("HSTRLEN", key, field))

}

func (c *Redis) HincrBy(key, field string, increment int64) (int64, error) {
	return typeInt64(c.commnddo("HINCRBY", key, field, increment))
}

func (c *Redis) HincrByFloat(key, field string, increment float64) (string, error) {
	return typeString(c.commnddo("HINCRBYFLOAT", key, field, increment))
}

func (c *Redis) Hmset(key string, params ...interface{}) (string, error) {

	return typeString(c.commnddo("HMSET", append([]interface{}{key},params...)...))
}

func (c *Redis) Hmget(key string,option ...string) ([]string, error) {
	return typeStrings( c.commnddo("HMGET", gconv.Interfaces(append([]string{key},option...))...))
}

func (c *Redis) Hkeys(key string) ([]string, error) {
	return typeStrings(c.commnddo("HKEYS", key))
}

func (c *Redis) Hvals(key string) ([]string, error) {
	return typeStrings(c.commnddo("HVALS", key))
}

func (c *Redis) HgetAll(key string) ([]string, error) {
	return typeStrings(c.commnddo("HGETALL", key))
}

//==============================================================================list
func (c *Redis) Lpush(key string, values ...interface{}) (int64, error) {
	return typeInt64(c.commnddo("LPUSH", append([]interface{}{key},values...)...))

}

func (c *Redis) Lpushx(key string, values interface{}) (int64, error) {
	return typeInt64(c.commnddo("LPUSHX", key, values))
}

func (c *Redis) Rpush(key string, values ...interface{}) (int64, error) {

	return typeInt64(c.commnddo("RPUSH", append([]interface{}{key},values...)...))
}

func (c *Redis) Rpushx(key string, values interface{}) (int64, error) {
	return typeInt64(c.commnddo("RPUSHX", key, values))
}

func (c *Redis) Lpop(key string) (string, error) {
	return typeString(c.commnddo("LPOP", key))
}

func (c *Redis) Rpop(key string) (string, error) {
	return typeString(c.commnddo("RPOP", key))
}

func (c *Redis) RpoplPush(source, destination string) (string, error) {
	return typeString(c.commnddo("RPOPLPUSH",   source, destination))
}

func (c *Redis) Lrem(key string, count int, value interface{}) (int64, error) {
	return typeInt64(c.commnddo("LREM", key, count, value))
}

func (c *Redis) Llen(key string) (int64, error) {
	return typeInt64(c.commnddo("LLEN", key))
}

func (c *Redis) Lindex(key string, index int64) (string, error) {
	return typeString(c.commnddo("LINDEX", key, index))
}

func (c *Redis) Linsert(key, layout, pivot string, value interface{}) (int64, error) {
	return typeInt64(c.commnddo("LINSERT", key, layout, pivot, value))
}

func (c *Redis) Lset(key  string, index int64, value interface{}) (string, error) {
	return typeString(c.commnddo("LSET", key, index, value))
}

func (c *Redis) Lrange(key string, start, stop int64) ([]string, error) {
	return typeStrings(c.commnddo("LRANGE", key, start, stop))
}

func (c *Redis) BlPop(key string, params ...interface{}) ([]string, error) {
	return typeStrings(c.commnddo("BLPOP", append([]interface{}{key},params...)...))
}

func (c *Redis) BrPop(key string, params ...interface{}) ([]string, error) {
	return typeStrings(c.commnddo("BRPOP", append([]interface{}{key},params...)...))
}

func (c *Redis) BrPoplPush(  source, destination string, timeout int) ([]string, error) {
	return typeStrings(c.commnddo("BRPOPLPUSH", source, destination, timeout))
}

//========================================================================================set
func (c *Redis) Sadd(key string, members ...interface{}) (int64, error) {

	return typeInt64(c.commnddo("SADD", append([]interface{}{key},members...)...))
}

func (c *Redis) SisMember(key, member string) (int, error) {
	return typeInt(c.commnddo("SISMEMBER", key, member))
}

func (c *Redis) Spop(key string) (string, error) {
	return typeString(c.commnddo("SPOP", key))
}

func (c *Redis) SrandMember(key string, count ...int) ([]string, error) {
	if len(count)==0{
		return  typeStrings(c.commnddo("SRANDMEMBER",key,1))
	}
	return  typeStrings(c.commnddo("SRANDMEMBER", key,count[0]))
}

func (c *Redis) Srem(key string,members ...string) (int, error) {
	return typeInt(c.commnddo("SREM", append([]interface{}{key},gconv.Interfaces(members)...)...))
}

func (c *Redis) Smove(source, destination, member string) (int, error) {
	return typeInt(c.commnddo("SMOVE", source, destination, member))
}

func (c *Redis) Scard(key string) (int64, error) {
	return typeInt64(c.commnddo("SCARD", key))
}

func (c *Redis) Smembers(key string) ([]string, error) {
	return typeStrings(c.commnddo("SMEMBERS", key))
}

func (c *Redis) Sinter(keys ...string) ([]string, error) {
	if len(keys)==0{
		return nil,errors.New("must have a key")
	}
	return typeStrings(c.commnddo("SINTER", gconv.Interfaces(keys)...))
}

func (c *Redis)  SinterStore(destination string, key string, keys ...string) (int64, error) {
	return typeInt64(c.commnddo("SINTERSTORE",append([]interface{}{destination,key},gconv.Interfaces(keys)...)...))
}

func (c *Redis) Sunion(key string, keys ...string) ([]string, error) {
	return typeStrings(c.commnddo("SUNION", append([]interface{}{key},gconv.Interfaces(keys)...)...))
}

func (c *Redis) SunionStore(destination string, key string, keys ...string) (int64, error) {
	return typeInt64(c.commnddo("SUNIONSTORE", append([]interface{}{destination,key},gconv.Interfaces(keys)...)...))
}

func (c *Redis) Sdiff(key string, keys ...string) ([]string, error) {
	return typeStrings(c.commnddo("SDIFF", append([]interface{}{key},gconv.Interfaces(keys)...)...))
}

func (c *Redis) SdiffStore(destination string, key string, keys ...string) (int64, error) {

	return typeInt64(c.commnddo("SDIFFSTORE", append([]interface{}{destination,key},gconv.Interfaces(keys)...)...))
}

//======================================================================================zset

func (c *Redis) Zadd(params ...interface{}) (int, error) {
	return typeInt(c.commnddo("ZADD", params...))
}

func (c *Redis) Zscore(key string, member interface{}) (string, error) {
	return typeString(c.commnddo("ZSCORE", key, member))
}

func (c *Redis) ZinCrby(key string, increment float64, member interface{}) (string, error) {
	return typeString(c.commnddo("ZINCRBY", key, increment, member))
}

func (c *Redis) Zcard(key string) (int64, error) {
	return typeInt64(c.commnddo("ZCARD", key))
}

func (c *Redis) Zcount(key string, min, max int64) (int64, error) {
	return typeInt64(c.commnddo("ZCOUNT",key, min, max))
}

func (c *Redis) Zrange(key string, start, stop int64, param ...string) ([]string, error) {
	if len(param)==0{
		return typeStrings( c.commnddo("ZRANGE", key,start, stop))
	}
	return typeStrings( c.commnddo("ZRANGE",key, start, stop, param[0]))
}

func (c *Redis) ZrevRange(key string, start, stop int64, param ...string) ([]string, error) {
	if len(param)==0{
		return typeStrings( c.commnddo("ZRANGE", key,start, stop))
	}
	return typeStrings(c.commnddo("ZREVRANGE",key, start, stop, param[0]))
}

func (c *Redis) ZrangeByScore(key , min, max string, options ...interface{}) ([]string, error) {

	return typeStrings(c.commnddo("ZRANGEBYSCORE",append([]interface{}{key,min,max},options...)...))
}

func (c *Redis) ZrevRangeByScore(key string, min, max string, options ...interface{}) ([]string, error) {
	return typeStrings(c.commnddo("ZREVRANGEBYSCORE",append([]interface{}{key,min,max},options...)...))
}

func (c *Redis) Zrank(key, member string) (int64, error) {
	return typeInt64(c.commnddo("ZRANK",key, member))
}

func (c *Redis) ZrevRank(key, member string) (int64, error) {
	return typeInt64(c.commnddo("ZREVRANK",key, member))
}

func (c *Redis) Zrem(key string, member ...interface{}) (int, error) {
	 if len(member)==0{
	 	return 0,errors.New("must have an one key")
	 }
	return typeInt(c.commnddo("ZREM",append([]interface{}{key},member...)...))
}

func (c *Redis) ZremRangeByRank(key string, start, stop int64) (int64, error) {
	return typeInt64(c.commnddo("ZREMRANGEBYRANK", key, start, stop))
}

func (c *Redis) ZremRangeByScore(key string, min, max float64) (int64, error) {
	return typeInt64(c.commnddo("ZREMRANGEBYSCORE", key, min, max))
}

func (c *Redis) ZrangeByLex(key, min, max string, options ...interface{}) ([]string, error) {
	return typeStrings(c.commnddo("ZRANGEBYLEX", append([]interface{}{key,min,max},options...)...))
}

func (c *Redis) ZlexCount(key, min, max string) (int64, error) {
	return typeInt64(c.commnddo("ZLEXCOUNT", key, min, max))
}

func (c *Redis) ZremRangeByLex(key, min, max string) (int64, error) {
	return typeInt64(c.commnddo("ZREMRANGEBYLEX", key, min, max))
}

func (c *Redis) ZunionStore(options ...interface{}) (int64, error) {
	if len(options)<3{
		return 0,errors.New("there must be three parameters")
	}
	return typeInt64(c.commnddo("ZUNIONSTORE", options...))
}
func (c *Redis) ZinterStore(options ...interface{}) (int64, error) {
	return typeInt64(c.commnddo("ZINTERSTORE", options...))
}

//================================================================HyperLogLog
func (c *Redis) PfAdd(key string, options ...interface{}) (int, error) {
	return typeInt(c.commnddo("PFADD", append([]interface{}{key},options...)...))
}

func (c *Redis) PfCount(keys ...string) (int64, error) {
	return typeInt64(c.commnddo("PFCOUNT", gconv.Interfaces(keys)...))
}

func (c *Redis) PfMerge(keys ...string) (string, error) {
	if len(keys)<2{
		return "",errors.New("need at least two keys")
	}
	return typeString(c.commnddo("PFMERGE", gconv.Interfaces(keys)...))
}

//================================================================================GEO
func (c *Redis) GeoAdd(key string, params ...interface{}) (int, error) {
	return typeInt(c.commnddo("GEOADD",append([]interface{}{key},params...)...))
}

func (c *Redis) GeoPos(key string, member ...interface{}) ([]GeoLocation, error) {
	return typeGeoLocation(c.commnddo("GEOPOS", append([]interface{}{key},member...)...))
}

func (c *Redis) GeoDist(key string, params ...string) (string, error) {
	return typeString(c.commnddo("GEODIST",  append([]interface{}{key},gconv.Interfaces(params)...)...) )
}

func (c *Redis) GeoRadius(key string, member ...interface{}) ([]GeoLocation, error) {
	if len(member)<5{
		return nil,errors.New("there are must have five keys")
	}
	return typeGeoLocationd(c.commnddo("GEORADIUS", append([]interface{}{key},member...)...))
}

func (c *Redis) GeoRadiusByMember(key string, member ...interface{}) ([]interface{}, error) {
	param := garray.NewArrayFrom(member).InsertBefore(0, key)
	return typeInterfacess(c.commnddo("GEORADIUSBYMEMBER", gconv.Interfaces(param)...))
}

func (c *Redis) GeoHash(key string, member ...interface{}) ([]interface{}, error) {
	param := garray.NewArrayFrom(member).InsertBefore(0, key)
	return typeInterfacess(c.commnddo("GEOHASH", gconv.Interfaces(param)...))
}

//============================================================================channel
func (c *Redis) PubList(channel, message string) (int, error) {
	return typeInt(c.commnddo("PUBLISH", channel, message))
}

func (c *Redis) SubScribe(channel ...string) (interface{}, error) {
	return c.commnddo("SUBSCRIBE", gconv.Interfaces(channel)...)
}

func (c *Redis) PsubScribe(pattern ...string) (interface{}, error) {
	return c.commnddo("PSUBSCRIBE", gconv.Interfaces(pattern)...)
}

func (c *Redis) UnSubScribe(pattern ...string) (interface{}, error) {
	return c.commnddo("UNSUBSCRIBE", gconv.Interfaces(pattern)...)
}

func (c *Redis) PubSubScribe(pattern ...string) (interface{}, error) {
	return c.commnddo("PUNSUBSCRIBE", gconv.Interfaces(pattern)...)
}
