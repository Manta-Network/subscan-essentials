package dao

import (
	"context"
	"reflect"

	"github.com/go-kratos/kratos/pkg/cache/redis"
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/itering/subscan/util"
	"github.com/pkg/errors"
)

func (d *Dao) SetMetadata(c context.Context, metadata map[string]interface{}) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	args := redis.Args{}.Add(RedisMetadataKey)
	if len(metadata) == 0 {
		return errors.New("ERR: nil metadata")
	}
	for k, v := range metadata {
		if reflect.ValueOf(v).Kind() == reflect.Int {
			args = args.Add(k).Add(util.IntToString(v.(int)))
		} else {
			args = args.Add(k).Add(v)
		}
	}
	_, err = conn.Do("HMSET", args...)
	if err != nil {
		log.Error("HMSET ERROR", err)
	}
	return
}

func (d *Dao) IncrMetadata(c context.Context, filed string, incrNum int) (err error) {
	if incrNum == 0 {
		return
	}
	conn := d.redis.Get(c)
	defer conn.Close()
	_, err = conn.Do("HINCRBY", RedisMetadataKey, filed, incrNum)
	if err != nil {
		log.Error("HINCRBY ERROR", err)
	}
	return
}

func (d *Dao) GetMetadata(c context.Context) (ms map[string]string, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	ms, err = redis.StringMap(conn.Do("HGETALL", RedisMetadataKey))
	return
}

func (d *Dao) GetBestBlockNum(c context.Context) (uint64, error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	return redis.Uint64(conn.Do("HGET", RedisMetadataKey, "blockNum"))
}

func (d *Dao) GetFinalizedBlockNum(c context.Context) (uint64, error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	return redis.Uint64(conn.Do("HGET", RedisMetadataKey, "finalized_blockNum"))
}
