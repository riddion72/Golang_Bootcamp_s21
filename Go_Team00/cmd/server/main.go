package main

import (
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"

	api "main/pkg/api/api/proto"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

// Prometheus metrics
var (
	sessionsActive = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "frequency_sessions_active",
			Help: "Number of active frequency generation sessions",
		},
	)

	messagesSent = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "frequency_messages_sent_total",
			Help: "Total number of frequency messages sent",
		},
	)

	streamErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "frequency_stream_errors_total",
			Help: "Total number of stream errors",
		},
	)

	bufferPoolUsage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "buffer_pool_usage",
			Help: "Current buffer pool usage",
		})
	grpcRequestsDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_request_duration_seconds",
			Help:    "gRPC request duration distribution",
			Buckets: []float64{0.1, 0.3, 0.5, 1, 2},
		},
		[]string{"method"},
	)

	grpcRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_total",
			Help: "Total number of gRPC requests",
		},
		[]string{"method", "status"},
	)
)

func init() {
	prometheus.MustRegister(
		sessionsActive,
		messagesSent,
		streamErrors,
		bufferPoolUsage,
		grpcRequestsDuration,
		grpcRequestsTotal,
	)
	prometheus.MustRegister(prometheus.NewProcessCollector(
		prometheus.ProcessCollectorOpts{
			Namespace: "frequency",
		},
	))
}

// server is used to implement helloworld.GreeterServer.
type server struct {
	api.UnimplementedFrequencyServiseServer
}

func (s *server) GenerateFrequency(req *api.Frequency, stream api.FrequencyServise_GenerateFrequencyServer) error {
	startTime := time.Now()
	methodName := "GenerateFrequency"
	grpcRequestsTotal.WithLabelValues(methodName, "started").Inc()
	defer func() {
		duration := time.Since(startTime).Seconds()
		grpcRequestsDuration.WithLabelValues(methodName).Observe(duration)
	}()
	sessionsActive.Inc()
	defer sessionsActive.Dec()
	uuid := uuid.New().String()
	mean := rand.Float64()*20 - 10
	stdDev := rand.Float64()*1.2 + 0.3

	log.Printf("New session: %s, mean: %.2f, stddev: %.2f", uuid, mean, stdDev)

	for {
		frequency := rand.NormFloat64()*stdDev + mean

		entry := &api.Frequency{
			SessionId: uuid,
			Frequency: frequency,
			Timestamp: time.Now().Unix(),
		}

		if err := stream.Send(entry); err != nil {
			grpcRequestsTotal.WithLabelValues(methodName, "error").Inc()
			streamErrors.Inc()
			return err
		}
		grpcRequestsTotal.WithLabelValues(methodName, "success").Inc()
		messagesSent.Inc()
		time.Sleep(time.Second)
	}
}

func main() {
	// Start metrics server
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Metrics server started on :8082")
		if err := http.ListenAndServe(":8082", nil); err != nil {
			log.Printf("Failed to start metrics server: %v", err)
		}
	}()

	// gRPC server
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal("Failed to create server listener: ", err)
	}

	serv := grpc.NewServer()
	api.RegisterFrequencyServiseServer(serv, &server{})

	log.Printf("gRPC server listening at %v", listener.Addr())
	if err := serv.Serve(listener); err != nil {
		log.Printf("Failed to serve: %v", err)
		serv.GracefulStop()
	}

	defer func() {
		listener.Close()
		// serv.GracefulStop()
		log.Println("server stopped")
		if err := recover(); err != nil {
			log.Println("Unknown panic happend: ", err)
		}
	}()

	// sigChan := make(chan os.Signal, 1)
	// signal.Notify(sigChan, os.Interrupt)
	// sig := <-sigChan
	// log.Println()
	// log.Println("Assept signal: ", sig)
	// serv.GracefulStop()
	// listener.Close()
	// log.Println("exiting")
}
