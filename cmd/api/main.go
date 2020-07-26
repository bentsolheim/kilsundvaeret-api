package main

import (
	"fmt"
	"github.com/bentsolheim/kilsundvaeret-api/internal/pkg/app"
	"log"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {

	config := app.ReadAppConfig()
	router := app.CreateGinEngine(config)
	return router.Run(fmt.Sprintf(":%s", config.ServerPort))
}
