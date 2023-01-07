package usecase

import (
	"context"

	http "github.com/gobox-preegnees/connection_controller/internal/controller/http"
	mb "github.com/gobox-preegnees/connection_controller/internal/controller/message_broker/kafka"
	entity "github.com/gobox-preegnees/connection_controller/internal/domain/entity"
	errors "github.com/gobox-preegnees/connection_controller/internal/errors"

	"github.com/sirupsen/logrus"
)

//go:generate mockgen -destination=../../mocks/domain/usecase/IStreamService/IStreamService.go -source=usecase.go
type IStreamService interface {
	DeleteStream(ctx context.Context, stream entity.Stream) (numberOfUsers int, err error)
	SaveStream(ctx context.Context, stream entity.Stream) (numberOfUsers int, err error)
}

//go:generate mockgen -destination=../../mocks/domain/usecase/usecase/IConsistencyService/IConsistencyService.go -source=usecase.go
type IConsistencyService interface {
	GetConsistency(ctx context.Context) (consistency entity.Consistency, err error)
	SaveConsistency(ctx context.Context, consistency entity.Consistency) (err error)
}

//go:generate mockgen -destination=../../mocks/domain/usecase/ISnapshotService/ISnapshotService.go -source=usecase.go
type ISnapshotService interface {
	SaveSnapshot(ctx context.Context, snapshot entity.Snapshot) (err error)
}

// usecase.
type usecase struct {
	log                *logrus.Logger
	streamService      IStreamService
	snapshotService    ISnapshotService
	consistencyService IConsistencyService
}

// CnfUsecase.
type CnfUsecase struct {
	Log                *logrus.Logger
	StreamService      IStreamService
	SnapshotService    ISnapshotService
	ConsistencyService IConsistencyService
}

// NewUsecase.
func NewUsecase(cnf CnfUsecase) *usecase {

	return &usecase{
		log:                cnf.Log,
		streamService:      cnf.StreamService,
		snapshotService:    cnf.SnapshotService,
		consistencyService: cnf.ConsistencyService,
	}
}

// GetConsistency.
func (u usecase) GetConsistency(ctx context.Context) (entity.Consistency, error) {

	return u.consistencyService.GetConsistency(ctx)
}

// SaveConsistency.
func (u usecase) SaveConsistency(ctx context.Context, consistency entity.Consistency) error {

	return u.consistencyService.SaveConsistency(ctx, consistency)
}

// SaveSnapshot.
func (u usecase) SaveSnapshot(ctx context.Context, snapshot entity.Snapshot) error {

	return u.snapshotService.SaveSnapshot(ctx, snapshot)
}

// SaveStream.
func (u usecase) SaveStream(ctx context.Context, stream entity.Stream) error {

	numberOfUsers, err := u.streamService.SaveStream(ctx, stream)
	if err != nil {
		return err
	}
	u.log.Debugf("numberOfUsers:%d on save streamId:%v", numberOfUsers, stream.StreamId)
	return nil
}

// DeleteStream.
func (u usecase) DeleteStream(ctx context.Context, stream entity.Stream) error {

	numberOfUsers, err := u.streamService.DeleteStream(ctx, stream)
	if err != nil {
		return err
	}
	u.log.Debugf("numberOfUsers:%d on delete streamId:%v", numberOfUsers, stream.StreamId)
	if numberOfUsers == -1 {
		return errors.ErrNoVisitors
	}
	return nil
}

var _ http.IUsecase = (*usecase)(nil)
var _ mb.IUsecase = (*usecase)(nil)
