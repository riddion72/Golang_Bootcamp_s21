package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"main/db"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Store interface {
	GetPlaces(limit int, offset int) ([]db.Place, int, error)
	GetPlacesRecommend(lat, lon float64) ([]db.Place, error)
}

var (
	base      Store
	err       error
	secretKey = []byte("secretKey4221")
)

func init() {
	base, err = db.NewElasticsearch()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/api/places", apiPlaces)
	http.HandleFunc("/api/recommend", apiRecommend)
	http.HandleFunc("/api/get_token", getToken)
	log.Fatal(http.ListenAndServe(":8888", nil))
}

type Data struct {
	Places []db.Place
	Total  int
	Page   int
	Last   int
}

func home(w http.ResponseWriter, r *http.Request) {
	var res Data
	pageStr := r.URL.Query().Get("page")
	if res.Page, err = strconv.Atoi(pageStr); err != nil {
		http.Error(w, fmt.Sprintf("\"error %v\" : Invalid 'page' value: '%v'", http.StatusBadRequest, pageStr), http.StatusBadRequest)
		return
	}
	limit := 10
	offset := (res.Page - 1) * limit

	res.Places, res.Total, err = base.GetPlaces(limit, offset)
	res.Last = int(math.Ceil(float64(res.Total) / float64(limit)))
	tmpl, err := template.New("index.html").Funcs(
		template.FuncMap{
			"sum": sum,
			"sub": sub,
		},
	).ParseFiles("template/index.html")

	if res.Page > res.Last || res.Page < 1 {
		http.Error(w, fmt.Sprintf("\"error %v\" : Invalid 'page' value: '%v'", http.StatusBadRequest, res.Page), http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, res)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func sum(x, y int) int {
	return x + y
}

func sub(x, y int) int {
	return x - y
}

func apiPlaces(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Data
	pageStr := r.URL.Query().Get("page")
	if res.Page, err = strconv.Atoi(pageStr); err != nil {
		http.Error(w, fmt.Sprintf("\"error %v\" : Invalid 'page' value: '%v'", http.StatusBadRequest, pageStr), http.StatusBadRequest)
		return
	}
	limit := 10
	offset := (res.Page - 1) * limit
	res.Places, res.Total, err = base.GetPlaces(limit, offset)
	res.Last = int(math.Ceil(float64(res.Total) / float64(limit)))

	if res.Page > res.Last || res.Page < 1 {
		log.Println(err)
		http.Error(w, fmt.Sprintf("\"error %v\" : Invalid 'page' value: '%v'", http.StatusBadRequest, pageStr), http.StatusBadRequest)
		return
	}

	response, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	func() { _, _ = w.Write(response) }()
}

func apiRecommend(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	str := r.Header.Get("Authorization")
	token := strings.Split(str, ".")
	if len(token) != 3 {
		log.Println(token)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	if err := parseToken(token[1]); err != nil {
		log.Println(token[1])
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	lat, err := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("\"error %v\" Invalid 'lat' number", http.StatusBadRequest), http.StatusBadRequest)
	}
	lon, err := strconv.ParseFloat(r.URL.Query().Get("lon"), 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("\"error %v\" Invalid 'lon' number", http.StatusBadRequest), http.StatusBadRequest)
	}
	place, err := base.GetPlacesRecommend(lat, lon)
	if err != nil {
		log.Println(err)
		return
	}
	res := struct {
		Name   string     `json:"name"`
		Places []db.Place `json:"places"`
	}{
		Name:   "Recommendation",
		Places: place,
	}
	response, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	func() { _, _ = w.Write(response) }()
}

func getToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	token, err := NewJWT()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	res := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}
	response, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func NewJWT() (string, error) {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(12 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func parseToken(accessToken string) error {
	claims := jwt.MapClaims{}
	_, _ = jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	return nil
}
