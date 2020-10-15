package main

import (
	"log"
	"runtime/debug"

	"github.com/caevv/ais-vessel-position/env"
	"github.com/caevv/ais-vessel-position/api"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("%s: %s", r, string(debug.Stack()))
		}
	}()

	app := api.New(env.New())
	app.Start()
}
