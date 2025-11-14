package server

import (
	"github.com/gorilla/mux"
)

func OldRegisterBusRoutes(r *mux.Router, b *Broker) {
	s := &BusServer{Broker: b}
	api := r.PathPrefix("/v1/bus").Subrouter()

	api.HandleFunc("/publish", s.handlePublish).Methods("POST")
	api.HandleFunc("/subscribe", s.handleWSSubscribe)
}
