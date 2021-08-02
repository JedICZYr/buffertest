package one

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)


type PrometheusServer struct {
	router *mux.Router
	msgRecevied prometheus.Counter
	msgProcessed prometheus.Counter
	msgWritten prometheus.Counter
	opsDelay prometheus.Gauge
}

func NewPrometheusServer() *PrometheusServer {
	router := mux.NewRouter()

	// Prometheus endpoint
	router.Path("/metrics").Handler(promhttp.Handler())

	fmt.Println("Serving requests on port 9000")
	go func() {
		err := http.ListenAndServe(":9000", router)
		log.Fatal(err)
	}()

	msgReceived := promauto.NewCounter(prometheus.CounterOpts{
		Name: "transformer_recevied_msg_total",
		Help: "The total number of received messages",
	})

	msgProcessed := promauto.NewCounter(prometheus.CounterOpts{
		Name: "transformer_processed_msg_total",
		Help: "The total number of processed messages",
	})

	msgWritten := promauto.NewCounter(prometheus.CounterOpts{
		Name: "transformer_written_msg_total",
		Help: "The total number of written messages",
	})

	opsDelay := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "transformer_ops_delay",
		Help:      "Maximum delay between receive and write.",
	})

	prometheus.MustRegister(opsDelay)

	return &PrometheusServer{router, msgReceived, msgProcessed, msgWritten, opsDelay}
}


