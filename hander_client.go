package bjorno

import (
	"encoding/json"
	"net/http"

	"github.com/kris-nova/bjorn/lib"
	"github.com/kris-nova/logger"
)

// ClientHandler is for the nivenly client API
type ClientHandler struct {
	Config  *ServerConfig
	HTTPDir http.Dir
}

func NewClientHandler(cfg *ServerConfig) *ClientHandler {
	return &ClientHandler{
		Config:  cfg,
		HTTPDir: http.Dir(cfg.ServeDirectory),
	}
}

// ServeHTTP is where the magic happens.
func (rh *ClientHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	clientMeta := lib.GetClientMeta(r)
	bytes, err := json.Marshal(clientMeta)
	if err != nil {
		logger.Warning(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(rh.Config.Content500)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
