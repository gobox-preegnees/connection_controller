package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	entity "github.com/gobox-preegnees/connection_controller/internal/domain/entity"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"github.com/r3labs/sse/v2"
	"github.com/sirupsen/logrus"
)

type IStatisticService interface {
	GetAverageTimeout() (float64, error)
}

type ISnapshotService interface {
	SaveSnapshot(ctx context.Context, snapshot entity.Snapshot) error
}

type IConsistencyService interface {
	Read() entity.Consistency
}

type IOwnerService interface {
	SaveOwner(ctx context.Context, owner entity.Owner) error
	DeleteOwner(ctx context.Context, owner entity.Owner)
}

type http1 struct {
	log                *logrus.Logger
	addr               string
	snapshotService    ISnapshotService
	consistencyService IConsistencyService
	ownerService       IOwnerService
	alg                string
	secret             string
}

type CnfhttpServer struct {
	Log                *logrus.Logger
	Addr               string
	SnapshotService    ISnapshotService
	ConsistencyService IConsistencyService
	OwnerService       IOwnerService
	JWTAlg             string
	JWTSecret          string
}

func NewhttpServer(cnf CnfhttpServer) *http1 {

	return &http1{
		log:                cnf.Log,
		addr:               cnf.Addr,
		snapshotService:    cnf.SnapshotService,
		consistencyService: cnf.ConsistencyService,
		ownerService:       cnf.OwnerService,
		alg:                cnf.JWTAlg,
		secret:             cnf.JWTSecret,
	}
}

func (h http1) Run() error {

	r := h.router()
	// TODO: add tls/ssl
	return http.ListenAndServe(h.addr, r)
}

func (h http1) router() http.Handler {

	r := chi.NewRouter()
	sseServer := sse.New()

	go func() {
		for {
			// TODO: добавить статистику по id или таймстемпам
			consistency := h.consistencyService.Read()
			streamId := fmt.Sprintf("%s_%s", consistency.Username, consistency.Folder)
			data, err := json.Marshal(consistency)
			if err != nil {
				h.log.Fatal(err)
			}
			h.log.Debugf("consistency:%s, streamId:%s, data:%s", consistency, streamId, data)

			sseServer.Publish(streamId, &sse.Event{
				Data: data,
			})
		}
	}()

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(jwtauth.New(h.alg, []byte(h.secret), nil)))
		r.Use(jwtauth.Authenticator)

		// TODO: добавить ручку для получения статистики по времени ответа
		r.Post("/snapshot", func(w http.ResponseWriter, r *http.Request) {
			shapshot := entity.Snapshot{}
			err := json.NewDecoder(r.Body).Decode(&shapshot)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			v := validator.New()
			if err := v.Struct(shapshot); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			ctx := context.Background()
			if err := h.snapshotService.SaveSnapshot(ctx, shapshot); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})

		r.Post("/event", func(w http.ResponseWriter, r *http.Request) {
			owner := entity.Owner{}
			err := json.NewDecoder(r.Body).Decode(&owner)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			v := validator.New()
			if err := v.Struct(owner); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			ctx := context.Background()
			if err := h.ownerService.SaveOwner(ctx, owner); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			go func() {
				_owner := owner
				streamId := fmt.Sprintf("%s_%s", _owner.Username, _owner.Folder)
				sseServer.CreateStream(streamId)
				<-r.Context().Done()
				h.log.Info("Client is disconnected")
				sseServer.RemoveStream(streamId)
				h.ownerService.DeleteOwner(ctx, _owner)
				return
			}()
			sseServer.ServeHTTP(w, r)
		})
	})
	return r
}
