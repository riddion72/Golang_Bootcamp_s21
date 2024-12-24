package utils

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sort"
	"time"

	"ex00/bd/config"
	"ex00/bd/internal/models"
)

func RunTimeBeat() {
	for {
		localAddres := config.HB.Address
		for _, addr := range localAddres {
			if addr == config.Flager.Address {
				continue
			}
			resp, err := http.Get("http://" + addr + "/HeartBeat")
			if err != nil {
				log.Print("Error in runTimeBeat() Get", err)
				DeleteAddress(addr)
				CheckMainAddres(addr)
				continue
			}
			body, _ := io.ReadAll(resp.Body)
			var buffer models.HeartBeat
			if err = json.Unmarshal(body, &buffer); err != nil {
				log.Print("Error in runTimeBeat() Unmarshal: ", err)
				continue
			}
			if buffer.MainAddres != config.HB.MainAddres {
				config.HB.MainAddres = buffer.MainAddres
			}
			for _, Addr := range buffer.Address {
				CheckStrager(Addr)
			}
			resp.Body.Close()
		}
		time.Sleep(time.Second * 1)
	}
}

func CheckMainAddres(addr string) {
	if addr == config.HB.MainAddres {
		sort.Strings(config.HB.Address)
		config.HB.MainAddres = config.HB.Address[0]
	}
}

// DeleteAddress удаляет адрес из списка.
func DeleteAddress(addr string) {
	buffer := make([]string, 0)
	for _, v := range config.HB.Address {
		if v != addr {
			buffer = append(buffer, v)
		}
	}
	config.HB.Address = buffer
}
