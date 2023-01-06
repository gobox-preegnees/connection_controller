package kafka

import (
	"context"
	"errors"
	"time"

	mbDTO "github.com/gobox-preegnees/connection_controller/internal/adapter/message_broker"
	service "github.com/gobox-preegnees/connection_controller/internal/domain/service"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/protocol"
	"github.com/sirupsen/logrus"
)

type producer struct {
	log       *logrus.Logger
	writer    *kafka.Writer
	attempts  int
	timeount  int
	sleeptime int
}

type ProducerConf struct {
	Log       *logrus.Logger
	ErrTopic  string
	Addrs     []string
	Attempts  int
	Timeout   int
	Sleeptime int
}

func NewProducer(cnf ProducerConf) *producer {

	w := &kafka.Writer{
		Addr:                   kafka.TCP(cnf.Addrs...),
		Topic:                  cnf.ErrTopic,
		AllowAutoTopicCreation: true,
		Logger:                 cnf.Log,
	}
	return &producer{
		log:       cnf.Log,
		writer:    w,
		attempts:  cnf.Attempts,
		timeount:  cnf.Timeout,
		sleeptime: cnf.Sleeptime,
	}
}

func (p producer) CreateOneSnapshot(req mbDTO.PublishSnapshotReqDTO) error {

	for i := 0; i < p.attempts; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.timeount)*time.Second)
		defer cancel()

		err := p.writer.WriteMessages(ctx, kafka.Message{
			Headers: []protocol.Header{
				{
					Key:   "RequestId",
					Value: []byte(req.RequestId),
				},
				{
					Key:   "StreamId",
					Value: []byte(req.StreamId),
				},
			},
			Value: req.Data,
			Time:  req.Timestamp,
		})
		if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
			time.Sleep(time.Duration(p.sleeptime) * time.Millisecond)
			continue
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (p producer) Close() error {

	return p.writer.Close()
}

var _ service.ISnapshotMessageBroker = (*producer)(nil)
