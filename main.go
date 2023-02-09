package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"example.com/SimplePaxos/paxos"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("config.env")
	viper.ReadInConfig()

	port := os.Args[1]
	//0 - proposer, 1 - acceptor, 2 - learner
	nodeType, _ := strconv.Atoi(os.Args[2])
	nodeId, _ := strconv.Atoi(os.Args[3])
	proposers := strings.Split(viper.GetString("PROPOSERS"), ",\n")
	acceptors := strings.Split(viper.GetString("ACCEPTORS"), ",\n")
	learners := strings.Split(viper.GetString("LEARNERS"), ",\n")

	n := paxos.NewNode(nodeId, nodeType, proposers, acceptors, learners)

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
