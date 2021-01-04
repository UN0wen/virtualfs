package controllers

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strconv"

	"github.com/UN0wen/virtualfs/server/api/models"
	"github.com/UN0wen/virtualfs/server/api/payloads"
	"github.com/UN0wen/virtualfs/server/utils"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// GetPath returns all files and dirs for a path, at a specific level.
// It accepts two optional queries
// path: string is the source path
// level: int is the number of levels deep to get
func GetPath(w http.ResponseWriter, r *http.Request) {
	var path string
	var level int
	path = chi.URLParam(r, "path")
	level, err := strconv.Atoi(r.URL.Query().Get("level"))

	// Can't parse level, just use default
	if err != nil {
		level = 1
		err = nil
	}
	parsedPath, err := base64.URLEncoding.DecodeString(path)

	if err != nil {
		render.Render(w, r, payloads.ErrInvalidRequest(errors.New("Path is not correctly encoded")))
		return
	}
	parsedPathString := string(parsedPath)
	// Make sure that path is absolute
	if parsedPathString != "" {
		if !utils.IsAbsolute(parsedPathString) {
			render.Render(w, r, payloads.ErrInvalidRequest(errors.New("Provided paths must be absolute")))
		}
	} else {
		parsedPathString = "/"
	}

	fileDirs, err := models.LayerInstance().FileDirs.GetAllPath(parsedPathString, level)

	if err != nil {
		render.Render(w, r, payloads.ErrInternalError(err))
		return
	}

	if err := render.RenderList(w, r, payloads.NewFileDirsListResponse(fileDirs)); err != nil {
		render.Render(w, r, payloads.ErrRender(err))
		return
	}
}

// GetExactPath returns a single file at a path.
// It uses a URLParam of path encoded in base64
func GetExactPath(w http.ResponseWriter, r *http.Request) {
	var path string
	path = chi.URLParam(r, "path")

	parsedPath, err := base64.URLEncoding.DecodeString(path)

	if err != nil {
		render.Render(w, r, payloads.ErrInvalidRequest(errors.New("Path is not correctly encoded")))
		return
	}

	parsedPathString := string(parsedPath)
	// Make sure that path is absolute
	if parsedPathString != "" {
		if !utils.IsAbsolute(parsedPathString) {
			render.Render(w, r, payloads.ErrInvalidRequest(errors.New("Provided paths must be absolute")))
		}
	} else {
		parsedPathString = "/"
	}

	fileDir, err := models.LayerInstance().FileDirs.GetPath(parsedPathString)

	if err != nil {
		render.Render(w, r, payloads.ErrInternalError(err))
		return
	}

	if err := render.Render(w, r, payloads.NewFileDirsResponse(fileDir)); err != nil {
		render.Render(w, r, payloads.ErrRender(err))
		return
	}
}

// CreatePath creates a new item at a specific place
func CreatePath(w http.ResponseWriter, r *http.Request) {
	data := &payloads.FileDirsRequest{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, payloads.ErrInvalidRequest(err))
		return
	}

	fileDirs := data.FileDirs

	inserted, err := models.LayerInstance().FileDirs.Insert(fileDirs, data.Path)

	if err != nil {
		render.Render(w, r, payloads.ErrInternalError(err))
		return
	}

	if err := render.Render(w, r, payloads.NewFileDirsResponse(inserted)); err != nil {
		render.Render(w, r, payloads.ErrRender(err))
		return
	}
}

// UpdateItem updates an item's value at a specific path
func UpdateItem(w http.ResponseWriter, r *http.Request) {
	data := &payloads.FileDirsRequest{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, payloads.ErrInvalidRequest(err))
		return
	}

	fileDirs := data.FileDirs

	updated, err := models.LayerInstance().FileDirs.Update(fileDirs, data.Path)

	if err != nil {
		render.Render(w, r, payloads.ErrInternalError(err))
		return
	}

	if err := render.Render(w, r, payloads.NewFileDirsResponse(updated)); err != nil {
		render.Render(w, r, payloads.ErrRender(err))
		return
	}
}

// UpdatePath moves an item from source to dest and returns all items changed
func UpdatePath(w http.ResponseWriter, r *http.Request) {
	data := &payloads.MoveRequest{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, payloads.ErrInvalidRequest(err))
		return
	}

	updated, err := models.LayerInstance().FileDirs.UpdatePath(data.Source, data.Dest)

	if err != nil {
		render.Render(w, r, payloads.ErrInternalError(err))
		return
	}

	if err := render.RenderList(w, r, payloads.NewFileDirsListResponse(updated)); err != nil {
		render.Render(w, r, payloads.ErrRender(err))
		return
	}
}

// DeletePath deletes an item from the path
// and all of its children
func DeletePath(w http.ResponseWriter, r *http.Request) {
	path := chi.URLParam(r, "path")
	parsedPath, err := base64.URLEncoding.DecodeString(path)

	if err != nil {
		render.Render(w, r, payloads.ErrInvalidRequest(errors.New("Path is not correctly encoded")))
		return
	}
	parsedPathString := string(parsedPath)

	numRows, err := models.LayerInstance().FileDirs.DeletePath(parsedPathString)
	if err != nil {
		render.Render(w, r, payloads.ErrInternalError(err))
		return
	}

	if err := render.Render(w, r, payloads.NewDeleteResponse(numRows)); err != nil {
		render.Render(w, r, payloads.ErrRender(err))
		return
	}
}
