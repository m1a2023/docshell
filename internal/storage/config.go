package storage

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Database variables
var (
	CONNECTION_STRING 		string 
	MAX_IDLE_TIME					int 
	MAX_CONNECTION_LIFE		int
	MAX_OPEN_CONNECTIONS 	int
	MAX_IDLE_CONNECTIONS 	int	 
)

var db *sql.DB

func init() {
	// Get environment variable
	CONNECTION_STRING = os.Getenv("DSN")
	
	// Open database 
	var err error
	db, err = sql.Open("postgres", CONNECTION_STRING)
	if err != nil {
		log.Fatal(err)
	}
	
	// Setting up database connection settings
	db.SetConnMaxIdleTime(time.Duration(MAX_IDLE_TIME) * time.Minute)
	db.SetConnMaxLifetime(time.Duration(MAX_CONNECTION_LIFE) * time.Minute)
	db.SetMaxOpenConns(MAX_OPEN_CONNECTIONS)
	db.SetMaxIdleConns(MAX_IDLE_CONNECTIONS)

	// Check connection 
	if ok, err := CheckConnecton(); !ok {
		log.Fatal(err)
	} else {
		log.Printf("Connection to \"%s\" established.\n", CONNECTION_STRING)
	}
}


func GetConnection() *sql.DB {
	return db
}


func CheckConnecton() (bool, error) {
	if err := db.Ping(); err != nil {
		return false, err 
	}
	return true, nil
}
