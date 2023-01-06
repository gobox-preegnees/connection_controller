package kafka

import (
	"context"
	"errors"
	"io"

	entity "github.com/gobox-preegnees/connection_controller/internal/domain/entity"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

// go:generate mockgen -destination=../../mocks/kafka/consumer/state/usecase.go -package=kafka -source=consumer.go
type IUsecase interface {
	SaveConsistency(ctx context.Context, consistency entity.Consistency) error
}

type consumer struct {
	ctx     context.Context
	log     *logrus.Logger
	reader  *kafka.Reader
	usecase IUsecase
}

type ConsumerCnf struct {
	Ctx       context.Context
	Log       *logrus.Logger
	Topic     string
	Addrs     []string
	GroupId   string
	Partition int
	Usecase   IUsecase
}

func NewConsumer(cnf ConsumerCnf) *consumer {

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
		MinBytes: 0,
	})

	return &consumer{
		ctx:     cnf.Ctx,
		log:     cnf.Log,
		reader:  reader,
		usecase: cnf.Usecase,
	}
}

func (k *consumer) Run() error {

	defer k.reader.Close()

	for {
		message, err := k.reader.FetchMessage(k.ctx)
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
		k.log.Debugf("time: %v, requestid: %v, streamid: %v, data: %v",
			consistency.Timestamp, consistency.RequestId, consistency.StreamId, consistency.Data)

		k.usecase.SaveConsistency(context.Background(), consistency)
		k.log.Debugf("success save consistency: %v", consistency)

		if err := k.reader.CommitMessages(k.ctx, message); err != nil {
			return err
		}
	}
}
