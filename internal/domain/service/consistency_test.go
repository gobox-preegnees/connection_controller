package service

import (
	"context"
	"errors"
	"testing"

	"github.com/gobox-preegnees/connection_controller/internal/domain/entity"
	"github.com/sirupsen/logrus"
)

func TestConsistensyService(t *testing.T) {

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetReportCaller(true)
	cService := NewConsistensyService(CnfConsistensyService{
		Log: logger,
	})

	inId := []string{"1", "2", "3"}

	initSlice := func() func(id string) []string {
		outid := make([]string, 0, len(inId))
		return func(id string) []string {
			if id == "" {
				return outid
			}
			outid = append(outid, id)
			return outid
		}
	}
	adder := initSlice()

	ctx := context.Background()

	done := make(chan struct{})

	go func() {
		for {
			c, err := cService.GetConsistency(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				t.Fatal(err)
			}
			adder(c.RequestId)
			done <- struct{}{}
		}
	}()

	go func() {
		for i := 0; i < len(inId); i++ {
			<-done
		}
		close(cService.consistensyCh)
		done<- struct{}{}
	}()

	for _, v := range inId {
		cService.SaveConsistency(ctx, entity.Consistency{
			RequestId: v,
		})
	}

	<-done
	t.Log(adder(""))
	if len(adder("")) != len(inId) {
		t.Fail()
	}
}
