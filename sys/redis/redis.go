package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

func newSys(options Options) (sys *Redis, err error) {
	sys = &Redis{options: options}
	err = sys.init()
	return
}

type Redis struct {
	options Options
	client  *redis.Client
}

func (this *Redis) init() (err error) {
	this.client = redis.NewClient(&redis.Options{
		Addr:     this.options.RedisUrl,
		Password: this.options.RedisPassword,
		DB:       this.options.RedisDB,
	})
	_, err = this.client.Ping(this.getContext()).Result()
	return
}

func (this *Redis) getContext() (ctx context.Context) {
	ctx, _ = context.WithTimeout(context.Background(), this.options.TimeOut)
	return
}

/* Key *******************************************************************************/

///删除redis key
func (this *Redis) Delete(key string) (err error) {
	err = this.client.Do(this.getContext(), "DEL", key).Err()
	return
}

///判断是否存在key
func (this *Redis) ExistsKey(key string) (iskeep bool, err error) {
	iskeep, err = this.client.Do(this.getContext(), "EXISTS", key).Bool()
	return
}

///设置key的过期时间 单位以秒级
func (this *Redis) ExpireKey(key string, expire int) (err error) {
	err = this.client.Do(this.getContext(), "EXPIRE", key, expire).Err()
	return
}

///设置key的过期时间戳 秒级时间戳
func (this *Redis) ExpireatKey(key string, expire_unix int64) (err error) {
	err = this.client.Do(this.getContext(), "EXPIREAT", key, expire_unix).Err()
	return
}

///设置key的过期时间 单位以毫秒级
func (this *Redis) Pexpirekey(key string, expire int) (err error) {
	err = this.client.Do(this.getContext(), "PEXPIRE", key, expire).Err()
	return
}

///设置key的过期时间戳 单位以豪秒级
func (this *Redis) PexpireatKey(key string, expire_unix int64) (err error) {
	err = this.client.Do(this.getContext(), "PEXPIREAT", key, expire_unix).Err()
	return
}

///移除Key的过期时间
func (this *Redis) PersistKey(key string) (err error) {
	err = this.client.Do(this.getContext(), "PERSIST", key).Err()
	return
}

///获取key剩余过期时间 单位毫秒
func (this *Redis) PttlKey(key string) (leftexpire int64, err error) {
	leftexpire, err = this.client.Do(this.getContext(), "PTTL", key).Int64()
	return
}

///获取key剩余过期时间 单位秒
func (this *Redis) TtlKey(key string) (leftexpire int64, err error) {
	leftexpire, err = this.client.Do(this.getContext(), "TTL", key).Int64()
	return
}

///重命名Key
func (this *Redis) RenameKye(oldkey string, newkey string) (err error) {
	err = this.client.Do(this.getContext(), "RENAME", oldkey, newkey).Err()
	return
}

///重命名key 在新的 key 不存在时修改 key 的名称
func (this *Redis) RenamenxKey(oldkey string, newkey string) (err error) {
	err = this.client.Do(this.getContext(), "RENAMENX", oldkey, newkey).Err()
	return
}

///判断是否存在key pattern:key*
func (this *Redis) Keys(pattern string) (keys []string, err error) {
	cmd := redis.NewStringSliceCmd(this.getContext(), "KEYS", pattern)
	this.client.Process(this.getContext(), cmd)
	keys, err = cmd.Result()
	return
}

/* String *******************************************************************************/

///命令用于设置给定 key 的值。如果 key 已经存储其他值， SET 就覆写旧值，且无视类型。
func (this *Redis) Set(key string, value interface{}) (err error) {
	// this.client.SetArgs()
	return
}

///命令用于设置给定 key 的值。如果 key 已经存储其他值， SET 就覆写旧值，且无视类型。
func (this *Redis) Get(key string, value interface{}) (err error) {
	// this.client.SetArgs()
	return
}

///
func (this *Redis) Getset(key string) (err error) {
	return
}
