package utils

import (
	"bytes"
	"log"
	"net/http"

	"ex00/bd/config"
)

func CheckStrager(addr string) {
	for _, str := range config.HB.Address {
		if str == addr {
			return
		}
	}
	config.HB.Address = append(config.HB.Address, addr)
	log.Println("Stranger detected: ", addr, config.HB.Address)
}

func ConsensusBuilding(end string, jsonRequest []byte, ststusCode int) (int, error) {
	var count int
	for _, addr := range config.HB.Address {
		if addr != config.Flager.Address {
			resp, err := http.Post("http://"+addr+end, "application/json", bytes.NewBuffer(jsonRequest))
			if err != nil {
				return count, err
			}
			if resp.StatusCode == ststusCode {
				count++
			}
			resp.Body.Close()
		}
	}
	return count, nil
}
