package kafka

import (
	"context"
	"errors"
	"io"

	entity "github.com/gobox-preegnees/connection_controller/internal/domain/entity"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -destination=../../../mocks/controller/message_broker/kafka/consumer/IUsecase/IUsecase.go -source=consumer.go
type IUsecase interface {
	SaveConsistency(ctx context.Context, consistency entity.Consistency) error
}

// consumer.
type consumer struct {
	ctx     context.Context
	log     *logrus.Logger
	reader  *kafka.Reader
	usecase IUsecase
}

// ConsumerCnf.
type ConsumerCnf struct {
	Ctx       context.Context
	Log       *logrus.Logger
	Topic     string
	Addrs     []string
	GroupId   string
	Partition int
	Usecase   IUsecase
}

// NewConsumer.
func NewConsumer(cnf ConsumerCnf) *consumer {

	// Проверка доступности
	if conn, err := kafka.Dial("tcp", cnf.Addrs[0]); err != nil {
		conn.Close()
		cnf.Log.Fatal(err)
	} else {
		conn.Close()
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cnf.Addrs,
		GroupID:  cnf.GroupId,
		Topic:    cnf.Topic,
		Logger:   cnf.Log,
	})

	cons :=  &consumer{
		ctx:     cnf.Ctx,
		log:     cnf.Log,
		reader:  reader,
		usecase: cnf.Usecase,
	}

	go func() {
		cons.stopOnDoneContext()
	}()

	return cons
}

// stopOnDoneContext.
func (c *consumer) stopOnDoneContext() {

	select {
	case <-c.ctx.Done():
		c.reader.Close()
		c.log.Debug("kafka consumer is stopped")
	}
}

// Run.
func (c *consumer) Run() error {

	for {
		message, err := c.reader.FetchMessage(c.ctx)
		if errors.Is(err, context.Canceled) {
			return nil
		} else if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		headers := message.Headers
		consistency := entity.Consistency{
			RequestId: string(headers[0].Value),
			StreamId:  string(headers[1].Value),
			Timestamp: message.Time,
			Data:      message.Value,
		}
		c.log.Debugf("time: %v, requestid: %v, streamid: %v, data: %v",
			consistency.Timestamp, consistency.RequestId, consistency.StreamId, consistency.Data)

		c.usecase.SaveConsistency(context.Background(), consistency)
		c.log.Debugf("success save consistency: %v", consistency)

		if err := c.reader.CommitMessages(c.ctx, message); err != nil {
			return err
		}
	}
}
