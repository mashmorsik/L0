package data

import (
	"context"
	"errors"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	log "github.com/mashmorsik/L0/pkg/logger"
	"os"

	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Data struct {
	Ctx context.Context
	db  *sql.DB
}

func NewData(ctx context.Context, db *sql.DB) *Data {
	if db == nil {
		panic("db is nil")
	}
	return &Data{Ctx: ctx, db: db}
}

func (r *Data) Master() *sql.DB {
	return r.db
}

func MustConnectPostgres(ctx context.Context) *sql.DB {
	connectionStr := "postgres://postgres:mysecretpassword@localhost:5432/postgres?sslmode=disable&application_name=l0&connect_timeout=5"

	connection, err := sql.Open("postgres", connectionStr)
	if err != nil {
		panic(err)
	}

	if err = connection.Ping(); err != nil {
		panic(err)
	}

	go func() {
		<-ctx.Done()
		err = connection.Close()
		if err != nil {
			log.Errf("can't close database connection, err: %s", err)
			return
		}
	}()

	log.Infof("connected to db: %v", connection)
	return connection
}

func MustMigrate(connection *sql.DB) {
	driver, err := postgres.WithInstance(connection, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	migrationPath := fmt.Sprintf("file://%s/migration", path)
	fmt.Printf("migrationPath : %s\n", migrationPath)

	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		"l0", driver)
	if err != nil {
		panic(err)
	}

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Infof("no changes in migration, skip")

		} else {
			panic(err)
		}
	}
}
