package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"ex00/bd/config"
	"ex00/bd/internal/models"
)

var Values = make(map[string]string)

func InitBeat() error {
	message, _ := json.Marshal(config.HB)
	url := "http://" + config.Flager.BDaddress + "/HeartBeat"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(message))
	if err != nil {
		log.Print("Error while sending heartbeat: ", err)
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print("Error while reading response body: ", err)
		return err
	}
	defer resp.Body.Close()

	var buffer models.HeartBeat
	if err = json.Unmarshal(body, &buffer); err != nil {
		log.Print("Error while unmarshalling response: ", err)
		return err
	}
	if config.HB.Replication != buffer.Replication {
		return errors.New("replication mismatch")
	}
	log.Print("Heartbeat received: ", config.HB.Address, config.HB.MainAddres)
	config.HB.Address = append(config.HB.Address, buffer.Address...)
	config.HB.MainAddres = buffer.MainAddres
	log.Print("Updated addresses: ", config.HB.Address, config.HB.MainAddres)
	return nil
}

// потом сотрем логи
func MigrateAddress() {
	resp, err := http.Get("http://" + config.Flager.BDaddress + "/Migrate")
	if err != nil {
		log.Print("Error while initiating address migration: ", err)
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print("Error while reading migration response body: ", err)
	}
	defer resp.Body.Close()

	buffer := make(map[string]string, 1)
	if err = json.Unmarshal(body, &buffer); err != nil {
		log.Print("Error while unmarshalling migration response: ", err)
	}
	log.Print("Current values: ", Values)
	Values = buffer
	log.Print("Updated values after migration: ", Values)
}
