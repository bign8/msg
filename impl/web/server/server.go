package server

import (
	"net/http"

	"github.com/bign8/msg"
)

// New constructs a new web based service handler.
// Register to a fixed HTTP location (be sure to use http.StripPrefix).
func New(svc *msg.Server) http.Handler {
	s := &server{
		ServeMux: http.NewServeMux(),
		svc:      svc,
	}
	s.ServeMux.HandleFunc("/sock", s.sock)  // Sample: GET /sock
	s.ServeMux.HandleFunc("/http/", s.http) // Sample: POST /http/:service/:method
	s.ServeMux.HandleFunc("/load/", s.load) // Sample: GET /load/:assetID
	return s
}

type server struct {
	*http.ServeMux
	svc *msg.Server
}

// sock creates and manages a websocket, designed for small messages
func (s *server) sock(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "TODO", http.StatusNotImplemented)
}

// http deals with  single method invocation RPC, designed for large messages
func (s *server) http(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "TODO", http.StatusNotImplemented)
}

// load returns a blob, this is used when an RPC was invoked via socket,
// but the resposne is too large to send across a websocket
func (s *server) load(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "TODO", http.StatusNotImplemented)
}
