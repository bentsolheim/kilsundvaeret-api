package app

import (
	"github.com/bentsolheim/go-app-utils/db"
	"github.com/bentsolheim/go-app-utils/utils"
)

type AppConfig struct {
	DbConfig        db.DbConfig
	LogLevel        string
	MetProxyUrl     string
	DataReceiverUrl string
	ServerPort      string
}

func ReadAppConfig() AppConfig {
	e := utils.GetEnvOrDefault
	return AppConfig{
		db.ReadDbConfig(db.DbConfig{
			User:     "root",
			Password: "devpass",
			Host:     "localhost",
			Port:     "3306",
			Name:     "kilsundvaeret",
		}),
		e("LOG_LEVEL", "debug"),
		e("MET_PROXY_URL", "http://localhost:8082"),
		e("DATA_RECEIVER_URL", "http://localhost:8081"),
		e("SERVER_PORT", "8080"),
	}
}
