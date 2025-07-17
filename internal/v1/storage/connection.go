package storage

import (
	"database/sql"
	docshell "docshell/internal/v1/config"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	// Expose config to local variable
	cfg := docshell.Config
	// Build connection string
	con := buildConnectionString(cfg)

	// Open database
	var err error
	db, err = sql.Open("postgres", con)
	if err != nil {
		log.Fatal(err)
	}

	// Setting up database setting
	db.SetConnMaxIdleTime(time.Duration(cfg.Service.DB.Settings.MaxIdleTime) * time.Minute)
	db.SetConnMaxLifetime(time.Duration(cfg.Service.DB.Settings.MaxConnLife) * time.Minute)
	db.SetMaxOpenConns(cfg.Service.DB.Settings.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Service.DB.Settings.MaxIdleConns)

	// Check connection
	if ok, err := CheckConnecton(); !ok {
		msg := "Connection to \"%s\" fault with %v"
		log.Fatalf(msg, con, err)
	}
	log.Printf("Connection to \"%s\" established.\n", con)
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

func buildConnectionString(cfg docshell.Configuration) string {
	format := "postgres://%s:%s@%s:%d/%s?sslmode=%s"
	c := cfg.Service.DB

	return fmt.Sprintf(
		format,
		c.Environment.PostgresUser,
		c.Environment.PostgresPassword,
		c.Host,
		c.Port,
		c.Environment.PostgresDB,
		c.SSLMode,
	)
}
