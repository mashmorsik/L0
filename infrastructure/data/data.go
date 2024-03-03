package data

import (
	"errors"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/mashmorsik/L0/pkg/models"

	log "github.com/mashmorsik/L0/pkg/logger"
	"os"

	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Data struct {
	db *sql.DB
}

func NewData(db *sql.DB) *Data {
	if db == nil {
		panic("db is nil")
	}
	return &Data{db: db}
}

func (r *Data) Db() *sql.DB {
	return r.db
}

func MustConnectPostgres() *sql.DB {
	connectionStr := "postgres://postgres:mysecretpassword@localhost:5432/postgres?sslmode=disable&application_name=l0&connect_timeout=5"

	connection, err := sql.Open("postgres", connectionStr)
	if err != nil {
		panic(err)
	}

	if err = connection.Ping(); err != nil {
		panic(err)
	}

	return connection
}

func MustMigrate(connection *sql.DB) error {
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
		log.Errf("can't migrate: %s", err)
		return err
	}

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Infof("no changes in migration, skip")

		} else {
			panic(err)
		}
	}
	return nil
}

func (r *Data) AddOrder(order models.Order) error {
	sqlAddOrder := `
	INSERT INTO order
	VALUES($1)`
	res, err := r.db.Exec(sqlAddOrder, order)
	if err != nil {
		log.Errf("can't add order, err: %s", err)
		return err
	}
	ra, _ := res.RowsAffected()
	fmt.Printf("rows affected: %v", ra)

	return nil
}
