package bjorno

import (
	"fmt"
	"github.com/kris-nova/logger"
	"net/http"
	"os"
	"path"
	"strings"
)

// RootHandler is a custom server that proxies whatever HTTPDir is set to to
// the / (root) of the HTTP(s) server.
type RootHandler struct {
	Config *ServerConfig
	HTTPDir http.Dir
}

func NewRootHandler(cfg *ServerConfig) *RootHandler {
	return &RootHandler{
		Config: cfg,
		HTTPDir: http.Dir(cfg.ServeDirectory),
	}
}

// ServeHTTP is where the magic happens.
//
// Here is where the custom logic begins
// We build a "custom" rootHandler that implements
// http.Handle so that we can add our own logic to the
// base HTTP server.
// This allows us to intercept new HTTP requests to the
// server and do whatever we want to with them.
//
// The whole value of Bjorno is built on this handler's
// ability to work "well"
func (rh *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// System to mutate the request path
	requestPath := r.URL.Path
	if !strings.HasPrefix(requestPath, "/"){
		requestPath = fmt.Sprintf("/%s", requestPath)
	}
	requestPath = path.Clean(r.URL.Path)
	if requestPath == "/" {
		// Special logic for root /
		requestPath = "/index.html"
		r.URL.Path = "index.html"
	}

	// Attempt to open the requested path
	file, err := rh.HTTPDir.Open(requestPath)
	//_, err := rh.HTTPDir.Open(requestPath)
	// Error reading requested file
	if err != nil {
		if os.IsNotExist(err) {
			// 404 File not found
			logger.Critical("404 looking up mutated path (%s) original path (%s)", requestPath, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			w.Write(rh.Config.Content404)
			return
		}else {
			// 500 Internal server
			logger.Warning("500 %v", err)
			//
			// @kris-nova
			//
			// We should do a lot more when this happens...
			//
			//
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(rh.Config.Content500)
			return
		}
	}else {
		stat, err := file.Stat()
		if err != nil {
			logger.Warning(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(rh.Config.Content500)
			return
		}
		// -------------------------------------
		interpolatedFile := InterpolateFile(file)
		http.ServeContent(w, r, stat.Name(), stat.ModTime(), interpolatedFile)
		// -------------------------------------
	}
}

func InterpolateFile(file http.File) http.File {
	// magical boops here
	return file
}


