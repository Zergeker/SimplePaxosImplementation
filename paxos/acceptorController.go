package paxos

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func StartAcceptorController(node *Node, port string) {
	acceptorNode := NewAcceptorNode(node)

	http.HandleFunc("/", acceptorResponse(acceptorNode))
	http.HandleFunc("/prepare", receivePrepare(acceptorNode))
	http.HandleFunc("/accept", receiveAccept(acceptorNode))

	http.ListenAndServe(":"+port, nil)
}

func acceptorResponse(n *AcceptorNode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		message := "Response from acceptor with id: " + strconv.Itoa(n.Node.NodeId) + "\n propose number: " +
			strconv.Itoa(n.N) + "\n propose value: " + strconv.Itoa(n.V)
		fmt.Fprint(w, message)
	}
}

func receivePrepare(acceptor *AcceptorNode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestBody, _ := ioutil.ReadAll(r.Body)

		var proposeStruct ProposeStruct
		json.Unmarshal(requestBody, &proposeStruct)

		if proposeStruct.N > acceptor.N {
			acceptor.N = proposeStruct.N
			acceptor.V = proposeStruct.V
		}

		responseBodyStruct := ProposeStruct{acceptor.N, acceptor.V}

		respBody, _ := json.Marshal(responseBodyStruct)

		fmt.Fprintf(w, "%s", respBody)
	}
}

func receiveAccept(acceptor *AcceptorNode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestBody, _ := ioutil.ReadAll(r.Body)

		var proposeStruct ProposeStruct
		json.Unmarshal(requestBody, &proposeStruct)

		if proposeStruct.N >= acceptor.N {
			acceptor.N = proposeStruct.N
			acceptor.V = proposeStruct.V
			fmt.Printf("New value %d has been accepted with propose № %d\n", acceptor.V, acceptor.N)
		} else {
			fmt.Printf("Propose № %d has been rejected, current N: %d\n", proposeStruct.N, acceptor.N)
		}

		w.WriteHeader(200)
	}
}
