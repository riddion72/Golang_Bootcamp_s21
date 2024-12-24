package api

import (
	"encoding/json"
	"net/http"

	"ex00/bd/internal/models"
	"ex00/bd/internal/utils"
)

func handlerInternalSET(rw http.ResponseWriter, r *http.Request) {
	var GR models.SetReq
	if err := json.NewDecoder(r.Body).Decode(&GR); err != nil {
		http.Error(rw, "Invalid request payload", http.StatusBadRequest)
		return
	}
	utils.Values[GR.UUID4] = GR.Value

	// log.Println(GR)
	// log.Println(utils.Values)

	rw.WriteHeader(http.StatusCreated)
}

func handlerInternalDELETE(rw http.ResponseWriter, r *http.Request) {
	var ID models.DelReq
	if err := json.NewDecoder(r.Body).Decode(&ID); err != nil {
		http.Error(rw, "Invalid request payload", http.StatusBadRequest)
		return
	}
	// if _, inMap := utils.Values[ID.UUID4]; inMap {
	// 	delete(utils.Values, ID.UUID4)
	// }
	delete(utils.Values, ID.UUID4) // Удаление без предварительной проверки НЕТ ПАНИКИ
	rw.WriteHeader(http.StatusOK)
}
