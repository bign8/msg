package main

import "net/http"

func main() {
	http.Handle("/sock", nil) // ws, used for small connections
	http.Handle("/port", nil) // HTTP post, used for larger blobs
	http.Handle("/serv", nil) // Serve large static assets

	// A connection will be pointed at a server, serving ^ 2 endpoints
	// will swap between ws for small RPC and post for large packets
	// servers can do pub-sub over sock
	// iff pub-sub payload it too large to send over sock:
	//     serve file locally and publish host url
}
