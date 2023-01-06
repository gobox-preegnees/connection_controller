package service

import (
	"context"

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

func (c consistensyService) GetConsistency(ctx context.Context) (entity.Consistency, error) {

	return <-c.consistensyCh, nil
}

func (c consistensyService) SaveConsistency(ctx context.Context, consistency entity.Consistency) error {

	c.consistensyCh <- consistency
	return nil
}

var _ usecase.IConsistencyService = (*consistensyService)(nil)
