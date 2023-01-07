package service

import (
	"context"
	"sync"
	"time"

	entity "github.com/gobox-preegnees/connection_controller/internal/domain/entity"
	usecase "github.com/gobox-preegnees/connection_controller/internal/domain/usecase"

	"github.com/sirupsen/logrus"
)

// consistensyService.
type consistensyService struct {
	isStopped     bool
	doneCh        chan struct{}
	once          sync.Once
	log           *logrus.Logger
	consistensyCh chan entity.Consistency
}

// CnfConsistensyService.
type CnfConsistensyService struct {
	Log *logrus.Logger
}

// NewConsistensyService.
func NewConsistensyService(cnf CnfConsistensyService) *consistensyService {

	return &consistensyService{
		log:           cnf.Log,
		consistensyCh: make(chan entity.Consistency),
		doneCh:        make(chan struct{}, 2),
	}
}

// GetConsistency. Это метод возвращает данные, которые приходят посредством метода SaveConsistency
func (c *consistensyService) GetConsistency(ctx context.Context) (entity.Consistency, error) {

	c.once.Do(func() {
		c.isStopped = false
		go func() {
			select {
			case <-ctx.Done():
				c.isStopped = true
				time.Sleep(1 * time.Second)
				
				// Чтобы избежать паники при закрытии конекста
				if len(c.doneCh) < 2 {
					return
				}
				close(c.consistensyCh)
				close(c.doneCh)
			}
		}()
	})

	// Если контекст закрылся
	if c.isStopped {
		c.doneCh <- struct{}{}
		return entity.Consistency{}, context.Canceled
	}

	consistensy, ok := <-c.consistensyCh
	if !ok {
		return entity.Consistency{}, context.Canceled
	}
	return consistensy, nil
}

// SaveConsistenc. Это метод сохраняет данные из брокера сообщений в канал (на http ответ)
func (c *consistensyService) SaveConsistency(ctx context.Context, consistency entity.Consistency) error {

	if c.isStopped {
		c.doneCh <- struct{}{}
		return context.Canceled
	}
	c.consistensyCh <- consistency
	return nil
}

var _ usecase.IConsistencyService = (*consistensyService)(nil)
