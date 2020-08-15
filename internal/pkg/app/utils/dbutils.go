package utils

import (
	"database/sql"
	dbutils "github.com/bentsolheim/go-app-utils/db"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/palantir/stacktrace"
	"github.com/prometheus/common/log"
	"github.com/sirupsen/logrus"
	"io"
)

func ConnectAndMigrateDatabase(config dbutils.DbConfig) (*sql.DB, error) {
	db, err := ConnectToDatabase(config)
	if err != nil {
		return nil, err
	}

	if err := ApplyMigrations(db, config); err != nil {
		return nil, stacktrace.Propagate(err, "unable to apply database migrations")
	}
	return db, nil
}

func ApplyMigrations(db *sql.DB, config dbutils.DbConfig) error {
	logrus.Info("Applying database migrations")
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		config.Name, driver)
	if err != nil {
		return stacktrace.Propagate(err, "unable to create migrate instance")
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return stacktrace.Propagate(err, "error while executing database migrations")
	}
	logrus.Info("Database migrations successfully applied")
	return nil
}

func ConnectToDatabase(config dbutils.DbConfig) (*sql.DB, error) {
	logrus.Info("Connecting to db: ", config.ConnectString("***"))
	db, err := sql.Open("mysql", config.ConnectString(""))
	if err != nil {
		return nil, stacktrace.Propagate(err, "error while connecting to database")
	}
	if err := db.Ping(); err != nil {
		return nil, stacktrace.Propagate(err, "unable to ping database")
	}
	return db, nil
}

func CloseSilently(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Warn(err)
	}
}
