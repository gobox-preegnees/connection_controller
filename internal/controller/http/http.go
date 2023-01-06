package http

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	entity "github.com/gobox-preegnees/connection_controller/internal/domain/entity"
	errs "github.com/gobox-preegnees/connection_controller/internal/errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/r3labs/sse/v2"
	"github.com/sirupsen/logrus"
)

type IUsecase interface {
	SaveSnapshot(ctx context.Context, snapshot entity.Snapshot) (err error)
	GetConsistency(ctx context.Context) (consistency entity.Consistency, err error)
	SaveOwner(ctx context.Context, owner entity.Owner) (err error)
	DeleteOwner(ctx context.Context, owner entity.Owner) (err error)
}

type http1 struct {
	ctx           context.Context
	log           *logrus.Logger
	addr          string
	usecase       IUsecase
	alg           string
	secret        string
	server        *http.Server
	crtPath       string
	kayPath       string
	cancelTimeout int
}

type CnfhttpServer struct {
	Ctx       context.Context
	Log       *logrus.Logger
	Addr      string
	Usecase   IUsecase
	JWTAlg    string
	JWTSecret string
	CrtPath   string
	KeyPath   string
}

func NewhttpServer(cnf CnfhttpServer) *http1 {

	return &http1{
		ctx:           cnf.Ctx,
		log:           cnf.Log,
		addr:          cnf.Addr,
		usecase:       cnf.Usecase,
		alg:           cnf.JWTAlg,
		secret:        cnf.JWTSecret,
		crtPath:       cnf.CrtPath,
		kayPath:       cnf.KeyPath,
		cancelTimeout: 20,
	}
}

func (h *http1) Run() {

	cer, err := tls.LoadX509KeyPair(h.crtPath, h.kayPath)
	if err != nil {
		h.log.Fatal(err)
	}

	h.server = &http.Server{
		Addr:    h.addr,
		Handler: h.router(),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cer},
		},
	}
	h.log.Info("server starting...")

	h.log.Fatal(h.server.ListenAndServe())
}

func (h *http1) Shutdown() error {

	return h.server.Shutdown(h.ctx)
}

func (h http1) router() http.Handler {

	r := chi.NewRouter()
	sseServer := sse.New()

	go func() {
		for {
			// TODO: добавить статистику по id или таймстемпам
			consistency, err := h.usecase.GetConsistency(h.ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				if err != nil {
					h.log.Fatal(err)
				}
			}

			streamId := consistency.StreamId
			requestId := consistency.RequestId
			sseServer.Publish(streamId, &sse.Event{
				Data: consistency.Data,
			})
			h.log.Debugf("consistency sent to streamId: %s, requestId: %s", streamId, requestId)
		}
	}()

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(jwtauth.New(h.alg, []byte(h.secret), nil)))
		r.Use(jwtauth.Authenticator)

		// TODO: добавить ручку для получения статистики по времени ответа
		r.Post("/snapshot", func(w http.ResponseWriter, r *http.Request) {
			snapshot := entity.Snapshot{}
			err := json.NewDecoder(r.Body).Decode(&snapshot)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			snapshot.RequestId = uuid.NewString()
			snapshot.StreamId = fmt.Sprintf("%s_%s", snapshot.Username, snapshot.Folder)
			snapshot.Timestamp = time.Now().UTC().Unix()

			v := validator.New()
			if err := v.Struct(snapshot); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(h.cancelTimeout)*time.Second)
			defer cancel()
			if err := h.usecase.SaveSnapshot(ctx, snapshot); err != nil {
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

			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(h.cancelTimeout)*time.Second)
			defer cancel()
			err = h.usecase.SaveOwner(ctx, owner)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			streamId := fmt.Sprintf("%s_%s", owner.Username, owner.Folder)
			h.log.Debugf("owener:%v connected to sse-server, streamid:%s", owner, streamId)


			go func() {
				<-r.Context().Done()
				h.log.Debugf("client: %s is disconnected", owner.Client)

				if err := h.usecase.DeleteOwner(ctx, owner); err != nil {
					if errors.Is(err, errs.ErrNoVisitors) {
						sseServer.RemoveStream(streamId)
					} else {
						h.log.Errorf("unable delete stream:%s", streamId)
					}
				}
				return
			}()
			sseServer.ServeHTTP(w, r)
		})
	})
	return r
}
