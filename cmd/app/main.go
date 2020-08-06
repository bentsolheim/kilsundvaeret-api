package main

import (
	"fmt"
	"github.com/bentsolheim/kilsundvaeret-api/internal/pkg/app"
	"github.com/bentsolheim/kilsundvaeret-api/internal/pkg/app/service"
	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {

	config := app.ReadAppConfig()

	if err := app.ConfigureLogging(config.LogLevel); err != nil {
		return err
	}

	db, err := app.ConnectAndMigrateDatabase(config.DbConfig)
	if err != nil {
		return err
	}
	defer db.Close()

	metricsService := service.NewMetricsService(config.DataReceiverUrl, prometheus.DefaultRegisterer)
	router := app.CreateGinEngine(config)

	go metricsService.UpdateMetricsForever("bua", 60)

	return router.Run(fmt.Sprintf(":%s", config.ServerPort))
}
