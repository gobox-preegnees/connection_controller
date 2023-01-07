package redis

import (
	"context"
	"fmt"

	daoDTO "github.com/gobox-preegnees/connection_controller/internal/adapter/dao"
	service "github.com/gobox-preegnees/connection_controller/internal/domain/service"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type redisClient struct {
	log    *logrus.Logger
	client *redis.Client
}

type CnfRedisClient struct {
	Ctx context.Context
	Log *logrus.Logger
	Url string
}

var incrBy = redis.NewScript(
	`
		local key = KEYS[1]

		local value = redis.call("GET", key)
		if not value then
			value = 0
		end

		value = value + 1
		redis.call("SET", key, value)

		return value
	`,
)

var decrBy = redis.NewScript(
	`
		local key = KEYS[1]

		local value = redis.call("GET", key)
		if not value then
			return -1
		end

		value = value - 1
		if value <= 0 then
			redis.call("DEL", key)
			return -1
		end

		redis.call("SET", key, value)

		return value
	`,
)

func NewRedisClient(cnf CnfRedisClient) *redisClient {

	opt, err := redis.ParseURL(cnf.Url)
	if err != nil {
		cnf.Log.Fatal(err)
	}

	client := redis.NewClient(opt)
	status := client.Ping(cnf.Ctx)
	if status.Err() != nil {
		cnf.Log.Fatal(status.Err())
	}
	return &redisClient{
		log:    cnf.Log,
		client: client,
	}
}

func (r redisClient) CreateOneStream(req daoDTO.CreateOneStreamReqDTO) (int, error) {

	num, err := incrBy.Run(
		req.Ctx,
		r.client,
		[]string{req.StreamId},
	).Int()
	if err != nil {
		return -1, err
	}
	r.log.Debugf("incr by stream=%s: current_connections=%d\n", req.StreamId, num)
	return num, nil
}

func (r redisClient) DeleteOneStream(req daoDTO.DeleteOneStreamReqDTO) (int, error) {

	num, err := decrBy.Run(
		req.Ctx,
		r.client,
		[]string{req.StreamId},
	).Int()
	if err != nil {
		return -1, err
	}
	fmt.Printf("decr by stream=%s: current_connections=%d\n", req.StreamId, num)
	if num == -1 {
		return -1, nil
	}
	return num, nil
}

var _ service.IStreamDao = (*redisClient)(nil)
