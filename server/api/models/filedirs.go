package models

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/UN0wen/virtualfs/server/db"
	"github.com/UN0wen/virtualfs/server/utils"
	"github.com/asaskevich/govalidator"
	"github.com/georgysavva/scany/pgxscan"

	"github.com/pkg/errors"
)

// TimeParam are variables where there needs to be
// conversion from time.Time to a string that can
// be parsed in timestamptz
// SkipParam are params that we skip on update
var (
	TimeParam = map[string]bool{
		"created": true,
		"updated": true,
	}
	SkipParam = map[string]bool{
		"created": true,
		"id":      true,
	}
)

// IsUndeclared uses reflection to see if the value of the field is set or not.
// It takes in an interface to reflect on and returns a boolean if the field is set or not.
func isUndeclared(v interface{}) bool {
	return reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface())
}

// FileDirsTableName is the name of the filedirs table in the db
const (
	FileDirsTableName = "filedirs"
)

// FileDirsTable represents the connection to the db instance
type FileDirsTable struct {
	connection *db.Db
}

// FileDirs represents a single row in the FileTable
type FileDirs struct {
	ID      int64     `valid:"-" json:"id"`
	Name    string    `valid:"required" json:"name"`
	Data    string    `valid:"-" json:"data"`
	Created time.Time `valid:"-" json:"created"`
	Updated time.Time `valid:"-" json:"updated"`
	Size    int64     `valid:"-" json:"size"`
	Type    string    `valid:"-" json:"filetype"`
	Path    string    `valid:"-" json:"path"`
}

// NewFileDirsTable creates a new table in the database for files and directories.
// It takes a reference to an open db connection and returns the constructed table
// This code doesn't actually run since the database is created using the migration script.
func NewFileDirsTable(db *db.Db) (fileDirsTable FileDirsTable, err error) {
	// Ensure connection is alive
	if db == nil {
		err = errors.New("Invalid database connection")
		return
	}
	fileDirsTable.connection = db
	// Construct query
	query := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		id serial PRIMARY KEY,
		name text NOT NULL,
		data text,
		created timestamptz NOT NULL DEFAULT now(),
		updated timestamptz NOT NULL DEFAULT now(),
		size int DEFAULT 0,
		TYPE filetype NOT NULL,
		path ltree
	);`, FileDirsTableName)
	// Create the actual table
	if err = fileDirsTable.connection.CreateTable(query); err != nil {
		err = errors.Wrapf(err, "Could not initialize table %s", FileDirsTableName)
	}
	return
}

// GetPath returns an entry in the FileDirs table
func (table *FileDirsTable) GetPath(path string) (item *FileDirs, err error) {
	var result []*FileDirs
	ltree, err := utils.PathToLtree(path)
	if err != nil {
		err = errors.Wrapf(err, "Could not convert path during GetPath")
		return
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE path='%s';", FileDirsTableName, ltree)

	utils.Sugar.Infof("SQL Query: %s", query)

	err = pgxscan.Select(context.Background(), table.connection.Pool, &result, query)
	if err != nil {
		err = errors.Wrapf(err, "GetPath query failed to execute")
		return
	}

	if len(result) == 0 {
		err = errors.New(fmt.Sprintf("No item at path %s found", path))
		return
	}

	item = result[0]
	return
}

// GetAllPath returns all entries under a directory in the FileDirs table
func (table *FileDirsTable) GetAllPath(path string, levels int) (items []*FileDirs, err error) {
	ltree, err := utils.PathToLtree(path)
	if err != nil {
		err = errors.Wrapf(err, "Could not convert path during GetPath")
		return
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE path~'%s.*{1,%d}';", FileDirsTableName, ltree, levels)

	utils.Sugar.Infof("SQL Query: %s", query)

	err = pgxscan.Select(context.Background(), table.connection.Pool, &items, query)
	if err != nil {
		err = errors.Wrapf(err, "GetPath query failed to execute")
		return
	}

	return
}

// Insert inserts either a file or a directory into the fs
func (table *FileDirsTable) Insert(newItem *FileDirs, path string) (insertedItem *FileDirs, err error) {
	var results []*FileDirs
	_, err = govalidator.ValidateStruct(newItem)
	if err != nil {
		err = errors.Wrapf(err, "Missing required fields")
		return
	}
	// Check parent dir path if it exists
	item, err := table.GetPath(path)
	if err != nil {
		err = errors.Wrapf(err, "Parent directory does not exist")
		return
	}

	// Set parent dir id and newItem's path
	newItem.Path = item.Path + "." + newItem.Name

	// Set item size
	if newItem.Data != "" {
		newItem.Size = int64(len(newItem.Data))
		newItem.Type = "file"
	} else {
		newItem.Size = 0
		newItem.Type = "directory"
	}
	// Set time
	currentTime := time.Now()
	newItem.Created = currentTime
	newItem.Updated = currentTime

	// Update parent dir using the retrieved item earlier
	item.Updated = currentTime
	parentPath, _ := utils.LtreeToPath(item.Path)
	_, err = table.Update(item, parentPath)

	if err != nil {
		err = errors.Wrapf(err, "Could not update parent directory")
		return
	}

	var values []interface{}
	var vStr, kStr bytes.Buffer
	vIdx := 1
	fields := reflect.ValueOf(*newItem)
	if fields.NumField() < 1 {
		err = errors.New("Invalid number of fields given")
		return
	}
	first := true
	for i := 0; i < fields.NumField(); i++ {
		k := strings.ToLower(fields.Type().Field(i).Name)
		v := fields.Field(i).Interface()
		// ID is automatically created, don't set ID
		if k == "id" {
			continue
		}
		if TimeParam[k] { // convert time types to String
			if t, ok := v.(time.Time); ok {
				v = t.Format(time.RFC3339)
			}
		}
		if first {
			first = false
		} else {
			vStr.WriteString(", ")
			kStr.WriteString(", ")
		}
		kStr.WriteString(k)
		vStr.WriteString(fmt.Sprintf("$%d", vIdx))

		values = append(values, v)

		vIdx++
	}
	var query bytes.Buffer
	query.WriteString(fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES (%s) RETURNING *;`, FileDirsTableName, kStr.String(), vStr.String()))

	utils.Sugar.Infof("SQL Query: %s", query.String())
	utils.Sugar.Infof("Values: %v", values)

	err = pgxscan.Select(context.Background(), table.connection.Pool, &results, query.String(), values...)
	if err != nil {
		err = errors.Wrapf(err, "Insertion query failed to execute")
		return
	}

	if len(results) == 0 {
		err = errors.New(fmt.Sprintf("item at path %s not inserted", path))
		return
	}

	insertedItem = results[0]
	return
}

// Update updates an item in the database
// This update function specifically does not update
// the path
func (table *FileDirsTable) Update(updates *FileDirs, path string) (updatedItem *FileDirs, err error) {
	var result []*FileDirs

	_, err = table.GetPath(path)
	if err != nil {
		err = errors.Wrapf(err, "Update error, file does not exist")
		return
	}

	ltree, err := utils.PathToLtree(path)
	if err != nil {
		err = errors.Wrapf(err, "Could not convert path for Update")
		return
	}

	// Set item size
	if updates.Data != "" {
		updates.Size = int64(len(updates.Data))
		updates.Type = "file"
	} else {
		updates.Size = 0
		updates.Type = "directory"
	}
	// Set time
	currentTime := time.Now()
	updates.Updated = currentTime

	var query bytes.Buffer
	query.WriteString(fmt.Sprintf("UPDATE %s SET", FileDirsTableName))

	var values []interface{}
	vIdx := 1
	fields := reflect.ValueOf(*updates)
	if fields.NumField() < 1 {
		err = errors.New("Invalid number of query fields")
		return
	}
	first := true
	for i := 0; i < fields.NumField(); i++ {
		k := strings.ToLower(fields.Type().Field(i).Name)
		v := fields.Field(i).Interface()
		// Skip unset fields on the incoming item
		// Also skip the path field since we dont  want to update it
		if isUndeclared(fields.Field(i).Interface()) || SkipParam[k] || k == "path" {
			continue
		} else if TimeParam[k] { // convert time types to String
			if t, ok := v.(time.Time); ok {
				v = t.Format(time.RFC3339)
			}
		}
		if first {
			query.WriteString(" ")
			first = false
		} else {
			query.WriteString(", ")
		}

		values = append(values, v)

		query.WriteString(fmt.Sprintf("%v=$%d", k, vIdx))
		vIdx++
	}
	query.WriteString(fmt.Sprintf(" WHERE path='%s' RETURNING *;", ltree))

	utils.Sugar.Infof("SQL Query: %s", query.String())
	utils.Sugar.Infof("Values: %v", values)

	err = pgxscan.Select(context.Background(), table.connection.Pool, &result, query.String(), values...)
	if err != nil {
		err = errors.Wrapf(err, "GetPath query failed to execute")
		return
	}

	if len(result) == 0 {
		err = errors.New(fmt.Sprintf("item at path %s not updated", path))
		return
	}

	updatedItem = result[0]
	return
}

// UpdatePath moves the file at source path to be under destination
// folder
func (table *FileDirsTable) UpdatePath(sourcePath, destPath string) (updatedItems []*FileDirs, err error) {
	if sourcePath == destPath {
		err = errors.New("Can't put a folder under itself")
		return
	}

	ltreeSource, err := utils.PathToLtree(sourcePath)
	if err != nil {
		err = errors.Wrapf(err, "Could not convert source path for UpdatePath")
		return
	}

	sourceItem, err := table.GetPath(sourcePath)
	if err != nil {
		err = errors.Wrapf(err, "UpdatePath error, source path does not exist")
		return
	}

	// check for parent dir
	_, err = table.GetPath(destPath)
	if err != nil {
		err = errors.Wrapf(err, "UpdatePath error, parent path of destination does not exist")
		return
	}

	// Example: destPath = "/var", item = "ex_file", file to check = "/var/ex_file"
	destItemPath := destPath + "/" + sourceItem.Name

	destItem, err := table.GetPath(destItemPath)
	if err == nil { // Dest path exists
		switch {
		case destItem.Type == "directory":
			destPath = destItemPath //  put file under dir
		case destItem.Type == "file" && sourceItem.Type == "file":
			ltree, _ := utils.LtreeToPath(destItem.Path)
			table.DeletePath(ltree) // delete dest item and replace it with current item
		default:
			err = errors.New("Cannot replace file with directory")
			return
		}
	}

	err = nil
	// Changing destpath now
	ltreeDest, err := utils.PathToLtree(destPath)
	if err != nil {
		err = errors.Wrapf(err, "Could not convert dest path for UpdatePath")
		return
	}

	// Set time
	currentTime := time.Now().Format(time.RFC3339)
	values := []string{ltreeDest, currentTime, ltreeSource}
	// Moving the entire branch to the new folder
	query := `UPDATE filedirs SET size= (CASE WHEN type='directory' THEN 0 ELSE size END),
	updated=$1, path=$2 || subpath(path, nlevel($3)-1) where path <@ $3 RETURNING *;`

	utils.Sugar.Infof("SQL Query: %s", query)
	utils.Sugar.Infof("Values: %v", values)
	err = pgxscan.Select(context.Background(), table.connection.Pool, &updatedItems, query, currentTime, ltreeDest, ltreeSource)
	if err != nil {
		err = errors.Wrapf(err, "GetPath query failed to execute")
		return
	}

	if len(updatedItems) == 0 {
		err = errors.New(fmt.Sprintf("UpdatePath did not update any values"))
		return
	}
	return
}

// DeletePath deletes a file/dir and everything under it
func (table *FileDirsTable) DeletePath(path string) (err error) {
	ltree, err := utils.PathToLtree(path)
	if err != nil {
		err = errors.Wrapf(err, "Could not convert path for Update")
		return
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE path<@'%s';", FileDirsTableName, ltree)

	utils.Sugar.Infof("SQL Query: %s", query)

	_, err = table.connection.Pool.Exec(context.Background(), query)
	if err != nil {
		err = errors.Wrapf(err, "DeletePath query failed to execute for path %s", path)
		return
	}

	return
}
