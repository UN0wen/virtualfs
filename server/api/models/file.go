package models

import (
	"fmt"
	"time"

	"github.com/UN0wen/virtualfs/server/db"
	"github.com/pkg/errors"
)

// FileTableName is the name of the file table in the db
const (
	FileTableName = "files"
)

// FileTable represents the connection to the db instance
type FileTable struct {
	connection *db.Db
}

// File represents a single row in the FileTable
type File struct {
	ID          int64     `valid:"required" json:"id"`
	ParentDirID int64     `valid:"required" json:"imageurl"`
	Name        string    `valid:"required" json:"name"`
	Data        string    `valid:"required" json:"description"`
	Created     time.Time `valid:"required" json:"url"`
	Updated     time.Time `valid:"required" json:"currency"`
}

// NewFileTable creates a new table in the database for items.
// It takes a reference to an open db connection and returns the constructed table
func NewFileTable(db *db.Db) (fileTable FileTable, err error) {
	// Ensure connection is alive
	if db == nil {
		err = errors.New("Invalid database connection")
		return
	}
	fileTable.connection = db
	// Construct query
	query := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
			id SERIAL PRIMARY KEY, 
			parentdirid int REFERENCES %s(id),
			name TEXT NOT NULL,
			data TEXT, 
			created TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated TIMESTAMPTZ NOT NULL DEFAULT now(),
		)`, FileTableName, DirectoryTableName)
	// Create the actual table
	if err = fileTable.connection.CreateTable(query); err != nil {
		err = errors.Wrapf(err, "Could not initialize table %s", FileTableName)
	}
	return
}
