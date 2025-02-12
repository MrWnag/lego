package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/liwei1dao/lego/sys/redis/cluster"
	"github.com/liwei1dao/lego/sys/redis/single"
	"google.golang.org/protobuf/proto"
)

func newSys(options Options) (sys *Redis, err error) {
	sys = &Redis{options: options}
	err = sys.init()
	return
}

type Redis struct {
	options Options
	client  IRedis
}

func (this *Redis) init() (err error) {
	if this.options.RedisType == Redis_Single {
		this.client, err = single.NewSys(
			this.options.Redis_Single_Addr,
			this.options.Redis_Single_Password,
			this.options.Redis_Single_DB,
			this.options.Redis_Single_PoolSize,
			this.options.TimeOut,
			this.Encode,
			this.Decode,
		)
	} else if this.options.RedisType == Redis_Cluster {
		this.client, err = cluster.NewSys(
			this.options.Redis_Cluster_Addr,
			this.options.Redis_Cluster_Password,
			this.options.TimeOut,
			this.Encode,
			this.Decode,
		)
	} else {
		err = fmt.Errorf("init Redis err:RedisType - %d", this.options.RedisType)
	}
	return
}

func (this *Redis) Close() (err error) {
	return this.client.Close()
}
func (this *Redis) Do(ctx context.Context, args ...interface{}) *redis.Cmd {
	return this.client.Do(ctx, args...)
}
func (this *Redis) Pipeline(ctx context.Context, fn func(pipe redis.Pipeliner) error) (err error) {
	return this.client.Pipeline(ctx, fn)
}
func (this *Redis) TxPipelined(ctx context.Context, fn func(pipe redis.Pipeliner) error) (err error) {
	return this.client.TxPipelined(ctx, fn)
}
func (this *Redis) Watch(ctx context.Context, fn func(*redis.Tx) error, keys ...string) (err error) {
	return this.client.Watch(ctx, fn)
}
func (this *Redis) Lock(key string, outTime int) (result bool, err error) {
	return this.client.Lock(key, outTime)
}
func (this *Redis) UnLock(key string) (err error) {
	return this.client.UnLock(key)
}

///数据编码
func (this *Redis) Encode(value interface{}) (result []byte, err error) {
	if this.options.RedisStorageType == JsonData {
		result, err = json.Marshal(value)
	} else {
		if _, ok := value.(proto.Message); ok {
			result, err = proto.Marshal(value.(proto.Message))
		} else {
			result, err = json.Marshal(value)
		}
	}
	return
}

func (this *Redis) Decode(value []byte, result interface{}) (err error) {
	if this.options.RedisStorageType == JsonData {
		err = json.Unmarshal(value, result)
	} else {
		if _, ok := result.(proto.Message); ok {
			err = proto.Unmarshal(value, result.(proto.Message))
		} else {
			err = json.Unmarshal(value, result)
		}
	}
	return
}

func (this *Redis) Delete(key string) (err error) {
	return this.client.Delete(key)
}

func (this *Redis) ExistsKey(key string) (iskeep bool, err error) {
	return this.client.ExistsKey(key)
}

func (this *Redis) ExpireKey(key string, expire int) (err error) {
	return this.client.ExpireKey(key, expire)
}
func (this *Redis) ExpireatKey(key string, expire_unix int64) (err error) {
	return this.client.ExpireatKey(key, expire_unix)
}
func (this *Redis) Pexpirekey(key string, expire int) (err error) {
	return this.client.Pexpirekey(key, expire)
}
func (this *Redis) PexpireatKey(key string, expire_unix int64) (err error) {
	return this.client.PexpireatKey(key, expire_unix)
}
func (this *Redis) PersistKey(key string) (err error) {
	return this.client.PersistKey(key)
}
func (this *Redis) PttlKey(key string) (leftexpire int64, err error) {
	return this.client.PttlKey(key)
}
func (this *Redis) TtlKey(key string) (leftexpire int64, err error) {
	return this.client.TtlKey(key)
}
func (this *Redis) RenameKye(oldkey string, newkey string) (err error) {
	return this.client.RenameKye(oldkey, newkey)
}
func (this *Redis) RenamenxKey(oldkey string, newkey string) (err error) {
	return this.client.RenamenxKey(oldkey, newkey)
}
func (this *Redis) Keys(pattern string) (keys []string, err error) {
	return this.client.Keys(pattern)
}

///获取键类型
func (this *Redis) Type(key string) (ty string, err error) {
	return this.client.Type(key)
}

/*String*/
func (this *Redis) Set(key string, value interface{}, expiration time.Duration) (err error) {
	return this.client.Set(key, value, expiration)
}
func (this *Redis) SetNX(key string, value interface{}) (result int64, err error) {
	return this.client.SetNX(key, value)
}
func (this *Redis) MSet(keyvalues map[string]interface{}) (err error) {
	return this.client.MSet(keyvalues)
}
func (this *Redis) MSetNX(keyvalues map[string]interface{}) (err error) {
	return this.client.MSetNX(keyvalues)
}
func (this *Redis) Incr(key string) (err error) {
	return this.client.Incr(key)
}
func (this *Redis) IncrBY(key string, value int) (err error) {
	return this.client.IncrBY(key, value)
}
func (this *Redis) Incrbyfloat(key string, value float32) (err error) {
	return this.client.Incrbyfloat(key, value)
}
func (this *Redis) Decr(key string, value int) (err error) {
	return this.client.Decr(key, value)
}
func (this *Redis) DecrBy(key string, value int) (err error) {
	return this.client.DecrBy(key, value)
}
func (this *Redis) Append(key string, value interface{}) (err error) {
	return this.client.Append(key, value)
}
func (this *Redis) Get(key string, value interface{}) (err error) {
	return this.client.Get(key, value)
}
func (this *Redis) GetSet(key string, value interface{}, result interface{}) (err error) {
	return this.client.GetSet(key, value, result)
}
func (this *Redis) MGet(keys ...string) (result []string, err error) {
	return this.client.MGet(keys...)
}
func (this *Redis) INCRBY(key string, amount int64) (result int64, err error) {
	return this.client.INCRBY(key, amount)
}

/*List*/
func (this *Redis) Lindex(key string, value interface{}) (err error) {
	return this.client.Lindex(key, value)
}
func (this *Redis) Linsert(key string, isbefore bool, tager interface{}, value interface{}) (err error) {
	return this.client.Linsert(key, isbefore, tager, value)
}
func (this *Redis) Llen(key string) (result int, err error) {
	return this.client.Llen(key)
}
func (this *Redis) LPop(key string, value interface{}) (err error) {
	return this.client.LPop(key, value)
}
func (this *Redis) LPush(key string, values ...interface{}) (err error) {
	return this.client.LPush(key, values...)
}
func (this *Redis) LPushX(key string, values ...interface{}) (err error) {
	return this.client.LPushX(key, values...)
}
func (this *Redis) LRange(key string, start, end int, valuetype reflect.Type) (result []interface{}, err error) {
	return this.client.LRange(key, start, end, valuetype)
}
func (this *Redis) LRem(key string, count int, target interface{}) (err error) {
	return this.client.LRem(key, count, target)
}
func (this *Redis) LSet(key string, index int, value interface{}) (err error) {
	return this.client.LSet(key, index, value)
}
func (this *Redis) Ltrim(key string, start, stop int) (err error) {
	return this.client.Ltrim(key, start, stop)
}
func (this *Redis) Rpop(key string, value interface{}) (err error) {
	return this.client.Rpop(key, value)
}
func (this *Redis) RPopLPush(oldkey string, newkey string, value interface{}) (err error) {
	return this.client.RPopLPush(oldkey, newkey, value)
}
func (this *Redis) RPush(key string, values ...interface{}) (err error) {
	return this.client.RPush(key, values...)
}
func (this *Redis) RPushX(key string, values ...interface{}) (err error) {
	return this.client.RPushX(key, values...)
}

/*Hash*/
func (this *Redis) HDel(key string, fields ...string) (err error) {
	return this.client.HDel(key, fields...)
}
func (this *Redis) HExists(key string, field string) (result bool, err error) {
	return this.client.HExists(key, field)
}
func (this *Redis) HGet(key string, field string, value interface{}) (err error) {
	return this.client.HGet(key, field, value)
}
func (this *Redis) HGetAll(key string, valuetype reflect.Type) (result []interface{}, err error) {
	return this.client.HGetAll(key, valuetype)
}
func (this *Redis) HIncrBy(key string, field string, value int) (err error) {
	return this.client.HIncrBy(key, field, value)
}
func (this *Redis) HIncrByFloat(key string, field string, value float32) (err error) {
	return this.client.HIncrByFloat(key, field, value)
}
func (this *Redis) Hkeys(key string) (result []string, err error) {
	return this.client.Hkeys(key)
}
func (this *Redis) Hlen(key string) (result int, err error) {
	return this.client.Hlen(key)
}
func (this *Redis) HMGet(key string, valuetype reflect.Type, fields ...string) (result []interface{}, err error) {
	return this.client.HMGet(key, valuetype, fields...)
}
func (this *Redis) HMSet(key string, value map[string]interface{}) (err error) {
	return this.client.HMSet(key, value)
}
func (this *Redis) HSet(key string, field string, value interface{}) (err error) {
	return this.client.HSet(key, field, value)
}
func (this *Redis) HSetNX(key string, field string, value interface{}) (err error) {
	return this.client.HSetNX(key, field, value)
}

/*Set*/
func (this *Redis) SAdd(key string, values ...interface{}) (err error) {
	return this.client.SAdd(key, values...)
}
func (this *Redis) SCard(key string) (result int64, err error) {
	return this.client.SCard(key)
}
func (this *Redis) SDiff(valuetype reflect.Type, keys ...string) (result []interface{}, err error) {
	return this.client.SDiff(valuetype, keys...)
}
func (this *Redis) SDiffStore(destination string, keys ...string) (result int64, err error) {
	return this.client.SDiffStore(destination, keys...)
}
func (this *Redis) SInter(valuetype reflect.Type, keys ...string) (result []interface{}, err error) {
	return this.client.SInter(valuetype, keys...)
}
func (this *Redis) SInterStore(destination string, keys ...string) (result int64, err error) {
	return this.client.SInterStore(destination, keys...)
}
func (this *Redis) Sismember(key string, value interface{}) (iskeep bool, err error) {
	return this.client.Sismember(key, value)
}
func (this *Redis) SMembers(valuetype reflect.Type, key string) (result []interface{}, err error) {
	return this.client.SMembers(valuetype, key)
}
func (this *Redis) SMove(source string, destination string, member interface{}) (result bool, err error) {
	return this.client.SMove(source, destination, member)
}
func (this *Redis) Spop(key string) (result string, err error) {
	return this.client.Spop(key)
}
func (this *Redis) Srandmember(key string) (result string, err error) {
	return this.client.Srandmember(key)
}
func (this *Redis) SRem(key string, members ...interface{}) (result int64, err error) {
	return this.client.SRem(key, members...)
}
func (this *Redis) SUnion(valuetype reflect.Type, keys ...string) (result []interface{}, err error) {
	return this.client.SUnion(valuetype, keys...)
}
func (this *Redis) Sunionstore(destination string, keys ...string) (result int64, err error) {
	return this.client.Sunionstore(destination, keys...)
}
func (this *Redis) Sscan(key string, _cursor uint64, match string, count int64) (keys []string, cursor uint64, err error) {
	return this.client.Sscan(key, _cursor, match, count)
}

/*ZSet*/
func (this *Redis) ZAdd(key string, members ...*redis.Z) (err error) {
	return this.client.ZAdd(key, members...)
}
func (this *Redis) ZCard(key string) (result int64, err error) {
	return this.client.ZCard(key)
}
func (this *Redis) ZCount(key string, min string, max string) (result int64, err error) {
	return this.client.ZCount(key, min, max)
}
func (this *Redis) ZIncrBy(key string, increment float64, member string) (result float64, err error) {
	return this.client.ZIncrBy(key, increment, member)
}
func (this *Redis) ZInterStore(destination string, store *redis.ZStore) (result int64, err error) {
	return this.client.ZInterStore(destination, store)
}
func (this *Redis) ZLexCount(key string, min string, max string) (result int64, err error) {
	return this.client.ZLexCount(key, min, max)
}
func (this *Redis) ZRange(valuetype reflect.Type, key string, start int64, stop int64) (result []interface{}, err error) {
	return this.client.ZRange(valuetype, key, start, stop)
}
func (this *Redis) ZRangeByLex(valuetype reflect.Type, key string, opt *redis.ZRangeBy) (result []interface{}, err error) {
	return this.client.ZRangeByLex(valuetype, key, opt)
}
func (this *Redis) ZRangeByScore(valuetype reflect.Type, key string, opt *redis.ZRangeBy) (result []interface{}, err error) {
	return this.client.ZRangeByScore(valuetype, key, opt)
}
func (this *Redis) ZRank(key string, member string) (result int64, err error) {
	return this.client.ZRank(key, member)
}
func (this *Redis) ZRem(key string, members ...interface{}) (result int64, err error) {
	return this.client.ZRem(key, members...)
}
func (this *Redis) ZRemRangeByLex(key string, min string, max string) (result int64, err error) {
	return this.client.ZRemRangeByLex(key, min, max)
}
func (this *Redis) ZRemRangeByRank(key string, start int64, stop int64) (result int64, err error) {
	return this.client.ZRemRangeByRank(key, start, stop)
}
func (this *Redis) ZRemRangeByScore(key string, min string, max string) (result int64, err error) {
	return this.client.ZRemRangeByScore(key, min, max)
}
func (this *Redis) ZRevRange(valuetype reflect.Type, key string, start int64, stop int64) (result []interface{}, err error) {
	return this.client.ZRevRange(valuetype, key, start, stop)
}
func (this *Redis) ZRevRangeByScore(valuetype reflect.Type, key string, opt *redis.ZRangeBy) (result []interface{}, err error) {
	return this.client.ZRevRangeByScore(valuetype, key, opt)
}
func (this *Redis) ZRevRank(key string, member string) (result int64, err error) {
	return this.client.ZRevRank(key, member)
}
func (this *Redis) ZScore(key string, member string) (result float64, err error) {
	return this.client.ZScore(key, member)
}
func (this *Redis) ZUnionStore(dest string, store *redis.ZStore) (result int64, err error) {
	return this.client.ZUnionStore(dest, store)
}
func (this *Redis) ZScan(key string, _cursor uint64, match string, count int64) (keys []string, cursor uint64, err error) {
	return this.client.ZScan(key, _cursor, match, count)
}
