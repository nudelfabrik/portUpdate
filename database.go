package portUpdate

import "errors"

var (
	ErrUnimplemented = errors.New("Database: Operation not implemented")
	ErrDBconnection  = errors.New("Database: Connection to DB failed")
)

type BackendService interface {
	AddEntries([]Entry) error
}
