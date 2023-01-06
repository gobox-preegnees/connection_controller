package service

import (
	"context"

	daoDTO "github.com/gobox-preegnees/connection_controller/internal/adapter/dao"
	entity "github.com/gobox-preegnees/connection_controller/internal/domain/entity"
	usecase "github.com/gobox-preegnees/connection_controller/internal/domain/usecase"

	"github.com/sirupsen/logrus"
)

type IOwnerDao interface {
	DeleteOneOwner(req daoDTO.DeleteOneOwnerReqDTO) (countOnwers int, err error)
	CreateOneOwner(req daoDTO.CreateOneOwnerReqDTO) (countOnwers int, err error)
}

type ownerService struct {
	log *logrus.Logger
	dao IOwnerDao
}

type CnfOwnerService struct {
	Log      *logrus.Logger
	OwnerDao IOwnerDao
}

func NewOwnerService(cnf CnfOwnerService) *ownerService {

	return &ownerService{
		log: cnf.Log,
		dao: cnf.OwnerDao,
	}
}

func (o ownerService) DeleteOwner(ctx context.Context, owner entity.Owner) (int, error) {

	return o.dao.DeleteOneOwner(daoDTO.DeleteOneOwnerReqDTO{
		Ctx:       ctx,
		Usernamme: owner.Username,
		Folder:    owner.Folder,
	})
}

func (o ownerService) SaveOwner(ctx context.Context, owner entity.Owner) (int, error) {

	return o.dao.CreateOneOwner(daoDTO.CreateOneOwnerReqDTO{
		Ctx:       ctx,
		Usernamme: owner.Username,
		Folder:    owner.Folder,
	})
}

var _ usecase.IOwnerService = (*ownerService)(nil)
