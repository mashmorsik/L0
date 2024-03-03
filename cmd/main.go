package main

import (
	"fmt"
	"github.com/mashmorsik/L0/infrastructure/data"
	"github.com/mashmorsik/L0/infrastructure/nats"
	log "github.com/mashmorsik/L0/pkg/logger"
)

func main() {
	log.Logger()

	err := nats.Connect()
	if err != nil {
		fmt.Println("error occurred", err)
	}

	conn := data.MustConnectPostgres()

	err = data.MustMigrate(conn)
	if err != nil {
		log.Errf("can't migrate, err: %s", err)
	}

	fmt.Printf("connected to db %v", conn)
}
