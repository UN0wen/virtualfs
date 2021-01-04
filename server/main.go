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

	utils.Sugar.Infof("Started server on port %s", utils.ServerPort)
	utils.Sugar.Fatal(server.ListenAndServe())
}
