package db

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

func EstablishDBConnection() *sql.DB {

	var (
		host       string
		port       int64
		portString string
		username   string
		password   string
		dbname     string
	)

	host, _ = os.LookupEnv("DB_HOST")
	username, _ = os.LookupEnv("DB_USER")
	password, _ = os.LookupEnv("DB_PASSWORD")
	portString, _ = os.LookupEnv("DB_PORT")
	dbname, _ = os.LookupEnv("DB_NAME")

	// port from string to int
	port, _ = strconv.ParseInt(portString, 10, 0)

	// prepare the DB connection string
	postgresConnectionString := fmt.Sprintf("host=%s port=%d password=%s user=%s dbname=%s sslmode=disable",
		host, port, password, username, dbname)

	// sql.open() validates the arguments BUT doesnt make an actual connection to the db
	// NOTE: sql.Open() function call never creates a connection to the database. Instead, it simply validates the arguments provided.
	db, err := sql.Open("postgres", postgresConnectionString)
	if err != nil {
		panic(err)
	}
	//defer db.Close()

	// Ping() opens a connection to the DB
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Established a successful Database connection!")

	return db
}
