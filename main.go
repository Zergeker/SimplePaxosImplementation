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
	//Reading config
	viper.SetConfigFile("config.env")
	viper.ReadInConfig()

	port := os.Args[1]
	//0 - proposer, 1 - acceptor, 2 - learner
	nodeType, _ := strconv.Atoi(os.Args[2])
	nodeId, _ := strconv.Atoi(os.Args[3])
	proposers := strings.Split(viper.GetString("PROPOSERS"), ",\n")
	acceptors := strings.Split(viper.GetString("ACCEPTORS"), ",\n")
	learners := strings.Split(viper.GetString("LEARNERS"), ",\n")
	minDelay := viper.GetInt("MIN_DELAY")
	maxDelay := viper.GetInt("MAX_DELAY")

	n := paxos.NewNode(nodeId, nodeType, proposers, acceptors, learners)

	//Choosing a role for the node
	switch nodeType {
	case 0:
		paxos.StartProposerController(n, port, minDelay, maxDelay)
	case 1:
		paxos.StartAcceptorController(n, port, minDelay, maxDelay)
	case 2:
		paxos.StartLearnerController(n, port)
	default:
		fmt.Println("Unknown node type entered")
		os.Exit(0)
	}
}
