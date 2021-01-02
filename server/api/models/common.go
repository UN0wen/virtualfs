package models

import (
	"sync"

	"github.com/UN0wen/virtualfs/server/db"
	"github.com/UN0wen/virtualfs/server/utils"
)

// Represents the layer for the model by exposing the
// different models' tables.
type layer struct {
	File      *FileTable
	Directory *DirectoryTable
}

// Singleton reference to the model layer.
var instance *layer

// Lock for running only once.
var once sync.Once

// LayerInstance gets the static singleton reference
// using double check synchronization.
// It returns the reference to the layer.
func LayerInstance() *layer {
	once.Do(func() {
		// Create DB only once
		db, err := db.Setup(db.Config{
			Host:     utils.DBHost,
			Port:     utils.DBPort,
			User:     utils.DBUser,
			Password: utils.DBPassword,
			Database: utils.DBName,
		})
		utils.CheckError(err)
		// Create all the tables
		directoryTable, err := NewDirectoryTable(&db)
		utils.CheckError(err)
		fileTable, err := NewFileTable(&db)
		utils.CheckError(err)
		// Create the layer only once
		instance = &layer{
			Directory: &directoryTable,
			File:      &fileTable,
		}
	})
	return instance
}
