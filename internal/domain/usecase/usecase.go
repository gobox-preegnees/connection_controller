package usecase

import (
	"context"
	"errors"

	http "github.com/gobox-preegnees/connection_controller/internal/controller/http"
	entity "github.com/gobox-preegnees/connection_controller/internal/domain/entity"
)

var ErrNoVisitors = errors.New("no visitors")

type IOwnerService interface {
	DeleteOwner(ctx context.Context, owner entity.Owner) (bool, error)
	SaveOwner(ctx context.Context, owner entity.Owner) error
}

type IConsistencyService interface {
	GetConsistency() entity.Consistency
	SaveConsistency(entity.Consistency) error
}

type ISnapshotService interface {
	SaveSnapshot(ctx context.Context, snapshot entity.Snapshot) error
}

type usecase struct {
	ownerService       IOwnerService
	snapshotService    ISnapshotService
	consistencyService IConsistencyService
}

func (u usecase) DeleteOwner(ctx context.Context, owner entity.Owner) error {

	isNoVisitors, err := u.ownerService.DeleteOwner(ctx, owner)
	if isNoVisitors {
		return ErrNoVisitors
	}
	return err
}

func (u usecase) SaveOwner(ctx context.Context, owner entity.Owner) error {
	
	return u.ownerService.SaveOwner(ctx, owner)
}

func (u usecase) SaveSnapshot(ctx context.Context, snapshot entity.Snapshot) error {
	
	return u.snapshotService.SaveSnapshot(ctx, snapshot)
}

func (u usecase) GetConsistency() entity.Consistency {
	
	return u.consistencyService.GetConsistency()
}

func (u usecase) SaveConsistency(consistency entity.Consistency) error {
	
	return u.consistencyService.SaveConsistency(consistency)
}

var _ http.IUsecase = (*usecase)(nil)
