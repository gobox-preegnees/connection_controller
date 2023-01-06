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

type IUsecase interface {
	SaveSnapshot(ctx context.Context, snapshot entity.Snapshot) error
	GetConsistency() entity.Consistency
	SaveOwner(ctx context.Context, owner entity.Owner) error
	DeleteOwner(ctx context.Context, owner entity.Owner) error
}

type http1 struct {
	log     *logrus.Logger
	addr    string
	usecase IUsecase
	alg     string
	secret  string
}

type CnfhttpServer struct {
	Log       *logrus.Logger
	Addr      string
	Usecase   IUsecase
	JWTAlg    string
	JWTSecret string
}

func NewhttpServer(cnf CnfhttpServer) *http1 {

	return &http1{
		log:     cnf.Log,
		addr:    cnf.Addr,
		usecase: cnf.Usecase,
		alg:     cnf.JWTAlg,
		secret:  cnf.JWTSecret,
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
			consistency := h.usecase.GetConsistency()
			streamId := fmt.Sprintf("%s_%s", consistency.Username, consistency.Folder)
			data, err := json.Marshal(consistency)
			if err != nil {
				h.log.Fatal(err)
			}
			h.log.Debugf("consistency:%s, streamId:%s, data:%s, timestamp: %d", consistency, streamId, data, consistency.Timestamp)

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
			if err := h.usecase.SaveSnapshot(ctx, shapshot); err != nil {
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
			if err := h.usecase.SaveOwner(ctx, owner); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			go func() {
				streamId := fmt.Sprintf("%s_%s", owner.Username, owner.Folder)
				sseServer.CreateStream(streamId)
				<-r.Context().Done()
				h.log.Info("Client is disconnected")
				sseServer.RemoveStream(streamId)
				h.usecase.DeleteOwner(ctx, owner)
				return
			}()
			sseServer.ServeHTTP(w, r)
		})
	})
	return r
}
