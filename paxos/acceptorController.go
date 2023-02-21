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

var msgCount int
var crashed bool
var minDelay int
var maxDelay int

func StartAcceptorController(node *Node, port string, minDel int, maxDel int) {
	minDelay = minDel
	maxDelay = maxDel

	msgCount = 0
	crashed = false

	acceptorNode := NewAcceptorNode(node)

	http.HandleFunc("/node-info", respondInfo(acceptorNode))
	http.HandleFunc("/prepare", receivePrepare(acceptorNode))
	http.HandleFunc("/accept", receiveAcceptAcceptor(acceptorNode, minDelay, maxDelay))
	http.HandleFunc("/recover", recoverNode())

	http.ListenAndServe(":"+port, nil)
}

func respondInfo(n *AcceptorNode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		acceptorInfo := AcceptorInfo{n.Node.NodeId, n.AcceptedN, n.AcceptedN, 1}
		respBody, _ := json.Marshal(acceptorInfo)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(respBody)
	}
}

func receivePrepare(acceptor *AcceptorNode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if crashed == false {
			//Incrementing the number of received messages, crashing the node when it overcomes the limit
			msgCount++
			if msgCount > acceptor.Node.msgLimit {
				crashed = true
				fmt.Println("Simulating crash state")
				w.WriteHeader(500)
			}

			requestBody, _ := ioutil.ReadAll(r.Body)

			var proposeStruct ProposeStruct
			json.Unmarshal(requestBody, &proposeStruct)

			//Acceptor changes its proposal number and value if receives a proposal with a higher number
			if proposeStruct.N > acceptor.PreparedN {
				acceptor.PreparedN = proposeStruct.N
				acceptor.PreparedV = proposeStruct.V

				//Acceptor responds with its last accepted proposal number and value
				responseBodyStruct := ProposeStruct{acceptor.AcceptedN, acceptor.AcceptedV, 1}
				respBody, _ := json.Marshal(responseBodyStruct)

				time.Sleep(time.Duration((rand.Intn(maxDelay-minDelay+1)+minDelay)*100) * time.Millisecond)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(200)
				w.Write(respBody)

				fmt.Printf("New value %d has been prepared for accept with propose № %d\n", acceptor.PreparedV, acceptor.PreparedN)
			} else {
				time.Sleep(time.Duration((rand.Intn(maxDelay-minDelay+1)+minDelay)*100) * time.Millisecond)
				w.WriteHeader(500)
				fmt.Printf("Proposal № %d has been rejected, current №: %d\n", proposeStruct.N, acceptor.PreparedN)
			}
		} else {
			w.WriteHeader(500)
		}
	}
}

func receiveAcceptAcceptor(acceptor *AcceptorNode, minDelay int, maxDelay int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if crashed == false {
			//Incrementing the number of received messages, crashing the node when it overcomes the limit
			msgCount++
			if msgCount > acceptor.Node.msgLimit {
				crashed = true
				fmt.Println("Simulating crash state")
				w.WriteHeader(500)
			}

			requestBody, _ := ioutil.ReadAll(r.Body)

			var proposeStruct ProposeStruct
			json.Unmarshal(requestBody, &proposeStruct)

			//Acceptor changes its last accepted proposal number and value if it receives a proposal with a higher or equal number
			if proposeStruct.N >= acceptor.AcceptedN {
				acceptor.AcceptedN = proposeStruct.N
				acceptor.AcceptedV = proposeStruct.V

				acceptRequest := AcceptorInfo{acceptor.Node.NodeId, proposeStruct.N, proposeStruct.V, proposeStruct.I}
				reqBodyAccept, _ := json.Marshal(acceptRequest)

				//Since it's an accept request, sends current proposal to each learner
				for _, address := range acceptor.Node.Learners {
					time.Sleep(time.Duration((rand.Intn(maxDelay-minDelay+1)+minDelay)*100) * time.Millisecond)
					http.Post("http://"+address+"/accept", "application/json", bytes.NewBuffer(reqBodyAccept))
				}

				fmt.Printf("New value %d has been accepted with propose № %d\n", acceptor.AcceptedV, acceptor.AcceptedN)
			} else {
				fmt.Printf("Propose № %d has been rejected, current №: %d\n", proposeStruct.N, acceptor.AcceptedN)
			}
			time.Sleep(time.Duration((rand.Intn(maxDelay-minDelay+1)+minDelay)*100) * time.Millisecond)
			w.WriteHeader(200)
		} else {
			w.WriteHeader(500)
		}
	}
}

func recoverNode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if crashed == false {
			w.Write([]byte("This node is not crashed"))
		} else {
			crashed = false
			fmt.Println("Recovering from crash state")
			w.Write([]byte("Node recovered"))
		}
	}
}
