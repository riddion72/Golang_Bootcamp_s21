/*
 * Candy Server
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package candyAPI

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	// "fmt"
)

func BuyCandy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// w.WriteHeader(http.StatusOK)
	switch r.Method {
	case "POST":
		var cPrice int32
		body := json.NewDecoder(r.Body)
		defer r.Body.Close()
		param := &Order{}
		if err := body.Decode(param); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if param.CandyType == "CE" {
			cPrice = 10
		} else if param.CandyType == "AA" {
			cPrice = 15
		} else if param.CandyType == "NT" {
			cPrice = 17
		} else if param.CandyType == "DE" {
			cPrice = 21
		} else if param.CandyType == "YR" {
			cPrice = 23
		}
		need := param.CandyCount * cPrice

		if param.Money <= 0 || param.CandyCount <= 0 || cPrice == 0 {
			w.WriteHeader(http.StatusBadRequest)
			Answer := &InlineResponse400{}
			Answer.Error_ = "some error in input data"
			ans := json.NewEncoder(w).Encode(Answer)
			if ans != nil {
				log.Fatal(ans)
			}
		} else if param.Money >= need {
			w.WriteHeader(http.StatusCreated)
			Answer := &InlineResponse201{}
			Answer.Thanks = "Thank you!"
			Answer.Change = param.Money - need
			ans := json.NewEncoder(w).Encode(Answer)
			if ans != nil {
				log.Fatal(ans)
			}
		} else if param.Money < need {
			w.WriteHeader(http.StatusPaymentRequired)
			Answer := &InlineResponse400{}
			Answer.Error_ = fmt.Sprintf("You need %d money!", need)
			ans := json.NewEncoder(w).Encode(Answer)
			if ans != nil {
				log.Fatal(ans)
			}
		}
	default:
		break
	}
}
