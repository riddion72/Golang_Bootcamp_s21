package main

import (
	"fmt"
	"log"
	"net/http"

	"main/internal/limits"
	api "main/internal/transport"
	"main/pkg/postgres"
)

func init() {
	postgres.MigrateDB()
}

func main() {
	defer func() {
		log.Println("server stopped")
		if err := recover(); err != nil {
			log.Println("Unknown panic happend: ", err)
		}
	}()
	router := api.NewRouter()

	fmt.Println("Server is running on port 8888")
	err := http.ListenAndServe(":8888", limits.Limit(router))
	if err != nil {
		log.Fatal(err)
	}

}
