package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"ex00/bd/config"
	"ex00/bd/internal/models"
	"ex00/bd/internal/utils"

	"github.com/google/uuid"
)

func handleHeartBeatGet(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(rw).Encode(config.HB)
	fmt.Println(config.HB)
}

func handleHeartBeatPost(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(rw).Encode(config.HB)
	var hb models.HeartBeat
	err := json.NewDecoder(r.Body).Decode(&hb)
	if err != nil {
		http.Error(rw, "Invalid request payload", http.StatusBadRequest)
		return
	}
	utils.CheckStrager(hb.Address[0])
}
func handlerSET(rw http.ResponseWriter, r *http.Request) {
	var GR models.SetReq
	err := json.NewDecoder(r.Body).Decode(&GR)
	if err != nil {
		http.Error(rw, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := uuid.Validate(GR.UUID4); err == nil {
		utils.Values[GR.UUID4] = GR.Value
		rw.WriteHeader(http.StatusCreated)
		rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
		message, _ := json.Marshal(GR)

		// log.Println(GR)
		// log.Println(string(message))

		count, _ := utils.ConsensusBuilding("/internal_set", message, http.StatusCreated)
		rw.Write([]byte(fmt.Sprintf("Created (%d replicas)", (count + 1))))
	} else {
		log.Println(errors.New("еrror: Key is not a proper UUID4"))
		http.Error(rw, "Invalid UUID", http.StatusBadRequest)
	}
}

func handlerGET(rw http.ResponseWriter, r *http.Request) {
	var ID models.GetReq
	if err := json.NewDecoder(r.Body).Decode(&ID); err != nil {
		http.Error(rw, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := uuid.Validate(ID.UUID4); err == nil {
		if val, inMap := utils.Values[ID.UUID4]; inMap {
			rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
			rw.Write([]byte(val)) //ТУТ
		} else {
			log.Println("Not found")
			http.Error(rw, "Value not found", http.StatusNotFound)
		}
	} else {
		log.Println(errors.New("error: Key is not a proper UUID4"))
		http.Error(rw, "Invalid UUID", http.StatusBadRequest)
	}
}

func handlerDELETE(rw http.ResponseWriter, r *http.Request) {
	var ID models.DelReq
	if err := json.NewDecoder(r.Body).Decode(&ID); err != nil {
		http.Error(rw, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := uuid.Validate(ID.UUID4); err == nil {
		if _, inMap := utils.Values[ID.UUID4]; inMap {
			rw.WriteHeader(http.StatusOK)
			rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
			delete(utils.Values, ID.UUID4)
			message, _ := json.Marshal(ID)
			count, _ := utils.ConsensusBuilding("/internal_delete", message, http.StatusOK)
			rw.Write([]byte(fmt.Sprintf("Deleted (%d replicas)", (count + 1))))
		} else {
			log.Println("Not found")
			http.Error(rw, "Value not found", http.StatusNotFound)
		}
	} else {
		log.Println(errors.New("error: Key is not a proper UUID4"))
		http.Error(rw, "Invalid UUID", http.StatusBadRequest)
	}
}

func handleMigrate(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(rw).Encode(utils.Values)
}
