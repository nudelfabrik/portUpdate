package portUpdate

import "errors"

var (
	ErrUnimplemented   = errors.New("Database: Operation not implemented")
	ErrOperationFailed = errors.New("Database: Operation failed")
	ErrDBconnection    = errors.New("Database: Connection to DB failed")
)

type BackendService interface {
	AddEntries([]Entry) error
}
