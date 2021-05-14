package bjorno

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/kris-nova/logger"
)

// RootHandler is a custom server that proxies whatever HTTPDir is set to to
// the / (root) of the HTTP(s) server.
type RootHandler struct {
	Config  *ServerConfig
	HTTPDir http.Dir
}

func NewRootHandler(cfg *ServerConfig) *RootHandler {
	return &RootHandler{
		Config:  cfg,
		HTTPDir: http.Dir(cfg.ServeDirectory),
	}
}

// RequestPath is a deterministic function that
// given an *http.Request will always return a
// request path to "search" for.
//
// Note: this will *NOT* check an inode for a directory (isDir())
// but will trust the request to identify a directory
// by POSIX convention.
//
// If the path ends with "/" it's a directory...
//
func RequestPath(r *http.Request) string {
	// Getting dot dot right
	requestPath := path.Clean(r.URL.Path)
	// Remove all periods "."
	requestPath = strings.Replace(requestPath, ".", "", -1)
	if !strings.HasPrefix(requestPath, "/") {
		requestPath = fmt.Sprintf("/%s", requestPath)
	}
	return requestPath

}

// FileDirectoryPath will take a set of default file strings, a request path, and a valid http.Dir and
// handles the logic for checking the filesystem for default files in directories such as index.html
//
// This will only return a file and stat if the file calculated actually exists.
func FileDirectoryPath(defaultFiles []string, requestPath string, httpDir http.Dir) (http.File, os.FileInfo, error) {
	var file http.File
	file, err := httpDir.Open(requestPath)
	if err != nil {
		return file, nil, fmt.Errorf("unable to open file %s: %v", requestPath, err)
	}
	stat, err := file.Stat()
	if err != nil {
		return file, nil, fmt.Errorf("unable to stat file %s: %v", requestPath, err)
	}
	if stat.IsDir() {
		// Very important to loop literally so that we will always return
		// the first file in the default files list found.
		for i := 0; i < len(defaultFiles); i++ {
			checkFile := defaultFiles[i]
			// Attempt to open file
			file, err := httpDir.Open(path.Join(requestPath, checkFile))
			if file != nil {
				// Found so let's break
				stat, err = file.Stat()
				if err != nil {
					return file, stat, fmt.Errorf("unable to stat file: %s", err)
				}
				return file, stat, nil
			}
			logger.Debug(err.Error())
		}
		// If we get here we have a directory but no default files are found
		return file, stat, fmt.Errorf("default file not found in directory: %s", requestPath)
	}
	if file == nil {
		return file, stat, fmt.Errorf("unable to find file or default file in list")
	}
	return file, stat, nil
}

func (rh *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// System to calculate "RequestPath"
	requestPath := RequestPath(r)
	// System to hit the filesystem and calculate default files in dirs
	file, stat, err := FileDirectoryPath(rh.Config.DefaultIndexFiles, requestPath, rh.HTTPDir)
	if err != nil {
		// 404
		w.WriteHeader(http.StatusNotFound)
		w.Write(rh.Config.Content404)
		return
	}
	logger.Info("Request: %s", stat.Name())
	interpolatedFile := InterpolateFile(file)
	http.ServeContent(w, r, stat.Name(), stat.ModTime(), interpolatedFile)
}

func InterpolateFile(file http.File) http.File {
	// magical boops here
	return file
}
