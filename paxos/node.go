package paxos

type Node struct {
	NodeId   int
	NodeType int
}

func NewNode(nodeId int, nodeType int) *Node {
	n := Node{nodeId, nodeType}
	return &n
}
