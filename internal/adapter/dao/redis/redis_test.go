package redis

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v8"
	storage "github.com/gobox-preegnees/connection_controller/internal/adapter/dao"
	"github.com/sirupsen/logrus"
)

var client *redis.Client

const url = "redis://default:password@localhost:6379/0"

var ctx context.Context = context.Background()

func TestMain(t *testing.M) {

	opt, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opt)
	status := client.Ping(ctx)
	if status.Err() != nil {
		panic(status.Err())
	}
	defer client.Close()

	t.Run()
}

func clear(r *redisClient) {
	r.client.FlushAll(ctx)
}

func TestCreateAndDelete(t *testing.T) {

	logger := logrus.New()
	r := NewRedisClient(CnfRedisClient{
		Ctx: ctx,
		Log: logger,
		Url: url,
	})
	clear(r)
	defer clear(r)

	username := "user1"
	folder := "folder1"

	data := []struct {
		Action   bool // true = create, false = delete
		ExeptNum int
		Err      bool
	}{
		{
			Action: true,
			ExeptNum: 1,
			Err: false,
		},
		{
			Action: true,
			ExeptNum: 2,
			Err: false,
		},
		{
			Action: false,
			ExeptNum: 1,
			Err: false,
		},
		{
			Action: false,
			ExeptNum: -1,
			Err: false,
		},
		{
			Action: false,
			ExeptNum: -1,
			Err: false,
		},
		{
			Action: true,
			ExeptNum: 1,
			Err: false,
		},
	}

	for _, d := range data {
		t.Run("test", func(t *testing.T) {
			var num int
			var err error
			if d.Action {
				num, err = r.CreateOneOwner(storage.CreateOneOwnerReqDTO{
					Ctx:       ctx,
					Usernamme: username,
					Folder:    folder,
				})
			} else {
				num, err = r.DeleteOneOwner(storage.DeleteOneOwnerReqDTO{
					Ctx:       ctx,
					Usernamme: username,
					Folder:    folder,
				})
			}

			if err != nil {
				if !d.Err {
					t.Errorf("expect none, got err:%v", err)
				} 
			}
			if num != d.ExeptNum {
				t.Errorf("expect %d, got %d", d.ExeptNum, num)
			}
		})
	}
}