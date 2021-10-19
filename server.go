package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/tidwall/evio"
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

func (e Server) ServeEcho() {
	var events evio.Events

	events.Data = func(c evio.Conn, in []byte) (out []byte, action evio.Action) {
		key := uuid.NewString()
		value := string(in)

		err := e.client.Set(key, value)
		if err != nil {
			out = []byte(err.Error())
			return
		}

		data, err := e.client.Get(key)
		if err != nil {
			out = []byte(err.Error())
			return
		}

		out = []byte(key + ":" + data)

		return
	}

	if err := evio.Serve(events, "tcp://0.0.0.0:"+e.settings.RedisEventPort); err != nil {
		panic(err.Error())
	}
}

func (e Server) ServeHTTP() {
	handler := NewHandler(e.settings, e.client)

	r := mux.NewRouter()
	r.HandleFunc("/redis/{key}", handler.Handle)
	http.ListenAndServe("0.0.0.0:80", r)
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
