package paxos

import (
	"fmt"
	"net/http"
	"strconv"
)

func StartLearnerController(node *Node, port string) {
	http.HandleFunc("/", learnerResponse(node))

	http.ListenAndServe(":"+port, nil)
}

func learnerResponse(n *Node) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		message := "Response from learner with id: " + strconv.Itoa(n.NodeId)
		fmt.Fprint(w, message)
	}
}
