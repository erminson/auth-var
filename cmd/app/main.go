package main

import (
	"github.com/erminson/auth-var/config"
	"github.com/erminson/auth-var/internal/app"
	"log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	app.Run(cfg)
}
