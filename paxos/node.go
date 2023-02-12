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
}

type AcceptorNode struct {
	Node *Node
	N    int
	V    int
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

func NewNode(nodeId int, nodeType int, proposers []string, acceptors []string, learners []string) *Node {
	n := Node{nodeId, nodeType, proposers, acceptors, learners}
	return &n
}

func NewAcceptorNode(node *Node) *AcceptorNode {
	n := 0
	v := 0

	newNode := AcceptorNode{node, n, v}
	return &newNode
}

func NewProposerNode(node *Node) *ProposerNode {
	rand.Seed(time.Now().UnixNano())

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
