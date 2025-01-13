package main

import (
	"fmt"
	"log"
	"messenger/internal"
)

const (
	host     = "postgres"
	port     = 5432
	user     = "myuser"
	password = "mypassword"
	dbname   = "messenger"

	key  = "*secret key*"
	path = "localhost"
)

func buildConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}

func main() {
	aa, err := internal.NewClientApi(buildConnectionString(), []byte(key), path)
	if err != nil {
		log.Fatal(err)
	}

	aa.Listen(":8080")
}
