package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"project-go/webook/internal/domain"
	"time"
)

var ErrKeyNotExist = redis.Nil

type UserCache struct {
	// 传单机 Redis可以
	// 传 cluster 的 Redis 也可以
	client     redis.Cmdable
	expiration time.Duration
}

// NewUserCache
// A 用到了B B一定是接口 => 这个是保证面向接口
// A 用到了B B一定是A的字段 => 规避包变量、包方法，都非常缺乏扩展性
// A 用到了B A 绝对不初始化B，而是外面注入 => 保持依赖注入和依赖反转
func NewUserCache(client redis.Cmdable) *UserCache {
	return &UserCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}

func (cache *UserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.key(id)
	val, err := cache.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	// 这里会将val 填充到 &u中
	err = json.Unmarshal(val, &u)
	return u, err
}

func (cache *UserCache) Set(ctx context.Context, u domain.User) error {
	// 把对象转换为json串
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}
	key := cache.key(u.Id)
	return cache.client.Set(ctx, key, val, cache.expiration).Err()
}

func (cache *UserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
