package main

import (
	"fmt"
	"log"
	"messenger/internal"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "messenger"
)

func buildConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}

func main() {
	aa, err := internal.NewAdminApi(buildConnectionString())
	if err != nil {
		log.Fatal(err)
	}

	aa.Listen(":8081")
}
