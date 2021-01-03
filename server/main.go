package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/UN0wen/virtualfs/server/api/models"
	"github.com/UN0wen/virtualfs/server/router"
	"github.com/UN0wen/virtualfs/server/utils"
)

func main() {
	// Setup DB

	if layer := models.LayerInstance(); layer == nil {
		utils.Sugar.Fatalf("Cannot connect to database at %s:%s", utils.DBHost, utils.DBPort)
	}

	// Setup Routes
	router := router.NewRouter()

	addr := fmt.Sprintf(":%s", utils.ServerPort)
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	item := models.FileDirs{}
	item.Name = "test_file"
	item.Type = "file"
	item.Data = "poggers"

	newItem, err := models.LayerInstance().FileDirs.GetPath("/test")
	utils.CheckError(err)
	utils.Sugar.Infof("%v", newItem)

	newItem2, err := models.LayerInstance().FileDirs.GetAllPath("/", 2)
	utils.CheckError(err)
	for _, i := range newItem2 {
		utils.Sugar.Infof("%v", *i)
	}

	err = models.LayerInstance().FileDirs.Insert(&item, "/test")
	utils.CheckError(err)

	utils.Sugar.Infof("Started server on port %s", utils.ServerPort)
	utils.Sugar.Fatal(server.ListenAndServe())
}
