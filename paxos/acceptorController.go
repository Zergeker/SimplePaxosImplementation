package paxos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

func StartAcceptorController(node *Node, port string, minDelay int, maxDelay int) {
	acceptorNode := NewAcceptorNode(node)

	http.HandleFunc("/node-info", respondInfo(acceptorNode))
	http.HandleFunc("/prepare", receivePrepare(acceptorNode))
	http.HandleFunc("/accept", receiveAcceptAcceptor(acceptorNode, minDelay, maxDelay))

	http.ListenAndServe(":"+port, nil)
}

func respondInfo(n *AcceptorNode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		acceptorInfo := AcceptorInfo{n.Node.NodeId, n.N, n.V}
		respBody, _ := json.Marshal(acceptorInfo)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(respBody)
	}
}

func receivePrepare(acceptor *AcceptorNode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestBody, _ := ioutil.ReadAll(r.Body)

		var proposeStruct ProposeStruct
		json.Unmarshal(requestBody, &proposeStruct)

		//Acceptor changes its proposal number and value if receives a proposal with a higher number
		if proposeStruct.N > acceptor.N {
			acceptor.N = proposeStruct.N
			acceptor.V = proposeStruct.V

			fmt.Printf("New value %d has been prepared for accept with propose № %d\n", acceptor.V, acceptor.N)
		} else {
			fmt.Printf("Proposal № %d has been rejected, current №: %d", proposeStruct.N, acceptor.N)
		}

		//Acceptor responds with its current proposal number and value
		responseBodyStruct := ProposeStruct{acceptor.N, acceptor.V}
		respBody, _ := json.Marshal(responseBodyStruct)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(respBody)
	}
}

func receiveAcceptAcceptor(acceptor *AcceptorNode, minDelay int, maxDelay int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestBody, _ := ioutil.ReadAll(r.Body)

		var proposeStruct ProposeStruct
		json.Unmarshal(requestBody, &proposeStruct)

		//Acceptor changes its proposal number and value if it receives a proposal with a higher or equal number
		if proposeStruct.N >= acceptor.N {
			acceptor.N = proposeStruct.N
			acceptor.V = proposeStruct.V

			acceptRequest := AcceptorInfo{acceptor.Node.NodeId, proposeStruct.N, proposeStruct.V}
			reqBodyAccept, _ := json.Marshal(acceptRequest)

			//Since it's an accept request, sends current proposal to each learner
			for _, address := range acceptor.Node.Learners {
				time.Sleep(time.Duration((rand.Intn(maxDelay-minDelay+1)+minDelay)*100) * time.Millisecond)
				http.Post("http://"+address+"/accept", "application/json", bytes.NewBuffer(reqBodyAccept))
			}

			fmt.Printf("New value %d has been accepted with propose № %d\n", acceptor.V, acceptor.N)
		} else {
			fmt.Printf("Propose № %d has been rejected, current №: %d\n", proposeStruct.N, acceptor.N)
		}

		w.WriteHeader(200)
	}
}
