package service

import (
	"context"
	"sync"

	entity "github.com/gobox-preegnees/connection_controller/internal/domain/entity"
	usecase "github.com/gobox-preegnees/connection_controller/internal/domain/usecase"

	"github.com/sirupsen/logrus"
)

type consistensyService struct {
	isStopped     bool
	once          sync.Once
	log           *logrus.Logger
	consistensyCh chan entity.Consistency
}

type CnfConsistensyService struct {
	Log *logrus.Logger
}

func NewConsistensyService(cnf CnfConsistensyService) *consistensyService {

	return &consistensyService{
		log:           cnf.Log,
		consistensyCh: make(chan entity.Consistency),
	}
}

func (c *consistensyService) GetConsistency(ctx context.Context) (entity.Consistency, error) {

	c.once.Do(func() {
		c.isStopped = false
		go func() {
			select {
			case <-ctx.Done():
				c.isStopped = true
			}
		}()
	})

	if c.isStopped {
		return entity.Consistency{}, context.Canceled
	}

	consistensy, ok := <-c.consistensyCh
	if !ok {
		return entity.Consistency{}, context.Canceled
	}
	return consistensy, nil
}

func (c *consistensyService) SaveConsistency(ctx context.Context, consistency entity.Consistency) error {

	if c.isStopped {
		return context.Canceled
	}
	c.consistensyCh <- consistency
	return nil
}

var _ usecase.IConsistencyService = (*consistensyService)(nil)
