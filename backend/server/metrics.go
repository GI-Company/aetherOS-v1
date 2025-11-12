package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	messagesPublished = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "aether_messages_published_total",
			Help: "Total number of messages published to the broker.",
		},
		[]string{"topic"},
	)
	messagesDropped = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "aether_messages_dropped_total",
			Help: "Total number of messages dropped by the broker.",
		},
		[]string{"topic"},
	)
	activeSubscribers = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aether_active_subscribers",
			Help: "Number of active subscribers.",
		},
		[]string{"topic"},
	)
)

func RegisterMetricsRoutes(r *mux.Router) {
	r.Handle("/metrics", promhttp.Handler())
}
