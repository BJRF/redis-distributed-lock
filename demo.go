package redis_distributed_lock

import (
	"context"
	_ "embed"
	"errors"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"time"
)

var (
	//go:embed unlock.lua
	ualUnlock string

	ErrTryLockFail = errors.New("加锁失败")
	ErrLockNotHold = errors.New("该用户未持有锁")
)

// Client 客户端结构
type Client struct {
	client redis.Cmdable
}

// NewClient 建立新的客户端
func NewClient(c redis.Cmdable) *Client {
	return &Client{
		client: c,
	}
}

// TryLock 尝试加锁
func (c *Client) TryLock(ctx context.Context, key string, expiration time.Duration) (*Lock, error) {
	value := uuid.New().String()
	res, err := c.client.SetNX(ctx, key, value, expiration).Result()
	//加锁错误
	if err != nil {
		return nil, err
	}
	//加锁没有错误，但是失败
	if !res {
		return nil, ErrTryLockFail
	}
	//加锁成功
	return newLock(c.client, key, value, expiration), nil
}

// Lock 锁的结构
type Lock struct {
	client     redis.Cmdable
	key        string
	value      string
	expiration time.Duration
}

// newLock new一个锁变量
func newLock(c redis.Cmdable, key string, value string, expiration time.Duration) *Lock {
	return &Lock{
		client:     c,
		key:        key,
		value:      value,
		expiration: expiration,
	}
}

// Refresh 续约分布式锁
func Refresh(ctx context.Context) error {
	//lua脚本判断是不是自己的锁，如果是则刷新

	return nil
}

// Unlock 解锁
func (l *Lock) Unlock(ctx context.Context, key string) error {
	////非原子操作，保证原子操作需要用到lua脚本
	//val, err := l.client.Get(ctx, key).Result()
	//if err != nil {
	//	return err
	//}
	//if l.value == val {
	//	_, err := l.client.Del(ctx, key).Result()
	//	if err != nil {
	//		return err
	//	}
	//}
	//return nil

	//lua脚本原子解锁操作
	res, err := l.client.Eval(ctx, "/unlock.lua", []string{l.key}, l.value).Int64()
	if err == redis.Nil {
		return ErrLockNotHold
	}
	if err != nil {
		return err
	}
	if res == 0 {
		//删除失败，锁不是你的或者key不存在
		return ErrLockNotHold
	}
	return nil
}
