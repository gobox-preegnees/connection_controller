package service

import (
	"context"
	"encoding/json"

	daoDTO "github.com/gobox-preegnees/connection_controller/internal/adapter/dao"
	entity "github.com/gobox-preegnees/connection_controller/internal/domain/entity"
	usecase "github.com/gobox-preegnees/connection_controller/internal/domain/usecase"

	"github.com/sirupsen/logrus"
)

type ISnapshotDao interface {
	CreateOneSnapshot(daoDTO.CreateOneSnapshotReqDTO) error
}

type snapshotService struct {
	log *logrus.Logger
	dao ISnapshotDao
}

type CnfSnapshotService struct {
	Log *logrus.Logger
	SnapshotDao ISnapshotDao
}

func NewShanpshotService(cnf CnfSnapshotService) *snapshotService {

	return &snapshotService{
		log: cnf.Log,
		dao: cnf.SnapshotDao,
	}
}

// TODO: тут придумать что нибудь по лучше, чем перекидывание байт
func (s snapshotService) SaveSnapshot(ctx context.Context, snapshot entity.Snapshot) error {

	data, err := json.Marshal(snapshot)
	if err!= nil {
        return err
    }
	return s.dao.CreateOneSnapshot(daoDTO.CreateOneSnapshotReqDTO{
		Snapshot: data,
	})
}

var _ usecase.ISnapshotService = (*snapshotService)(nil)
