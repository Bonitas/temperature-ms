package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	guuid "github.com/google/uuid"
	"it.schwarz/city/app/model"
	"it.schwarz/city/config"
)

var cities []*model.City

func initCityList() {
	cities = cities[:0]
	cities = append(cities, model.NewCity(guuid.New().String(), "Berlin", 10, "DE"))
	cities = append(cities, model.NewCity(guuid.New().String(), "Stuttgart", 15, "DE"))

	cities = append(cities, model.NewCity(guuid.New().String(), "Sofia", 11, "BG"))
	cities = append(cities, model.NewCity(guuid.New().String(), "Varna", 8, "BG"))

	cities = append(cities, model.NewCity(guuid.New().String(), "Paris", 21, "FR"))
	cities = append(cities, model.NewCity(guuid.New().String(), "Lyon", 28, "FR"))
}

func filter(vs []*model.City, f func(string) bool) []*model.City {
	vsf := make([]*model.City, 0)
	for _, v := range vs {
		if f(v.Country) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

// GetcityMember get city members
func GetCities(config *config.Config, w http.ResponseWriter, r *http.Request) {
	initCityList()
	query := r.URL.Query().Get("country")
	if len(query) > 0 {
		ccities := filter(cities, func(v string) bool {
			return strings.Contains(v, query)
		})
		ResponseWriter(w, http.StatusOK, ccities)
	} else {
		ResponseWriter(w, http.StatusOK, cities)
	}
}

func AddNewCity(config *config.Config, w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		country := r.URL.Query().Get("country")
		name := r.URL.Query().Get("city")
		temp, err := strconv.Atoi(r.URL.Query().Get("temp"))
		if err != nil {
			log.Fatal(err)
		}

		newCity := model.NewCity(guuid.New().String(), name, temp, country)

		cities = append(cities, newCity)

		filtered := filter(cities, func(v string) bool {
			return strings.Contains(v, country)
		})

		ResponseWriter(w, http.StatusCreated, filtered)
	} else {
		ResponseWriter(w, http.StatusBadRequest, "BAD REQUEST")
	}
}

// ResponseWriter will write result in http.ResponseWriter
func ResponseWriter(res http.ResponseWriter, statusCode int, data interface{}) error {
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(statusCode)
	err := json.NewEncoder(res).Encode(data)
	return err
}
