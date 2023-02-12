Nodes should be started with command line parameters: port, node type, node id.

Node types: 0 - proposer, 1 - acceptor, 2 - learner.

Starting proposer with the right node id is very important since it is used for forming the unique proposal number.
Each proposal should have an id, that is greater that the id of previous proposer by 1.

For example, the command below will start a proposer with id = 1 on port 1111.
go run main.go 1111 0 1 

All of the nodes addresses must be stated in configuration file before setting up the system.

Groups of nodes should be initialized in this order for the correct functioning:
1. Acceptors
2. Learners
3. Proposers