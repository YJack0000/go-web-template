package main

import (
	"log"

    "golang_backend_template/config"
	"golang_backend_template/internal"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)
}
