package app

import (
	"github.com/bentsolheim/go-app-utils/db"
	"github.com/bentsolheim/go-app-utils/utils"
	"github.com/palantir/stacktrace"
	"strconv"
)

type AppConfig struct {
	DbConfig                   db.DbConfig
	LogLevel                   string
	StracktraceInErrorMessages bool
	MetProxyUrl                string
	DataReceiverUrl            string
	ServerPort                 string
}

func ReadAppConfig() (*AppConfig, error) {
	e := utils.GetEnvOrDefault
	stacktraceInErrorMessages, err := strconv.ParseBool(e("STACKTRACE_IN_ERROR_MESSAGES", "true"))
	if err != nil {
		return nil, stacktrace.Propagate(err, "unable to read STACKTRACE_IN_ERROR_MESSAGES config")
	}
	return &AppConfig{
		db.ReadDbConfig(db.DbConfig{
			User:     "root",
			Password: "devpass",
			Host:     "localhost",
			Port:     "3306",
			Name:     "kilsundvaeret",
		}),
		e("LOG_LEVEL", "debug"),
		stacktraceInErrorMessages,
		e("MET_PROXY_URL", "http://localhost:8082"),
		e("DATA_RECEIVER_URL", "http://localhost:8081"),
		e("SERVER_PORT", "8080"),
	}, nil
}
