package models

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/UN0wen/virtualfs/server/utils"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/google/go-cmp/cmp"
	_ "github.com/jackc/pgx/stdlib"
)

var (
	fixtures *testfixtures.Loader
)

func TestMain(m *testing.M) {
	// connect to db
	cfgDNS := fmt.Sprintf(
		"user=%s password=%s dbname=virtualfs_test host=%s port=%s sslmode=disable",
		utils.DBUser, utils.DBPassword, utils.DBHost, utils.DBPort)
	// Create all the tables

	db, err := sql.Open("pgx", cfgDNS)
	if err != nil {
		utils.Sugar.Fatal(err)
	}
	fixtures, err = testfixtures.New(testfixtures.Database(db), testfixtures.Dialect("postgresql"), testfixtures.Directory("fixtures"))
	if err != nil {
		utils.Sugar.Fatal(err)
	}
	os.Exit(m.Run())
}

func prepareTestDatabase() {
	if err := fixtures.Load(); err != nil {
		utils.Sugar.Fatal(err)
	}
}
func TestGet(t *testing.T) {
	prepareTestDatabase()

	currentTime := time.Now()
	expected := FileDirs{
		ID:      5,
		Name:    "mv",
		Data:    "move func",
		Size:    9,
		Type:    "file",
		Path:    "root.bin.mv",
		Created: currentTime,
		Updated: currentTime,
	}

	item, err := LayerInstance().FileDirs.GetPath("/bin/mv")
	if err != nil {
		t.Fatal(err)
	}

	// unset the time fields
	item.Created = currentTime
	item.Updated = currentTime

	if !cmp.Equal(*item, expected) {
		t.Errorf("GET(/bin/mv) = %v; want %v", item, expected)
	}
}

func TestCreate(t *testing.T) {
	prepareTestDatabase()

	expected := FileDirs{
		ID:   10001, // testfixtures start sequences at 10000
		Name: "testdir",
		Size: 0,
		Type: "directory",
		Path: "root.lib.testdir",
	}

	inserted, err := LayerInstance().FileDirs.Insert(&expected, "/lib")
	if err != nil {
		t.Fatal(err)
	}

	// unset the time fields
	inserted.Created = time.Time{}
	inserted.Updated = time.Time{}
	expected.Created = time.Time{}
	expected.Updated = time.Time{}

	if !cmp.Equal(*inserted, expected) {
		t.Errorf("INSERT(/bin/mv) = %v; want %v", inserted, expected)
	}
}

func TestUpdatePath(t *testing.T) {
	// /var -> /lib/var
	updated, err := LayerInstance().FileDirs.UpdatePath("/var", "/lib")
	if err != nil {
		t.Fatal(err)
	}

	if len(updated) != 4 {
		t.Errorf("UPDATEPATH(/var, /lib) = %d changes; want 4", len(updated))
	}
}

func TestDelete(t *testing.T) {
	// /var -> /lib/var
	_, err := LayerInstance().FileDirs.DeletePath("/lib")
	if err != nil {
		t.Fatal(err)
	}

	// Check that this deletes all path beneath var

	items, err := LayerInstance().FileDirs.GetAllPath("/var", 5)
	if err == nil && len(items) > 0 {
		t.Errorf("GETALL(/var) after delete returns values; want error")
	}
}
