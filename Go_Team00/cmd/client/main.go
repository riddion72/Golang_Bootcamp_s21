package main

import (
	"context"
	"flag"
	"fmt"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	loger "main/internal/logger"
	api "main/pkg/api/api/proto"
	cnf "main/pkg/config"
)

var (
	anomaliesDetected = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "client_anomalies_detected_total",
			Help: "Total number of detected anomalies",
		},
	)

	processingTime = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "client_message_processing_seconds",
			Help:    "Time spent processing messages",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1},
		},
	)

	connectionStatus = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "client_connection_status",
			Help: "Client connection status (1 = connected, 0 = disconnected)",
		},
	)

	currentMean = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "client_current_mean",
			Help: "Current mean value of frequencies",
		},
	)

	currentStdDev = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "client_current_stddev",
			Help: "Current standard deviation of frequencies",
		},
	)

	flagK float64 // Флаг k

	bufPool = sync.Pool{
		New: func() any {
			return new(api.Frequency)
		},
	}
)

func init() {
	flag.Float64Var(&flagK, "k", 1.0, "Value of anomaly coefficient")
	prometheus.MustRegister(
		anomaliesDetected,
		processingTime,
		connectionStatus,
		currentMean,
		currentStdDev,
	)
}

type Statistics struct {
	Mean   float64
	StdDev float64
	Count  int
	Sum    float64
	Sumsq  float64
}

func (s *Statistics) update(newValue float64) {
	s.Count++
	s.Sum += newValue
	s.Sumsq += newValue * newValue
	s.Mean = s.Sum / float64(s.Count)
	if s.Count > 1 {
		variance := (s.Sumsq / float64(s.Count)) - (s.Mean * s.Mean)
		if variance < 0 {
			variance = 0
		}
		s.StdDev = math.Sqrt(variance)
	}
}

func (s *Statistics) findAnomaly(value float64) bool {
	res := false
	if s.Count > 10 {
		lowerBound := s.Mean - flagK*s.StdDev
		upperBound := s.Mean + flagK*s.StdDev
		res = value < lowerBound || value > upperBound
	}
	return res
}

func process(setings *cnf.Config) {
	// Запускаем HTTP-сервер для метрик
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8083", nil)
	}()
	// Устанавливаем соединение
	conn, err := grpc.NewClient(setings.ServerHost+":"+setings.ServerPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		loger.WriteLog(fmt.Sprintf("Failed to connect to server: %v", err))
		connectionStatus.Set(0)
		return
	}
	defer conn.Close()

	connectionStatus.Set(1)
	loger.WriteLog("Процесс запущен")
	client := api.NewFrequencyServiseClient(conn)

	stream, err := client.GenerateFrequency(context.Background(), &api.Frequency{})
	if err != nil {
		loger.WriteLog(fmt.Sprintf("Error when calling a function GenerateFrequency: %v", err))
		return
	}

	stats := &Statistics{}
	entry := &api.Frequency{}

	for {
		start := time.Now()
		entry = bufPool.Get().(*api.Frequency)
		entry.Reset()
		entry, err = stream.Recv()
		if err != nil {
			if strings.Contains(err.Error(), "EOF") {
				loger.WriteLog(fmt.Sprintf("1Count: %v, mean: %.2f, stdDev: %.2f \n", stats.Count, stats.Mean, stats.StdDev))
				loger.WriteLog(fmt.Sprintf("Сервер пекратил работу: %v", err))
			} else {
				loger.WriteLog(fmt.Sprintf("2Count: %v, mean: %.2f, stdDev: %.2f \n", stats.Count, stats.Mean, stats.StdDev))
				loger.WriteLog(fmt.Sprintf("Ошибка при получении данных: %v", err))
			}
			return
		}

		// Обновляем статистику с новым значением частоты
		stats.update(entry.Frequency)

		// Обновляем статистику каждые 10 значений
		if stats.Count%10 == 0 {
			loger.WriteLog(fmt.Sprintf("Count: %v, mean: %v, stdDev: %v \n", stats.Count, stats.Mean, stats.StdDev))
		}

		currentMean.Set(stats.Mean)
		currentStdDev.Set(stats.StdDev)

		// Проверка на аномалию, если количество значений больше 10
		if stats.Count > 10 {
			if stats.findAnomaly(entry.Frequency) {
				anomaliesDetected.Inc()
				processingTime.Observe(time.Since(start).Seconds())
				loger.WriteLog(fmt.Sprintf("Session ID: %s, Frequency: %f, Timestamp: %d\n", entry.SessionId, entry.Frequency, entry.Timestamp))
			}
		}

		bufPool.Put(entry)
		time.Sleep(time.Second) // Пауза
	}
}

func main() {

	flag.Parse()

	if loger.PrepareLogger() != nil {
		return
	}

	defer func() {
		loger.WriteLog("Клиент прекращает работу")
		if err := recover(); err != nil {
			loger.WriteLog(fmt.Sprint("Unknown panic happend: ", err))
		}
		loger.CloseLogger()
	}()

	setings, err := cnf.LoadConfig()
	if err != nil {
		loger.WriteLog(fmt.Sprintf("Error Config: %s", err))
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		process(setings)
		wg.Done()
	}()
	wg.Wait()
}
