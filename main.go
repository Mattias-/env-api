package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gorilla/handlers"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	version = "v0.0.0"

	infoMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "env_api",
			Name:      "info",
			Help:      "Information about the env-api service.",
		},
		[]string{
			// Which version is running?
			"version",
		},
	)
)

func jsonEnv(w http.ResponseWriter, req *http.Request) {
	datas := map[string]string{}
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		datas[pair[0]] = pair[1]
	}
	d, _ := json.Marshal(datas)
	w.Header().Set("Content-Type", "application/json")
	w.Write(d)
}

func health(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
}

func versionEndpoint(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, version)
}

func main() {
	log.Println("Starting env-api application...")

	http.HandleFunc("/", jsonEnv)
	http.HandleFunc("/health", health)
	http.HandleFunc("/version", versionEndpoint)
	http.Handle("/metrics", promhttp.Handler())

	prometheus.MustRegister(infoMetric)
	infoMetric.WithLabelValues(version).Set(1)

	s := http.Server{
		Addr:    ":8080",
		Handler: handlers.LoggingHandler(os.Stdout, http.DefaultServeMux)}
	go func() {
		log.Fatal(s.ListenAndServe())
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Println("Shutdown signal received, exiting...")

	s.Shutdown(context.Background())
}
