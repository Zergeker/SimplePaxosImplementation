package paxos

import (
	"fmt"
	"net/http"
	"strconv"
)

func StartAcceptorController(node *Node, port string) {
	http.HandleFunc("/", acceptorResponse(node))

	http.ListenAndServe(":"+port, nil)
}

func acceptorResponse(n *Node) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		message := "Response from acceptor with id: " + strconv.Itoa(n.NodeId)
		fmt.Fprint(w, message)
	}
}
