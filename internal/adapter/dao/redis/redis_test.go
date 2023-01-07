package redis

import (
	"context"
	"testing"

	storage "github.com/gobox-preegnees/connection_controller/internal/adapter/dao"
	"github.com/sirupsen/logrus"
)

const url = "redis://default:password@localhost:6379/0"

var ctx context.Context = context.Background()

func TestMain(t *testing.M) {

	t.Run()
}

func clear(r *redisClient) {
	r.client.FlushAll(ctx)
}

func TestCreateAndDelete(t *testing.T) {

	logger := logrus.New()
	r, err := NewRedisClient(CnfRedisClient{
		Ctx: ctx,
		Log: logger,
		Url: url,
	})
	if err != nil {
		t.Fatal(err)
	}
	clear(r)
	defer clear(r)

	streamId := "1"

	data := []struct {
		Action   bool // true = create, false = delete
		ExeptNum int
		Err      bool
	}{
		{
			Action:   true,
			ExeptNum: 1,
			Err:      false,
		},
		{
			Action:   true,
			ExeptNum: 2,
			Err:      false,
		},
		{
			Action:   false,
			ExeptNum: 1,
			Err:      false,
		},
		{
			Action:   false,
			ExeptNum: -1,
			Err:      false,
		},
		{
			Action:   false,
			ExeptNum: -1,
			Err:      false,
		},
		{
			Action:   true,
			ExeptNum: 1,
			Err:      false,
		},
	}

	for _, d := range data {
		t.Run("test", func(t *testing.T) {
			var num int
			var err error
			if d.Action {
				num, err = r.CreateOneStream(storage.CreateOneStreamReqDTO{
					Ctx:      ctx,
					StreamId: streamId,
				})
			} else {
				num, err = r.DeleteOneStream(storage.DeleteOneStreamReqDTO{
					Ctx:      ctx,
					StreamId: streamId,
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
