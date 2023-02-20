package paxos

import (
	"math/rand"
	"time"
)

type Node struct {
	NodeId    int
	NodeType  int
	Proposers []string
	Acceptors []string
	Learners  []string
	msgLimit  int
}

type AcceptorNode struct {
	Node      *Node
	PreparedN int
	PreparedV int
	AcceptedN int
	AcceptedV int
}

type AcceptorInfo struct {
	AcceptorId int
	N          int
	V          int
}

type ProposerNode struct {
	Node    *Node
	N       int
	V       int
	Accepts []*ProposeStruct
}

type LearnerNode struct {
	Node      *Node
	AcceptedN int
	AcceptedV int
	Accepts   []*AcceptorInfo
}

type ProposeStruct struct {
	N int
	V int
}

func NewNode(nodeId int, nodeType int, proposers []string, acceptors []string, learners []string, crashCount int) *Node {
	n := Node{nodeId, nodeType, proposers, acceptors, learners, crashCount}
	return &n
}

func NewAcceptorNode(node *Node) *AcceptorNode {
	preparedN := 0
	preparedV := 0
	acceptedN := 0
	acceptedV := 0

	newNode := AcceptorNode{node, preparedN, preparedV, acceptedN, acceptedV}
	return &newNode
}

func NewProposerNode(node *Node) *ProposerNode {
	rand.Seed(time.Now().UnixNano())

	//Generating proposal number and value for the first iteration
	n := node.NodeId
	v := rand.Intn(100)
	var accepts []*ProposeStruct

	newNode := ProposerNode{node, n, v, accepts}
	return &newNode
}

func NewLearnerNode(node *Node) *LearnerNode {
	var accepts []*AcceptorInfo

	newNode := LearnerNode{node, 0, 0, accepts}
	return &newNode
}
