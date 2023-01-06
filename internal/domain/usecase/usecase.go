package usecase

import (
	"context"

	http "github.com/gobox-preegnees/connection_controller/internal/controller/http"
	entity "github.com/gobox-preegnees/connection_controller/internal/domain/entity"
	errors "github.com/gobox-preegnees/connection_controller/internal/errors"

	"github.com/sirupsen/logrus"
)

type IOwnerService interface {
	DeleteOwner(ctx context.Context, owner entity.Owner) (ownersCount int, err error)
	SaveOwner(ctx context.Context, owner entity.Owner) (ownersCount int, err error)
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
	ownerService       IOwnerService
	snapshotService    ISnapshotService
	consistencyService IConsistencyService
}

type CnfUsecase struct {
	Log                *logrus.Logger
	OwnerService       IOwnerService
	SnapshotService    ISnapshotService
	ConsistencyService IConsistencyService
}

func NewUsecase(cnf CnfUsecase) *usecase {

	return &usecase{
		log:                cnf.Log,
		ownerService:       cnf.OwnerService,
		snapshotService:    cnf.SnapshotService,
		consistencyService: cnf.ConsistencyService,
	}
}

func (u usecase) GetConsistency(ctx context.Context) (entity.Consistency, error) {

	return u.consistencyService.GetConsistency(ctx)
}

func (u usecase) SaveSnapshot(ctx context.Context, snapshot entity.Snapshot) error {

	return u.snapshotService.SaveSnapshot(ctx, snapshot)
}

func (u usecase) SaveOwner(ctx context.Context, owner entity.Owner) error {

	ownersCount, err := u.ownerService.SaveOwner(ctx, owner)
	if err != nil {
		return err
	}
	u.log.Debugf("ownerCount:%d on save owner:%v", ownersCount, owner)
	return nil
}

func (u usecase) DeleteOwner(ctx context.Context, owner entity.Owner) error {

	ownersCount, err := u.ownerService.DeleteOwner(ctx, owner)
	if err != nil {
		return err
	}
	u.log.Debugf("ownerCount:%d on delete owner:%v", ownersCount, owner)
	if ownersCount == -1 {
		return errors.ErrNoVisitors
	}
	return nil
}

var _ http.IUsecase = (*usecase)(nil)
