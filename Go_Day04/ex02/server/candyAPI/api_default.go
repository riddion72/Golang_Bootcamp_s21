package candyAPI

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

var CandyPrices = map[string]int32{
	"CE": 10,
	"AA": 15,
	"NT": 17,
	"DE": 21,
	"YR": 23,
}

func err_code(err error) {
	if err != nil {
		log.Fatal(err)
	}
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
			err_code(ans)
			return
		}
		need := param.CandyCount * cPrice

		if param.Money <= 0 || param.CandyCount <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			Answer := &InlineResponse400{}
			Answer.Error_ = "some error in input numerical values"
			ans := json.NewEncoder(w).Encode(Answer)
			err_code(ans)
		} else if param.Money >= need {
			cmd := exec.Command("./CowSay", "Thank you!")
			stdout, err := cmd.CombinedOutput()
			err_code(err)
			w.WriteHeader(http.StatusOK)
			Answer := &InlineResponse200{}
			Answer.Thanks = string(stdout)
			Answer.Change = param.Money - need
			ans := json.NewEncoder(w).Encode(Answer)
			err_code(ans)
		} else if param.Money < need {
			w.WriteHeader(http.StatusPaymentRequired)
			Answer := &InlineResponse400{}
			Answer.Error_ = fmt.Sprintf("You need %d money!", need)
			ans := json.NewEncoder(w).Encode(Answer)
			err_code(ans)
		}
	default:
		break
	}
}
