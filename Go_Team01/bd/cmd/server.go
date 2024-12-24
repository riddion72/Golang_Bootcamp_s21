package main

import (
	"fmt"
	"log"
	"net/http"

	"ex00/bd/api"
	"ex00/bd/config"
	"ex00/bd/internal/utils"
)

func init() {
	config.InitFlags()
}
func main() {
	router := api.NewRouter()
	fmt.Println(config.Flager.Address)
	if config.CheckConfig() {
		if err := utils.InitBeat(); err != nil {
			log.Fatal(err)
		}
		utils.MigrateAddress()
	}
	go utils.RunTimeBeat()
	http.ListenAndServe(config.Flager.Address, router)
}
