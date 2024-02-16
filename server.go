package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/zufardhiyaulhaq/echo-redis/pkg/settings"

	redis_client "github.com/zufardhiyaulhaq/echo-redis/pkg/redis"
)

type Server struct {
	settings settings.Settings
	client   redis_client.RedisClient
}

func NewServer(settings settings.Settings, client redis_client.RedisClient) Server {
	return Server{
		settings: settings,
		client:   client,
	}
}

func (e Server) ServeHTTP() {
	handler := NewHandler(e.settings, e.client)

	r := mux.NewRouter()

	r.HandleFunc("/redis/{key}", handler.Handle)
	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello!"))
	})
	r.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello!"))
	})

	err := http.ListenAndServe(":"+e.settings.HTTPPort, r)
	if err != nil {
		log.Fatal().Err(err)
	}
}

type Handler struct {
	settings settings.Settings
	client   redis_client.RedisClient
}

func NewHandler(settings settings.Settings, client redis_client.RedisClient) Handler {
	return Handler{
		settings: settings,
		client:   client,
	}
}

func (h Handler) Handle(w http.ResponseWriter, req *http.Request) {
	key := uuid.NewString()
	value := mux.Vars(req)["key"]

	err := h.client.Set(key, value)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(err.Error()))
		return
	}

	data, err := h.client.Get(key)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(key + ":" + data))
}
