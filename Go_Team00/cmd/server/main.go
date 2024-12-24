package main

import (
	"log"
	"math/rand"
	"net"
	"time"

	api "main/pkg/api/api/proto"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	api.UnimplementedFrequencyServiseServer
}

func (s *server) GenerateFrequency(req *api.Frequency, stream api.FrequencyServise_GenerateFrequencyServer) error {
	// Генерируем уникальный идентификатор для каждой сессии
	uuid := uuid.New().String()

	// Генерируем случайное среднее значение и стандартное отклонение
	// (для генерации нормального распределения частот)
	mean := rand.Float64()*20 - 10
	stdDev := rand.Float64()*1.2 + 0.3

	// Логируем сгенерированные значения
	log.Printf("Generating frequency mean: %.2f, stddev: %.2f\n", mean, stdDev)

	// Генерация частоты и отправка данных клиенту в потоке
	for {
		frequency := rand.NormFloat64()*stdDev + mean

		entry := &api.Frequency{
			SessionId: uuid,
			Frequency: frequency,
			Timestamp: time.Now().Unix(),
		}

		if err := stream.Send(entry); err != nil {
			return err
		}
		// log.Printf("%.2f\n", frequency)
		time.Sleep(time.Second)
	}

}

func main() {
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal("Faeled create server listener: ", err)
	}

	serv := grpc.NewServer()

	api.RegisterFrequencyServiseServer(serv, &server{})
	log.Printf("server listening at %v", listener.Addr())
	if err := serv.Serve(listener); err != nil {
		log.Printf("failed to serve: %v", err)
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
