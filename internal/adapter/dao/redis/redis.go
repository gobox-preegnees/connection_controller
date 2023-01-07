package redis

import (
	"context"
	"time"

	daoDTO "github.com/gobox-preegnees/connection_controller/internal/adapter/dao"
	service "github.com/gobox-preegnees/connection_controller/internal/domain/service"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

// redisClient.
type redisClient struct {
	ctx             context.Context
	log             *logrus.Logger
	client          *redis.Client
	shutdonwTimeout int
}

// CnfRedisClient.
type CnfRedisClient struct {
	Ctx             context.Context
	Log             *logrus.Logger
	Url             string
	ShutdownTimeout int
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

// NewRedisClient.
func NewRedisClient(cnf CnfRedisClient) (*redisClient, error) {

	opt, err := redis.ParseURL(cnf.Url)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)
	status := client.Ping(cnf.Ctx)
	if status.Err() != nil {
		return nil, status.Err()
	}

	rCli := &redisClient{
		ctx:             cnf.Ctx,
		log:             cnf.Log,
		client:          client,
		shutdonwTimeout: cnf.ShutdownTimeout,
	}

	go func() {
		rCli.stopOnDoneContext()
	}()

	return rCli, nil
}

// stopOnDoneContext.
func (r *redisClient) stopOnDoneContext() {

	select {
	case <-r.ctx.Done():
		r.client.Close()

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.shutdonwTimeout)*time.Second)
		defer cancel()
		r.client.Shutdown(ctx)
		r.log.Debug("redis client is stopped")
	}
}

// CreateOneStream. Если такого стрима (ключа) еще нет в базе, то он создастся и значение будет равно "1".
// Если такой стрим уже есть в базе, то значение увеличится на 1.
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

// DeleteOneStream. Декреметит значение при удалении стрима, если значение становится <= 0, то в ответ содержит -1
func (r redisClient) DeleteOneStream(req daoDTO.DeleteOneStreamReqDTO) (int, error) {

	num, err := decrBy.Run(
		req.Ctx,
		r.client,
		[]string{req.StreamId},
	).Int()
	if err != nil {
		return -1, err
	}
	r.log.Debugf("decr by stream=%s: current_connections=%d\n", req.StreamId, num)
	if num == -1 {
		return -1, nil
	}
	return num, nil
}

var _ service.IStreamDao = (*redisClient)(nil)
