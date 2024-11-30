package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"project-go/webook/internal/domain"
	"time"
)

type UserCache struct {
	// 传单机 Redis可以
	// 传 cluster 的 Redis 也可以
	client     redis.Cmdable
	expiration time.Duration
}

// NewUserCache
// A 用到了B B一定是接口
// A 用到了B B一定是A的字段
// A 用到了B A 绝对不初始化B，而是外面注入
func NewUserCache(client redis.Cmdable) *UserCache {
	return &UserCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}

//func (cache *UserCache) GetUser(ctx context.Context, id int64) (domain.User, error) {
//
//}

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
