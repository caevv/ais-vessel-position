package api

import (
	"fmt"
	"github.com/caevv/ais-vessel-position/internal/app/aisvesselposition/api/repository"
	"github.com/caevv/ais-vessel-position/pkg/aisvesselposition"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (app *Application) GetVesselPosition(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imo, err := strconv.Atoi(vars["imo"])
	if err != nil {
		respondError(w, fmt.Errorf("could not identify vessel IMO: %w", err))
		return
	}

	vesselRepository := repository.New(app.config.RepositoryJsonPath, app.config.Files)

	positions, err := vesselRepository.Positions(imo)
	if err != nil {
		respondError(w, err)
		return
	}

	distance := aisvesselposition.CalculateDistance(positions)

	respondWithJSON(
		w,
		http.StatusOK,
		map[string]float64{
			"statuteMiles":  distance.StatuteMiles,
			"kilometers":    distance.Kilometer,
			"nauticalMiles": distance.NauticalMiles,
		})
}
