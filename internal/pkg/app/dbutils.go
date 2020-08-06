package app

import (
	"database/sql"
	"github.com/bentsolheim/go-app-utils/db"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/palantir/stacktrace"
	"github.com/sirupsen/logrus"
)

func ConnectAndMigrateDatabase(config db.DbConfig) (*sql.DB, error) {
	db, err := ConnectToDatabase(config)
	if err != nil {
		return nil, err
	}

	if err := ApplyMigrations(db, config); err != nil {
		return nil, stacktrace.Propagate(err, "unable to apply database migrations")
	}
	return db, nil
}

func ApplyMigrations(db *sql.DB, config db.DbConfig) error {
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

func ConnectToDatabase(config db.DbConfig) (*sql.DB, error) {
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
