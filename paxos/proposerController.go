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

func StartProposerController(node *Node, port string, minDelay int, maxDelay int) {
	proposerNode := NewProposerNode(node)
	acceptorsLen := len(proposerNode.Node.Acceptors)

	http.HandleFunc("/node-info", proposerResponse(proposerNode))

	go http.ListenAndServe(":"+port, nil)

	//Proposer sends requests in a loop
	for true == true {
		fmt.Println("Beginning new iteration")
		rand.Seed(time.Now().UnixNano())

		requestPropose := ProposeStruct{proposerNode.N, proposerNode.V}
		requestBody, _ := json.Marshal(requestPropose)

		//Sends prepare request to each acceptor
		for _, address := range proposerNode.Node.Acceptors {
			time.Sleep(time.Duration((rand.Intn(maxDelay-minDelay+1)+minDelay)*100) * time.Millisecond)
			fmt.Printf("Sending preparing proposal № %d with value %d to %s\n", proposerNode.N, proposerNode.V, address)
			resp, _ := http.Post("http://"+address+"/prepare", "application/json", bytes.NewBuffer(requestBody))
			if resp.StatusCode == 200 {
				respBody, _ := ioutil.ReadAll(resp.Body)

				var prepareResp ProposeStruct
				json.Unmarshal(respBody, &prepareResp)

				fmt.Printf("Prepare response: N=%d V=%d\n", prepareResp.N, prepareResp.V)

				//Appends a response from acceptor to the list of responses
				proposerNode.Accepts = append(proposerNode.Accepts, &prepareResp)
			}
		}

		//If acceptors quorum responded to the proposer, searches for the proposal with the highest number
		if len(proposerNode.Accepts) >= (acceptorsLen/2 + 1) {
			for _, p := range proposerNode.Accepts {
				if p.N > proposerNode.N {
					fmt.Printf("Changing the proposers value %d to the value %d of a greater-numbered proposal № %d\n", proposerNode.V, p.V, p.N)
					proposerNode.V = p.V
				}
			}

			requestPropose = ProposeStruct{proposerNode.N, proposerNode.V}
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

		//Updating proposal number and new value generation
		proposerNode.N += len(proposerNode.Node.Proposers)
		proposerNode.V = rand.Intn(100)
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
