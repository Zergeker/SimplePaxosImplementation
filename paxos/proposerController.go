package paxos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func StartProposerController(node *Node, port string) {
	proposerNode := NewProposerNode(node)
	acceptorsLen := len(proposerNode.Node.Acceptors)

	http.HandleFunc("/", proposerResponse(proposerNode))

	go http.ListenAndServe(":"+port, nil)

	for true == true {
		rand.Seed(time.Now().UnixNano())

		proposerNode.N = proposerNode.Node.NodeId
		proposerNode.V = rand.Intn(100)

		requestPropose := ProposeStruct{proposerNode.N, proposerNode.V}
		requestBody, _ := json.Marshal(requestPropose)

		for _, address := range proposerNode.Node.Acceptors {

			resp, _ := http.Post("http://"+address+"/prepare", "application/json", bytes.NewBuffer(requestBody))
			respBody, _ := ioutil.ReadAll(resp.Body)
			var prepareResp ProposeStruct
			json.Unmarshal(respBody, &prepareResp)

			proposerNode.Accepts = append(proposerNode.Accepts, &prepareResp)
		}

		if len(proposerNode.Accepts) > (acceptorsLen/2 + 1) {
			for _, p := range proposerNode.Accepts {
				if p.N > proposerNode.N {
					proposerNode.V = p.V
				}
			}

			requestPropose = ProposeStruct{proposerNode.N, proposerNode.V}
			requestBody, _ = json.Marshal(requestPropose)

			for _, address := range proposerNode.Node.Acceptors {
				http.Post("http://"+address+"/accept", "application/json", bytes.NewBuffer(requestBody))
			}
		}

		time.Sleep(time.Duration((rand.Intn(4)+7)*100) * time.Millisecond)
		proposerNode.Accepts = nil
		proposerNode.Node.NodeId += len(proposerNode.Node.Proposers)
	}
}

func proposerResponse(n *ProposerNode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		message := "Response from proposer with id: " + strconv.Itoa(n.Node.NodeId)
		fmt.Fprint(w, message)
	}
}
