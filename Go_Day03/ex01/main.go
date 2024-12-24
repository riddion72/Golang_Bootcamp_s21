package main

import (
	"html/template"
	"log"
	"main/db"
	"math"
	"net/http"
	"strconv"
)

type Store interface {
	GetPlaces(limit int, offset int) ([]db.Place, int, error)
}

var (
	base Store
	err  error
)

func init() {
	base, err = db.NewElasticsearch()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/", home)
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
	if res.Page, err = strconv.Atoi(r.URL.Query().Get("page")); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	limit := 10
	offset := (res.Page - 1) * limit

	res.Places, res.Total, err = base.GetPlaces(limit, offset)
	res.Last = int(math.Ceil(float64(res.Total) / float64(limit)))

	if res.Page > res.Last || res.Page < 1 {
		http.Error(w, "Erorr 400\n BadRequest \nInvalid 'page' value: 'foo'", http.StatusBadRequest)
		return
	}

	tmpl, err := template.New("index.html").Funcs(
		template.FuncMap{
			"sum": sum,
			"sub": sub,
		},
	).ParseFiles("template/index.html")

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
