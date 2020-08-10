package postgres

import (
	"database/sql"
	"fmt"
	"log"

	// Import SQL driver.
	_ "github.com/lib/pq"
	pu "github.com/nudelfabrik/portUpdate"
)

type Service struct {
	pgdb     *sql.DB
	username string
	dbname   string
}

func NewBackendService() (pu.BackendService, error) {
	pgs := Service{}
	pgs.username = "bene"
	pgs.dbname = "pgspu"

	err := pgs.init()

	return &pgs, err
}

func (pgs *Service) init() error {
	var err error

	login := fmt.Sprintf("user=%s dbname=%s sslmode=disable", pgs.username, pgs.dbname)
	pgs.pgdb, err = sql.Open("postgres", login)

	if err != nil {
		log.Println(err)

		return pu.ErrDBconnection
	}

	err = pgs.pgdb.Ping()

	if err != nil {
		log.Println(err)

		return pu.ErrDBconnection
	}

	return err
}

func (pgs *Service) AddEntries([]pu.Entry) error {
	return pu.ErrUnimplemented
}
