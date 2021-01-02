package models

import (
	"fmt"
	"time"

	"github.com/UN0wen/virtualfs/server/db"
	"github.com/pkg/errors"
)

// DirectoryTableName is the name of the dir table in the db
const (
	DirectoryTableName = "directories"
)

// DirectoryTable represents the connection to the db instance
type DirectoryTable struct {
	connection *db.Db
}

// Directory represents a single row in the FileTable
type Directory struct {
	ID          int64     `valid:"required" json:"id"`
	ParentDirID int64     `valid:"required" json:"imageurl"`
	Name        string    `valid:"required" json:"name"`
	Created     time.Time `valid:"required" json:"url"`
	Updated     time.Time `valid:"required" json:"currency"`
}

// NewDirectoryTable creates a new table in the database for items.
// It takes a reference to an open db connection and returns the constructed table
func NewDirectoryTable(db *db.Db) (directoryTable DirectoryTable, err error) {
	// Ensure connection is alive
	if db == nil {
		err = errors.New("Invalid database connection")
		return
	}
	directoryTable.connection = db
	// Construct query
	query := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
			id SERIAL PRIMARY KEY, 
			parentdirid int REFERENCES %s(id),
			name TEXT NOT NULL,
			data TEXT NOT NULL, 
			created TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated TIMESTAMPTZ NOT NULL DEFAULT now(),
		)`, DirectoryTableName, DirectoryTableName)
	// Create the actual table
	if err = directoryTable.connection.CreateTable(query); err != nil {
		err = errors.Wrapf(err, "Could not initialize table %s", FileTableName)
	}
	return
}
