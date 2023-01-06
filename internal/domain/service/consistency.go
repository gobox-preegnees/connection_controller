package service

import (
	entity "github.com/gobox-preegnees/connection_controller/internal/domain/entity"
	usecase "github.com/gobox-preegnees/connection_controller/internal/domain/usecase"

	"github.com/sirupsen/logrus"
)

type consistensyService struct {
	log           *logrus.Logger
	consistensyCh chan entity.Consistency
}

type CnfConsistensyService struct {
	Log *logrus.Logger
}

func New(cnf CnfConsistensyService) *consistensyService {

	return &consistensyService{
		log:           cnf.Log,
		consistensyCh: make(chan entity.Consistency),
	}
}

func (c consistensyService) GetConsistency() entity.Consistency {

	return <-c.consistensyCh
}

func (c consistensyService) SaveConsistency(consistency entity.Consistency) error {

	// TODO: можно сделать таймаут через контекста или еще что
	c.consistensyCh <- consistency
	return nil
}

var _ usecase.IConsistencyService = (*consistensyService)(nil)
