package postgres

import (
	"database/sql"
	"fmt"
	"log"

	// Import SQL driver.
	_ "github.com/lib/pq"
	pu "github.com/nudelfabrik/portUpdate"
)

type BackendService struct {
	pgdb     *sql.DB
	username string
	dbname   string
}

// Create and Setup a new BackendService with Postgres
func NewBackendService() (pu.BackendService, error) {
	pgs := BackendService{}
	pgs.username = "bene"
	pgs.dbname = "pgspu"

	err := pgs.init()

	return &pgs, err
}

// Establish connection with database and create Tables.
func (pgs *BackendService) init() (err error) {
	// Establish connection
	login := fmt.Sprintf("user=%s dbname=%s sslmode=disable", pgs.username, pgs.dbname)
	pgs.pgdb, err = sql.Open("postgres", login)

	if err != nil {
		log.Println(err)

		return pu.ErrDBconnection
	}

	// Test if connection works
	err = pgs.pgdb.Ping()

	if err != nil {
		log.Println(err)

		return pu.ErrDBconnection
	}

	// Capture panicing Table creations
	defer func() {
		if r := recover(); r != nil {
			log.Println("Table Setup failed: ", r)

			err = pu.ErrOperationFailed
		}
	}()

	// Create Tables
	pgs.createTable("entries", "id SERIAL PRIMARY KEY, date DATE, author text, ports text[], description text")

	return nil
}

// Create a single new Table, if it does not exist.
func (pgs *BackendService) createTable(name, columns string) {
	str := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", name, columns)

	_, err := pgs.pgdb.Exec(str)
	if err != nil {
		panic(err)
	}
}
