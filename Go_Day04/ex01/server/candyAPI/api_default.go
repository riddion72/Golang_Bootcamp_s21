package candyAPI

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var CandyPrices = map[string]int32{
	"CE": 10,
	"AA": 15,
	"NT": 17,
	"DE": 21,
	"YR": 23,
}

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

		cPrice, ok := CandyPrices[param.CandyType]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			Answer := &InlineResponse400{}
			Answer.Error_ = "Invalid candy type"
			ans := json.NewEncoder(w).Encode(Answer)
			if ans != nil {
				log.Fatal(ans)
			}
			return
		}
		need := param.CandyCount * cPrice

		if param.Money <= 0 || param.CandyCount <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			Answer := &InlineResponse400{}
			Answer.Error_ = "some error in input numerical values"
			ans := json.NewEncoder(w).Encode(Answer)
			if ans != nil {
				log.Fatal(ans)
			}
		} else if param.Money >= need {
			w.WriteHeader(http.StatusOK)
			Answer := &InlineResponse200{}
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
