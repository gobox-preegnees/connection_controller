package service

import (
	"context"

	daoDTO "github.com/gobox-preegnees/connection_controller/internal/adapter/dao"
	entity "github.com/gobox-preegnees/connection_controller/internal/domain/entity"
	usecase "github.com/gobox-preegnees/connection_controller/internal/domain/usecase"

	"github.com/sirupsen/logrus"
)

type IStreamDao interface {
	DeleteOneStream(req daoDTO.DeleteOneStreamReqDTO) (numberOfUsers int, err error)
	CreateOneStream(req daoDTO.CreateOneStreamReqDTO) (numberOfUsers int, err error)
}

type streamService struct {
	log *logrus.Logger
	dao IStreamDao
}

type CnfStreamService struct {
	Log      *logrus.Logger
	StreamDao IStreamDao
}

func NewStreamService(cnf CnfStreamService) *streamService {

	return &streamService{
		log: cnf.Log,
		dao: cnf.StreamDao,
	}
}

func (o streamService) DeleteStream(ctx context.Context, stream entity.Stream) (int, error) {

	return o.dao.DeleteOneStream(daoDTO.DeleteOneStreamReqDTO{
		Ctx:      ctx,
		StreamId: stream.StreamId,
	})
}

func (o streamService) SaveStream(ctx context.Context, stream entity.Stream) (int, error) {

	return o.dao.CreateOneStream(daoDTO.CreateOneStreamReqDTO{
		Ctx:       ctx,
		StreamId: stream.StreamId,
	})
}

var _ usecase.IStreamService = (*streamService)(nil)
