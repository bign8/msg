package blob

import (
	"net/http"

	"github.com/bign8/msg"
	"github.com/bign8/msg/rand"
)

// NewServer constructs a new blob server
func NewServer(conn *msg.Conn) *Server {
	return &Server{
		rand: rand.New(),
		conn: conn,
		data: make(map[string][]byte, 1),
	}
}

// Server is a blob server
type Server struct {
	rand func() string
	conn *msg.Conn
	data map[string][]byte // TODO: TTL
}

// Start actually executes the server
func (s *Server) Start(ip string) (err error) {
	err = s.conn.Publish(nil, &msg.Msg{
		Title: "sd_register",        // service discovery
		Body:  []byte("blob:" + ip), // I am the blob service
	})
	if err != nil {
		return err
	}
	http.HandleFunc("/save", s.save)
	http.HandleFunc("/load/", s.load)
	return http.ListenAndServe(ip, nil)
}

func (s *Server) save(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "todo", http.StatusNotImplemented)
}

func (s *Server) load(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/load/"):]
	bits, ok := s.data[id]
	if !ok {
		http.NotFound(w, r)
		return
	}
	w.Write(bits)
}
