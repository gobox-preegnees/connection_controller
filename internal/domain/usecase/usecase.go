package usecase

import (
	"context"

	http "github.com/gobox-preegnees/connection_controller/internal/controller/http"
	mb "github.com/gobox-preegnees/connection_controller/internal/controller/message_broker/kafka"
	entity "github.com/gobox-preegnees/connection_controller/internal/domain/entity"
	errors "github.com/gobox-preegnees/connection_controller/internal/errors"

	"github.com/sirupsen/logrus"
)

type IStreamService interface {
	DeleteStream(ctx context.Context, stream entity.Stream) (ownersCount int, err error)
	SaveStream(ctx context.Context, stream entity.Stream) (ownersCount int, err error)
}

type IConsistencyService interface {
	GetConsistency(ctx context.Context) (consistency entity.Consistency, err error)
	SaveConsistency(ctx context.Context, consistency entity.Consistency) (err error)
}

type ISnapshotService interface {
	SaveSnapshot(ctx context.Context, snapshot entity.Snapshot) (err error)
}

type usecase struct {
	log                *logrus.Logger
	streamService      IStreamService
	snapshotService    ISnapshotService
	consistencyService IConsistencyService
}

type CnfUsecase struct {
	Log                *logrus.Logger
	StreamService      IStreamService
	SnapshotService    ISnapshotService
	ConsistencyService IConsistencyService
}

func NewUsecase(cnf CnfUsecase) *usecase {

	return &usecase{
		log:                cnf.Log,
		streamService:      cnf.StreamService,
		snapshotService:    cnf.SnapshotService,
		consistencyService: cnf.ConsistencyService,
	}
}

func (u usecase) GetConsistency(ctx context.Context) (entity.Consistency, error) {

	return u.consistencyService.GetConsistency(ctx)
}

func (u usecase) SaveConsistency(ctx context.Context, consistency entity.Consistency) error {

	return u.consistencyService.SaveConsistency(ctx, consistency)
}

func (u usecase) SaveSnapshot(ctx context.Context, snapshot entity.Snapshot) error {

	return u.snapshotService.SaveSnapshot(ctx, snapshot)
}

func (u usecase) SaveStream(ctx context.Context, stream entity.Stream) error {

	ownersCount, err := u.streamService.SaveStream(ctx, stream)
	if err != nil {
		return err
	}
	u.log.Debugf("ownerCount:%d on save streamId:%v", ownersCount, stream.StreamId)
	return nil
}

func (u usecase) DeleteStream(ctx context.Context, stream entity.Stream) error {

	ownersCount, err := u.streamService.DeleteStream(ctx, stream)
	if err != nil {
		return err
	}
	u.log.Debugf("ownerCount:%d on delete streamId:%v", ownersCount, stream.StreamId)
	if ownersCount == -1 {
		return errors.ErrNoVisitors
	}
	return nil
}

var _ http.IUsecase = (*usecase)(nil)
var _ mb.IUsecase = (*usecase)(nil)
