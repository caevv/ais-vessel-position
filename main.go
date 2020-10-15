package main

import (
	"log"
	"runtime/debug"

	"github.com/caevv/ais-vessel-position/configs"
	"github.com/caevv/ais-vessel-position/api"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("%s: %s", r, string(debug.Stack()))
		}
	}()

	app := api.New(configs.New())
	app.Start()
}
