package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterBusRoutes(r *mux.Router, b *Broker) {
	s := &BusServer{Broker: b}
	api := r.PathPrefix("/v1/bus").Subrouter()

	// wrap bus endpoints with JWT middleware
	api.Handle("/publish", JWTAuthMiddleware(http.HandlerFunc(s.handlePublish))).Methods("POST")
	api.Handle("/subscribe", JWTAuthMiddleware(http.HandlerFunc(s.handleWSSubscribe)))
}
