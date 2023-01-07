package main

import (
	"context"

	redisAdapter "github.com/gobox-preegnees/connection_controller/internal/adapter/dao/redis"
	kafkaAdapter "github.com/gobox-preegnees/connection_controller/internal/adapter/message_broker/kafka"

	httpController "github.com/gobox-preegnees/connection_controller/internal/controller/http"
	kafkaController "github.com/gobox-preegnees/connection_controller/internal/controller/message_broker/kafka"

	service "github.com/gobox-preegnees/connection_controller/internal/domain/service"
	usecase "github.com/gobox-preegnees/connection_controller/internal/domain/usecase"

	"github.com/sirupsen/logrus"
	// "golang.org/x/sync/errgroup"
)

func main() {

	const url = "redis://default:password@localhost:6379/0"
	ctx := context.TODO()
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	// logger.SetReportCaller(true)

	streamDao := redisAdapter.NewRedisClient(redisAdapter.CnfRedisClient{
		Ctx: ctx,
		Log: logger,
		Url: url,
	})
	snapshotMessageBroker := kafkaAdapter.NewProducer(kafkaAdapter.ProducerConf{
		Log:       logger,
		Topic:     "snapshot",
		Addrs:     []string{"localhost:29092"},
		Attempts:  3,
		Timeout:   10,
		Sleeptime: 250,
	})

	streamService := service.NewStreamService(service.CnfStreamService{
		Log:       logger,
		StreamDao: streamDao,
	})
	snapshotService := service.NewShanpshotService(service.CnfSnapshotService{
		Log:                   logger,
		SnapshotMessageBroker: snapshotMessageBroker,
	})
	consistencyService := service.NewConsistensyService(service.CnfConsistensyService{
		Log: logger,
	})

	uc := usecase.NewUsecase(usecase.CnfUsecase{
		Log:                logger,
		StreamService:      streamService,
		SnapshotService:    snapshotService,
		ConsistencyService: consistencyService,
	})

	hController := httpController.NewhttpServer(httpController.CnfhttpServer{
		Ctx:       ctx,
		Log:       logger,
		Addr:      "localhost:6060",
		Usecase:   uc,
		JWTAlg:    "HS256",
		JWTSecret: "secret",
		CrtPath:   "server.crt",
		KeyPath:   "server.key",
	})
	_ = kafkaController.NewConsumer(kafkaController.ConsumerCnf{
		Ctx:       ctx,
		Log:       logger,
		Topic:     "consistency",
		Addrs:     []string{"localhost:29092"},
		GroupId:   "groupId",
		Partition: 0,
		Usecase:   uc,
	})

	hController.Run()
	// g1 := new(errgroup.Group)
	// g1.Go(func() error {

	// 	return nil
	// })
	// g1.Go(func() error {
	// 	return kController.Run()
	// })

	// if err := g1.Wait(); err != nil {
	// 	logger.Fatal(err)
	// }
}
