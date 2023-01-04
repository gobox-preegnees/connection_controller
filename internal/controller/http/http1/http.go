package http1

import (
	"net/http"
	"context"
	"encoding/json"

	entity "github.com/gobox-preegnees/connection_controller/internal/domain/entity"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/sirupsen/logrus"
	"github.com/r3labs/sse/v2"
	"github.com/go-playground/validator/v10"
)

type ISnapshotService interface {
	SaveSnapshot(ctx context.Context, snapshot entity.Snapshot) error
	SaveOwner(ctx context.Context, snapshot entity.Snapshot) error
}

type IConsistencyService interface {
}

type http1 struct {
	log                *logrus.Logger
	addr               string
	snapshotService    ISnapshotService
	consistencyService IConsistencyService
	alg                string
	secret             string
}

type CnfHttp1Server struct {
	Log                *logrus.Logger
	Addr               string
	SnapshotService    ISnapshotService
	ConsistencyService IConsistencyService
	JWTAlg             string
	JWTSecret          string
}

func NewHttp1Server(cnf CnfHttp1Server) *http1 {

	return &http1{
		log:                cnf.Log,
		addr:               cnf.Addr,
		snapshotService:    cnf.SnapshotService,
		consistencyService: cnf.ConsistencyService,
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

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(jwtauth.New(h.alg, []byte(h.secret), nil)))
		r.Use(jwtauth.Authenticator)
		
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
			
			go func() {
				<-r.Context().Done()
				h.log.Info("Client is disconnected")
				return
			}()
			sseServer.ServeHTTP(w, r)
		})
	})
	return r
}
