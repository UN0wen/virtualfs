package payloads

import (
	"net/http"

	"github.com/go-chi/render"
)

//--
// Error response payloads & renderers
//--

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

// Render is required to implement the Renderer interface
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// ErrInvalidRequest is a response payload with status code 400.
func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

// ErrRender is a response payload with status code 422.
func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

// ErrUnauthorized is a response payload with status code 401.
func ErrUnauthorized(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 401,
		StatusText:     "Unauthorized access",
		ErrorText:      err.Error(),
	}
}

// ErrInternalError is a response payload with status code 500.
func ErrInternalError(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 500,
		StatusText:     "Internal server error",
		ErrorText:      err.Error(),
	}
}

// ErrNotImplemented is a response payload with status code 501.
// This error is for creating items without the corresponding scraper
var ErrNotImplemented = &ErrResponse{HTTPStatusCode: 501, StatusText: "This website is not supported."}

// ErrNotFound is a standard response for a 404 code
var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}
