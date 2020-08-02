package app

import (
	"github.com/bentsolheim/go-app-utils/utils"
)

type AppConfig struct {
	ServerPort    string
	DataLoggerUrl string
	DataLoggerId  string
}

func ReadAppConfig() AppConfig {
	e := utils.GetEnvOrDefault
	return AppConfig{
		e("SERVER_PORT", "8080"),
		e("DATALOGGER_URL", "http://localhost:8081"),
		e("DATALOGGER_ID", "bua"),
	}
}
