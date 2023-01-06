package service

import (
	"context"

	daoDTO "github.com/gobox-preegnees/connection_controller/internal/adapter/dao"
	entity "github.com/gobox-preegnees/connection_controller/internal/domain/entity"
	usecase "github.com/gobox-preegnees/connection_controller/internal/domain/usecase"

	"github.com/sirupsen/logrus"
)

type IOwnerDao interface {
	DeleteOneOwner(daoDTO.DeleteOneOwnerReqDTO) error
	// TODO: нужно что то возращать, чтобы было понятно создали или нет
	CreateOneOwner(daoDTO.CreateOneOwnerReqDTO) error
}

type ownerService struct {
	log *logrus.Logger
	dao IOwnerDao
}

type CnfOwnerService struct {
	Log *logrus.Logger
	OwnerDao IOwnerDao
}

func NewOwnerService(cnf CnfOwnerService) *ownerService {

	return &ownerService{
		log: cnf.Log,
		dao: cnf.OwnerDao,
	}
}

func (*ownerService) DeleteOwner(ctx context.Context, owner entity.Owner) (bool, error) {

	return false, nil
}

func (*ownerService) SaveOwner(ctx context.Context, owner entity.Owner) error {

	return nil
}

var _ usecase.IOwnerService = (*ownerService)(nil)
