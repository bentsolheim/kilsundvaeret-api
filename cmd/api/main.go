package main

import (
	"fmt"
	"github.com/bentsolheim/kilsundvaeret-api/internal/pkg/app"
	"github.com/bentsolheim/kilsundvaeret-api/internal/pkg/app/service"
	"github.com/prometheus/client_golang/prometheus"
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
	metricsService := service.NewMetricsService(config.DataLoggerUrl, prometheus.DefaultRegisterer)
	router := app.CreateGinEngine(config)

	go metricsService.UpdateMetricsForever()

	return router.Run(fmt.Sprintf(":%s", config.ServerPort))
}
