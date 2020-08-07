package postgres

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	pu "github.com/nudelfabrik/portUpdate"
)

type PostgresService struct {
	pgdb     *sql.DB
	username string
	dbname   string
}

func NewBackendService() (pu.BackendService, error) {

	pgs := PostgresService{}
	pgs.username = "bene"
	pgs.dbname = "test"

	err := pgs.init()

	return &pgs, err
}

func (pgs *PostgresService) init() error {
	var err error

	login := fmt.Sprintf("user=%s dbname=%s", pgs.username, pgs.dbname)
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

func (pgs *PostgresService) AddEntries([]pu.Entry) error {
	return pu.ErrUnimplemented
}
