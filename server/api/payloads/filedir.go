package payloads

import (
	"errors"
	"net/http"

	"github.com/UN0wen/virtualfs/server/api/models"
	"github.com/UN0wen/virtualfs/server/utils"
	"github.com/go-chi/render"
)

// FileDirsRequest is the request payload for the FileDirs data model
type FileDirsRequest struct {
	FileDirs *models.FileDirs `json:"filedirs"`

	CurrentWorkDir string `json:"cwd"` // the current working directory in the clientside

	Path string `json:"path"`
}

// Bind is the postprocessing for the FileDirsRequest after the request is unmarshalled
func (fd *FileDirsRequest) Bind(r *http.Request) error {
	if fd.FileDirs == nil {
		return errors.New("missing required FileDirs fields")
	}

	if !utils.IsAbsolute(fd.Path) {
		return errors.New("Path must be present and must be absolute")
	}
	return nil
}

// MoveRequest is the request payload for move requests
type MoveRequest struct {
	Source string `json:"source"`
	Dest   string `json:"dest"`
}

// Bind is the postprocessing for the MoveRequest after the request is unmarshalled
func (mr *MoveRequest) Bind(r *http.Request) error {
	if !utils.IsAbsolute(mr.Source) {
		return errors.New("Source path must be present and must be absolute")
	}

	if !utils.IsAbsolute(mr.Dest) {
		return errors.New("Destination path must be present and must be absolute")
	}
	return nil
}

// FileDirsResponse is the response payload for the FileDirs data model.
type FileDirsResponse struct {
	FileDirs *models.FileDirs `json:"filedirs"`
}

// NewFileDirsResponse generate a Response for FileDirs object
func NewFileDirsResponse(filedir *models.FileDirs) *FileDirsResponse {
	resp := &FileDirsResponse{FileDirs: filedir}

	return resp
}

// NewFileDirsListResponse generates a list of renders for Items
func NewFileDirsListResponse(items []*models.FileDirs) []render.Renderer {
	list := []render.Renderer{}
	for i := range items {
		list = append(list, NewFileDirsResponse(items[i]))
	}

	return list
}

// Render is preprocessing before the response is marshalled
func (rd *FileDirsResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}
