package main

import (
	"context"
	"os"

	redisAdapter "github.com/gobox-preegnees/connection_controller/internal/adapter/dao/redis"
	kafkaAdapter "github.com/gobox-preegnees/connection_controller/internal/adapter/message_broker/kafka"

	httpController "github.com/gobox-preegnees/connection_controller/internal/controller/http"
	kafkaController "github.com/gobox-preegnees/connection_controller/internal/controller/message_broker/kafka"

	service "github.com/gobox-preegnees/connection_controller/internal/domain/service"
	usecase "github.com/gobox-preegnees/connection_controller/internal/domain/usecase"

	cnf "github.com/gobox-preegnees/connection_controller/internal/config"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func main() {

	path := os.Args[1]
	if path == "" {
		panic("path (arg 1) is empty")
	}

	ctx, cancel := context.WithCancel(context.Background())
	logger := logrus.New()

	config := cnf.GetGonfig(path)
	logger.Debugf("config: %v", config)

	if config.Debug {
		logger.SetLevel(logrus.DebugLevel)
		logger.SetReportCaller(true)
	}

	streamDao, err := redisAdapter.NewRedisClient(redisAdapter.CnfRedisClient{
		Ctx: ctx,
		Log: logger,
		Url: config.Redis.Url,
	})
	if err != nil {
		logger.Fatal(err)
	}

	snapshotMessageBroker, err := kafkaAdapter.NewProducer(kafkaAdapter.ProducerConf{
		Log:       logger,
		Topic:     config.Kafka.Producer.SnapshotTopic,
		Addrs:     config.Kafka.Addrs,
		Attempts:  config.Kafka.Producer.Attempts,
		Timeout:   config.Kafka.Producer.Timeout,
		Sleeptime: config.Kafka.Producer.Sleeptime,
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
		Addr:      config.Http.Addr,
		Usecase:   uc,
		JWTAlg:    config.Http.JWTAlg,
		JWTSecret: config.Http.Secret,
		CrtPath:   config.Http.CrtPath,
		KeyPath:   config.Http.KeyPath,
	})
	kController := kafkaController.NewConsumer(kafkaController.ConsumerCnf{
		Ctx:       ctx,
		Log:       logger,
		Topic:     config.Kafka.Consumer.ConsistencyTopic,
		Addrs:     config.Kafka.Addrs,
		GroupId:   config.Kafka.Consumer.GroupId,
		Partition: config.Kafka.Consumer.Partition,
		Usecase:   uc,
	})

	g := new(errgroup.Group)
	stopCh := make(chan struct{})

	g.Go(func() error {
		errCh := make(chan error)
		go func() {
			errCh <- hController.Run()
		}()
		select {
		case <-stopCh:
			cancel()
			return nil
		case err := <-errCh:
			stopCh <- struct{}{}
			return err
		}
	})

	g.Go(func() error {
		errCh := make(chan error)
		go func() {
			errCh <- kController.Run()
		}()
		select {
		case <-stopCh:
			cancel()
			return nil
		case err := <-errCh:
			stopCh <- struct{}{}
			return err
		}
	})

	if err := g.Wait(); err != nil {
		close(stopCh)
		logger.Fatal(err)
	}
}
