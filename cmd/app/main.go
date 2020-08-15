package main

import (
	"fmt"
	"github.com/bentsolheim/kilsundvaeret-api/internal/pkg/app"
	"github.com/bentsolheim/kilsundvaeret-api/internal/pkg/app/controller"
	"github.com/bentsolheim/kilsundvaeret-api/internal/pkg/app/service"
	"github.com/bentsolheim/kilsundvaeret-api/internal/pkg/app/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/palantir/stacktrace"
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

	config, err := app.ReadAppConfig()
	if err != nil {
		return stacktrace.Propagate(err, "error while reading application configuration")
	}
	if !config.StracktraceInErrorMessages {
		stacktrace.DefaultFormat = stacktrace.FormatBrief
	}

	if err := utils.ConfigureLogging(config.LogLevel); err != nil {
		return err
	}

	db, err := utils.ConnectAndMigrateDatabase(config.DbConfig)
	if err != nil {
		return err
	}
	defer utils.CloseSilently(db)

	metricsService := service.NewMetricsService(config.DataReceiverUrl, prometheus.DefaultRegisterer)
	weatherMetricsService := service.NewWeatherMetricsService(db)

	weatherMetricsController := controller.NewWeatherMetricsController(weatherMetricsService)
	weatherMetricsReportController := controller.NewWeatherMetricsReportController(weatherMetricsService)

	router := app.CreateGinEngine(*config, weatherMetricsController, weatherMetricsReportController)

	go metricsService.UpdateMetricsForever("bua", 60)

	return router.Run(fmt.Sprintf(":%s", config.ServerPort))
}
