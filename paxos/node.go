package paxos

type Node struct {
	NodeId    int
	NodeType  int
	Proposers []string
	Acceptors []string
	Learners  []string
}

func NewNode(nodeId int, nodeType int, proposers []string, acceptors []string, learners []string) *Node {
	n := Node{nodeId, nodeType, proposers, acceptors, learners}
	return &n
}
