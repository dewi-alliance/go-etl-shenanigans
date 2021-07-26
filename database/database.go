package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv" // package used to read the .env file
	_ "github.com/lib/pq"      // postgres golang driver
)

// DB connection is global
var DB *sql.DB

// Start database
func Start() {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Open the connection
	DB, err = sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	// check the connection
	err = DB.Ping()

	if err != nil {
		panic(err)
	}

	log.Println("Database Successfully connected!")
}
