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

var consensusReached bool

func StartProposerController(node *Node, port string, minDelay int, maxDelay int) {
	consensusReached = false

	iteration := 0

	proposerNode := NewProposerNode(node)
	acceptorsLen := len(proposerNode.Node.Acceptors)

	http.HandleFunc("/node-info", proposerResponse(proposerNode))
	http.HandleFunc("/consensus", receiveConsensusMessage())

	go http.ListenAndServe(":"+port, nil)

	//Proposer sends requests in a loop
	for consensusReached == false {
		iteration++
		fmt.Println("Beginning new iteration")
		rand.Seed(time.Now().UnixNano())

		requestPropose := ProposeStruct{proposerNode.N, proposerNode.V, iteration}
		requestBody, _ := json.Marshal(requestPropose)

		//Sends prepare request to each acceptor
		for _, address := range proposerNode.Node.Acceptors {
			time.Sleep(time.Duration((rand.Intn(maxDelay-minDelay+1)+minDelay)*100) * time.Millisecond)
			fmt.Printf("Sending preparing proposal № %d with value %d to %s\n", proposerNode.N, proposerNode.V, address)
			resp, err := http.Post("http://"+address+"/prepare", "application/json", bytes.NewBuffer(requestBody))
			if resp.StatusCode == 200 && err == nil {
				respBody, _ := ioutil.ReadAll(resp.Body)

				var prepareResp ProposeStruct
				json.Unmarshal(respBody, &prepareResp)

				//Appends a response from acceptor to the list of responses
				proposerNode.Accepts = append(proposerNode.Accepts, &prepareResp)
			}
		}

		//If acceptors quorum responded to the proposer, searches for the proposal with the highest number
		if len(proposerNode.Accepts) >= (acceptorsLen/2 + 1) {
			i := 0
			var val int

			for _, p := range proposerNode.Accepts {
				if p.N > i {
					i = p.N
					val = p.V
				}
			}

			if i != 0 {
				proposerNode.V = val
			}

			requestPropose = ProposeStruct{proposerNode.N, proposerNode.V, iteration}
			requestBody, _ = json.Marshal(requestPropose)

			//Sends accept request to each acceptor
			for _, address := range proposerNode.Node.Acceptors {
				time.Sleep(time.Duration((rand.Intn(maxDelay-minDelay+1)+minDelay)*100) * time.Millisecond)
				fmt.Printf("Sending acceptance proposal № %d with value %d to %s\n", proposerNode.N, proposerNode.V, address)
				http.Post("http://"+address+"/accept", "application/json", bytes.NewBuffer(requestBody))
			}
		}

		//The list of acceptors responses is cleared
		proposerNode.Accepts = nil

		//Updating proposal number
		proposerNode.N += len(proposerNode.Node.Proposers)
	}
}

func proposerResponse(n *ProposerNode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respBody, _ := json.Marshal(n)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(respBody)
	}
}

func receiveConsensusMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		consensusReached = true
		w.WriteHeader(200)
	}
}
