package main

import (
	"context"
	"flag"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"

	loger "main/internal/logger"
	api "main/pkg/api/api/proto"
	cnf "main/pkg/config"
	model "main/pkg/model"
	rep "main/pkg/repository"
)

var (
	flagK       float64 // Флаг k
	logFileName string  = "assets/statistics.log"

	bufPool = sync.Pool{
		New: func() any {
			return new(api.Frequency)
		},
	}
)

func init() {
	flag.Float64Var(&flagK, "k", 1.0, "Value of anomaly coefficient")
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

func pushAnomaly(db *gorm.DB, entry *api.Frequency) error {
	// Сохранение аномалии в базу данных
	anomaly := model.Anomalies{
		SessionID: entry.SessionId,
		Frequency: entry.Frequency,
		Timestamp: time.Unix(entry.Timestamp, 0),
	}

	result := db.Create(&anomaly)
	if result.Error != nil {
		loger.WriteLog(fmt.Sprintf("Ошибка при сохранении аномалии в БД: %v", result.Error))
		return result.Error
	} else {
		loger.WriteLog(fmt.Sprintf("Аномалия успешно сохранена в ID: %d", anomaly.ID))
		return nil
	}
}

func process(setings *cnf.Config, db *gorm.DB) {
	// Устанавливаем соединение
	conn, err := grpc.NewClient(setings.ServerHost+":"+setings.ServerPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		loger.WriteLog(fmt.Sprintf("Failed to connect to server: %v", err))
		return
	}
	defer conn.Close()

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

		// Вывод полученного значения
		// fmt.Printf("Session ID: %s, Frequency: %f, Timestamp: %d\n", entry.SessionId, entry.Frequency, entry.Timestamp)

		// Обновляем статистику с новым значением частоты
		stats.update(entry.Frequency)

		// Обновляем статистику каждые 10 значений
		if stats.Count%10 == 0 {
			loger.WriteLog(fmt.Sprintf("Count: %v, mean: %v, stdDev: %v \n", stats.Count, stats.Mean, stats.StdDev))
		}

		// Проверка на аномалию, если количество значений больше 10
		if stats.Count > 10 {
			if stats.findAnomaly(entry.Frequency) {
				if pushAnomaly(db, entry) != nil {
					loger.WriteLog(fmt.Sprintf("Ошибка при попытке сохранить аномалию в БД: %v", err))
					return
				}
			}
		}

		bufPool.Put(entry)
		time.Sleep(time.Second) // Пауза
	}
}

func main() {

	flag.Parse()

	if loger.PrepareLogger(logFileName) != nil {
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

	db, err := rep.PrepareDB(setings)
	if err != nil {
		loger.WriteLog(fmt.Sprintf("Error DB: %s", err))
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		process(setings, db) // Передаем db в Process
		wg.Done()
	}()
	wg.Wait()
}
