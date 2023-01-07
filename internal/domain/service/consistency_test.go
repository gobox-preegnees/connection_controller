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
		done <- struct{}{}
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

func TestCtxCancel(t *testing.T) {

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetReportCaller(true)
	cService := NewConsistensyService(CnfConsistensyService{
		Log: logger,
	})

	okRecv := false
	okSend := false

	in := []string{"1", "2", "3", "4", "5"}
	out := make([]string, 0, len(in)-1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		for {
			val, err := cService.GetConsistency(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					okRecv = true
					return 
				}
			}
			out = append(out, val.RequestId)
		}
	}()

	for i, v := range in {
		if err := cService.SaveConsistency(ctx, entity.Consistency{
			RequestId: v,
		}); err != nil {
			if errors.Is(err, context.Canceled) {
				okSend = true
			}
		}
		if i == 1 {
			cancel()
		}
	}

	if okRecv != okSend {
		t.Fatal("okRecv != okSend && len(in) - len(out) != 1")
	}
}
