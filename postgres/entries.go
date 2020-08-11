package postgres

import (
	"log"

	"github.com/lib/pq"
	pu "github.com/nudelfabrik/portUpdate"
)

// Add new Entries to Backend
func (pgs *BackendService) AddEntries(entries []pu.Entry) error {
	// Start Transaction
	tx, err := pgs.pgdb.Begin()
	if err != nil {
		log.Println(err)
		return pu.ErrOperationFailed
	}

	// Prepare Statement
	stmt, err := tx.Prepare(pq.CopyIn("entries", "date", "author", "ports", "description"))
	if err != nil {
		log.Println(err)
		return pu.ErrOperationFailed
	}

	// Add Entries
	for _, entry := range entries {
		_, err = stmt.Exec(entry.Date, entry.Author, pq.Array(entry.Ports), entry.Description)
		if err != nil {
			log.Println(err)
			tx.Rollback()
			return pu.ErrOperationFailed
		}
	}

	// Flush all data with empty exec
	_, err = stmt.Exec()
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return pu.ErrOperationFailed
	}

	// Close and Commit Transaction
	err = stmt.Close()
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return pu.ErrOperationFailed
	}
	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return pu.ErrOperationFailed
	}

	return nil
}
