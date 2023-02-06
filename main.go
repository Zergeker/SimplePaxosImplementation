package main

import (
	"fmt"
	"os"
	"strconv"

	"example.com/SimplePaxos/paxos"
)

func main() {
	port := os.Args[1]
	//0 - proposer, 1 - acceptor, 2 - learner
	nodeType, _ := strconv.Atoi(os.Args[2])
	nodeId, _ := strconv.Atoi(os.Args[3])

	n := paxos.NewNode(nodeId, nodeType)

	switch nodeType {
	case 0:
		paxos.StartProposerController(n, port)
	case 1:
		paxos.StartAcceptorController(n, port)
	case 2:
		paxos.StartLearnerController(n, port)
	default:
		fmt.Println("Unknown node type entered")
		os.Exit(0)
	}
}
