package paxos

import (
	"fmt"
	"net/http"
	"strconv"
)

func StartProposerController(node *Node, port string) {
	http.HandleFunc("/", proposerResponse(node))

	http.ListenAndServe(":"+port, nil)
}

func proposerResponse(n *Node) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		message := "Response from proposer with id: " + strconv.Itoa(n.NodeId)
		fmt.Fprint(w, message)
	}
}
