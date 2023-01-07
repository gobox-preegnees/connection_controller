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

// producer.
type producer struct {
	ctx       context.Context
	log       *logrus.Logger
	writer    *kafka.Writer
	attempts  int
	timeount  int
	sleeptime int
}

// ProducerConf.
type ProducerConf struct {
	Ctx       context.Context
	Log       *logrus.Logger
	Topic     string
	Addrs     []string
	Attempts  int
	Timeout   int
	Sleeptime int
}

// NewProducer.
func NewProducer(cnf ProducerConf) (*producer, error) {

	// Проверка доступности
	if conn, err := kafka.Dial("tcp", cnf.Addrs[0]); err != nil {
		return nil, err
	} else {
		conn.Close()
	}

	w := &kafka.Writer{
		Addr:                   kafka.TCP(cnf.Addrs...),
		Topic:                  cnf.Topic,
		AllowAutoTopicCreation: true,
		Logger:                 cnf.Log,
	}

	pr := &producer{
		ctx:       cnf.Ctx,
		log:       cnf.Log,
		writer:    w,
		attempts:  cnf.Attempts,
		timeount:  cnf.Timeout,
		sleeptime: cnf.Sleeptime,
	}

	go func() {
		pr.stopOnDoneContext()
	}()

	return pr, nil
}

// stopOnDoneContext.
func (p *producer) stopOnDoneContext() {

	select {
	case <-p.ctx.Done():
		p.writer.Close()
		p.log.Debug("kafka producer is stopped")
	}
}

// SaveSnapshot.
func (p producer) SaveSnapshot(req mbDTO.PublishSnapshotReqDTO) error {

	// https://github.com/segmentio/kafka-go#missing-topic-creation-before-publication
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
		p.log.Debugf("producer sent message:%v", req)
	}
	return nil
}

var _ service.ISnapshotMessageBroker = (*producer)(nil)
