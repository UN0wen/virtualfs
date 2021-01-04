package utils

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// IsAbsolute checks if the provided path is absolute
func IsAbsolute(path string) (abs bool) {
	return strings.HasPrefix(path, "/")
}

// CleanPath gives us a clean path from a potentially relative path
func CleanPath(path string, currentWorkDir string) (absolutePath string, err error) {
	if strings.HasPrefix(path, "/") {
		absolutePath = path
		return
	}

	absolutePath, err = filepath.Abs(currentWorkDir + path)
	if err != nil {
		errors.Wrap(err, "Could not parse path")
	}

	return
}

// PathToLtree returns a Ltree representation of an absolute path
func PathToLtree(path string) (ltree string, err error) {
	if !strings.HasPrefix(path, "/") {
		err = errors.New("Relative paths cannot be converted to Ltree")
		return
	}

	ltree = "root"

	if path == "/" {
		return
	}

	ltree = ltree + strings.ReplaceAll(path, "/", ".")
	return
}

// LtreeToPath returns a fs path representation of an Ltree path
func LtreeToPath(ltree string) (path string, err error) {
	if !strings.HasPrefix(ltree, "root") {
		err = errors.New("Ltree format is incorrect")
		return
	}

	if ltree == "root" {
		path = "/"
		return
	}

	noRoot := strings.Replace(ltree, "root.", "/", 1)

	path = strings.ReplaceAll(noRoot, ".", "/")
	return
}

// GetParentPath expects an absolute path with no / at the end
// and return the parent path of that path
func GetParentPath(inPath string) string {
	return path.Dir(inPath)
}
