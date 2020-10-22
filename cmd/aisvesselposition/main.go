package main

import (
	"github.com/caevv/ais-vessel-position/internal/app/aisvesselposition/api"
	"os"
	"os/signal"
	"syscall"

	"github.com/caevv/ais-vessel-position/configs"
)

func main() {
	stopServer := make(chan bool, 1)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-stop
		stopServer <- true
	}()

	app := api.New(configs.New())
	app.Start(stopServer, make(chan bool, 1))
}
