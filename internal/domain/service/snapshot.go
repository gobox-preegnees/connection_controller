package service

import (
	"context"
	"encoding/json"
	"time"

	mbDTO "github.com/gobox-preegnees/connection_controller/internal/adapter/message_broker"
	entity "github.com/gobox-preegnees/connection_controller/internal/domain/entity"
	usecase "github.com/gobox-preegnees/connection_controller/internal/domain/usecase"

	"github.com/sirupsen/logrus"
)

//go:generate mockgen -destination=../../mocks/domain/service/snapshot/ISnapshotMessageBroker/ISnapshotMessageBroker.go -source=snapshot.go
type ISnapshotMessageBroker interface {
	SaveSnapshot(mbDTO.PublishSnapshotReqDTO) error
}

// snapshotService.
type snapshotService struct {
	log *logrus.Logger
	mb ISnapshotMessageBroker
}

// CnfSnapshotService.
type CnfSnapshotService struct {
	Log *logrus.Logger
	SnapshotMessageBroker ISnapshotMessageBroker
}

// NewShanpshotService.
func NewShanpshotService(cnf CnfSnapshotService) *snapshotService {

	return &snapshotService{
		log: cnf.Log,
		mb: cnf.SnapshotMessageBroker,
	}
}

// SaveSnapshot. 
func (s snapshotService) SaveSnapshot(ctx context.Context, snapshot entity.Snapshot) error {

	data, err := json.Marshal(snapshot)
	if err!= nil {
        return err
    }
	
	return s.mb.SaveSnapshot(mbDTO.PublishSnapshotReqDTO{
		RequestId: snapshot.RequestId,
		StreamId: snapshot.StreamId,
		Timestamp: time.Unix(snapshot.Timestamp, 0).UTC(),
		Data: data,
	})
}

var _ usecase.ISnapshotService = (*snapshotService)(nil)
