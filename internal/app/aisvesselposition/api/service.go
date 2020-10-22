package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/caevv/ais-vessel-position/configs"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type Application struct {
	Router *mux.Router
	config *configs.Config
	server *http.Server
}

func New(config *configs.Config) *Application {
	router := mux.NewRouter()
	router.Use(loggingMiddleware)

	app := &Application{
		Router: router,
		config: config,
		server: &http.Server{
			Handler: router,
			Addr:    fmt.Sprintf(":%d", config.Port),
		},
	}

	router.HandleFunc("/vessel/position/{imo}", app.GetVesselPosition).Methods("GET")
	router.HandleFunc("/health", HealthCheckHandler)

	return app
}

func (app *Application) Start(stopServer chan bool, startedSignal chan bool) {
	go func() {
		log.Infof("Server started on %s", app.server.Addr)
		startedSignal <- true

		<-stopServer
		log.Info("server is shutting down")
		err := app.server.Shutdown(context.Background())
		log.Info(fmt.Errorf("failed to shutdown server: %w", err))
	}()

	if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debugln(r.RequestURI)

		next.ServeHTTP(w, r)
	})
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err := io.WriteString(w, `{"alive": true}`)
	if err != nil {
		log.Info(fmt.Errorf("could not write health check response: %w", err))
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

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
