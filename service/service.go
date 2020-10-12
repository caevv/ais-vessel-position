package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/caevv/ais-vessel-position/data"
	"github.com/caevv/ais-vessel-position/env"
	"github.com/caevv/ais-vessel-position/repository"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type Application struct {
	Router *mux.Router
	config *env.Config
}

func New(config *env.Config) *Application {
	router := mux.NewRouter()

	app := &Application{Router: router, config: config}

	router.HandleFunc("/vessel/position/{imo}", app.GetVesselPosition).Methods("GET")

	return app
}

func (a *Application) Start() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", a.config.Port), a.Router))
}

func (app *Application) GetVesselPosition(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)

	imo, err := strconv.Atoi(vars["imo"])
	if err != nil {
		respondError(w, errors.Wrap(err, "could not identify vessel IMO "))
		return
	}

	vesselRepository := repository.New(app.config.RepositoryJsonPath)

	positions, err := vesselRepository.Positions(imo)
	if err != nil {
		respondError(w, err)
		return
	}

	distance := data.CalculateDistance(positions)

	respondWithJSON(
		w,
		http.StatusCreated,
		map[string]float64{
			"statuteMiles":  distance.StatuteMiles,
			"kilometers":    distance.Kilometer,
			"nauticalMiles": distance.NauticalMiles,
		})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		log.Print(err)
	}
}

func respondError(w http.ResponseWriter, err error) {
	respondWithJSON(w, http.StatusBadRequest, map[string]string{
		"error": err.Error(),
	})
}
