package redis_distributed_lock

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v9"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mocks "redis-distributed-lock/mock"
	"testing"
	"time"
)

//单元测试使用mock
//先下载gomock:
//go get github.com/golang/mock/gomock
//然后安装:
//go install github.com/golang/mock/gomock
//安装好mock后在本地目录使用命令mockgen生成mock文件
//mockgen -package=mocks -destination=/mock/redis_cmdable.mock.go github.com/go-redis/redis/v9 Cmdable
//destination：mockgen生成的文件存放的位置以及文件的名字。
//package：生成的mock文件的包名。
//goMockTest/order：module名字+ 包名字。
//OrderService：需要mock的接口的名称。

func TestClient_TryLock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	testCases := []struct {
		name string // 测试用例的场景
		//输入
		key          string        // 测试的输入
		expiration   time.Duration // 测试的输入
		expectedErr  error         // 测试预期输出
		expectedLock *Lock         // 测试预期输出

		mock func() redis.Cmdable // 设置mock的数据
	}{
		{
			//模拟正常加锁成功
			name:       "locked",
			key:        "locked-key",
			expiration: time.Minute,
			mock: func() redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				//mock中预期的res
				res := redis.NewBoolResult(true, nil)
				//执行mock的预期
				cmd.EXPECT().SetNX(gomock.Any(), "locked-key", gomock.Any(), time.Minute).Return(res)
				return cmd
			},
			expectedLock: &Lock{
				key: "locked-key",
			},
		},
		{
			//模拟加锁错误（网络错误）
			name:       "network error",
			key:        "network-key",
			expiration: time.Minute,
			mock: func() redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				//mock中预期的res
				res := redis.NewBoolResult(true, errors.New("network error"))
				//执行mock的预期
				cmd.EXPECT().SetNX(gomock.Any(), "network-key", gomock.Any(), time.Minute).Return(res)
				return cmd
			},
			expectedLock: &Lock{
				key: "network-key",
			},
			expectedErr: errors.New("network error"),
		},
		{
			//模拟加锁没有错误，但是失败
			name:       "fail to lock",
			key:        "fail-key",
			expiration: time.Minute,
			mock: func() redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				//mock中预期的res
				res := redis.NewBoolResult(true, ErrTryLockFail)
				//执行mock的预期
				cmd.EXPECT().SetNX(gomock.Any(), "fail-key", gomock.Any(), time.Minute).Return(res)
				return cmd
			},
			expectedLock: &Lock{
				key: "fail-key",
			},
			expectedErr: ErrTryLockFail,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// mock模拟Client
			c := NewClient(tc.mock())
			// mock模拟调用Trylock
			l, err := c.TryLock(context.Background(), tc.key, tc.expiration)
			//testcase中的预期结果与mock模拟的结果做比较
			//比较err
			assert.Equal(t, tc.expectedErr, err)
			//如果有err可以直接return不用再比较其他的了
			if err != nil {
				return
			}
			assert.NotEmpty(t, l.client)
			assert.Equal(t, tc.expectedLock.key, l.key)
			assert.NotEmpty(t, l.value)
		})
	}
}
