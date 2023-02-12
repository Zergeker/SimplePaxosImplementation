package paxos

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func StartLearnerController(node *Node, port string) {
	learnerNode := NewLearnerNode(node)
	initializeAccepts(learnerNode)

	http.HandleFunc("/node-info", learnerResponse(learnerNode))
	http.HandleFunc("/accept", receiveAcceptLearner(learnerNode))

	http.ListenAndServe(":"+port, nil)
}

func initializeAccepts(node *LearnerNode) {
	//Getting each acceptors info to form a list of acceptors proposals
	for _, address := range node.Node.Acceptors {
		resp, _ := http.Get("http://" + address + "/node-info")
		respBody, _ := ioutil.ReadAll(resp.Body)

		var acceptorInfo AcceptorInfo
		json.Unmarshal(respBody, &acceptorInfo)

		node.Accepts = append(node.Accepts, &acceptorInfo)
	}
}

func learnerResponse(n *LearnerNode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respBody, _ := json.Marshal(n)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(respBody)
	}
}

func receiveAcceptLearner(n *LearnerNode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBody, _ := ioutil.ReadAll(r.Body)
		var acceptInfo AcceptorInfo
		json.Unmarshal(reqBody, &acceptInfo)

		//If proposal number by acceptor is greater than the last accepted proposal number, learner assumes it as a new possible acceptance
		if acceptInfo.N > n.AcceptedN {
			quorum := 1
			for _, v := range n.Accepts {
				//Counting all of the accepts with the same number
				if v.N == acceptInfo.N {
					quorum += 1
				}

				//Changing values in the accepts list for the acceptor that send the proposal
				if v.AcceptorId == acceptInfo.AcceptorId {
					v.AcceptorId = acceptInfo.AcceptorId
					v.N = acceptInfo.N
					v.V = acceptInfo.V
				}
			}

			//If quorum is reached, new value is accepted
			if quorum >= len(n.Node.Acceptors)/2+1 {
				n.AcceptedN = acceptInfo.N
				n.AcceptedV = acceptInfo.V

				fmt.Printf("Value %d has been accepted with proposal № %d\n", acceptInfo.V, acceptInfo.N)
				w.WriteHeader(200)
			} else {
				fmt.Printf("Proposal № %d has been rejected, quorum has not been reached\n", acceptInfo.N)
				w.WriteHeader(500)
			}
		} else {
			fmt.Printf("Proposal № %d has been rejected, last accepted proposal has greater or equal number\n", acceptInfo.N)
			w.WriteHeader(500)
		}
	}
}
