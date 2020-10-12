package main

import (
	"log"
	"runtime/debug"

	"github.com/caevv/ais-vessel-position/env"
	"github.com/caevv/ais-vessel-position/service"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("%s: %s", r, string(debug.Stack()))
		}
	}()

	app := service.New(env.New())
	app.Start()
}
