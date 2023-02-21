Nodes should be started with command line parameters: port, node type, node id, crash limit.

Node types: 0 - proposer, 1 - acceptor, 2 - learner.

Starting proposer with the right node id is very important since it is used for forming the unique proposal number.
Each proposal should have an id, that is greater that the id of previous proposer by 1.

Crash limit parameter stands for the number of messages the node can receive until going into the crash state. It should have a value of an unsigned integer. For now, only acceptors crashing can be simulated. Unfortunately, making this parameter optional is TBD, it's compulsory for now.

For example, the command below will start an acceptor with id = 1 on port 1111, which would crash after 15 received messages.
go run main.go 1111 1 1 15

All of the nodes addresses must be stated in configuration file before setting up the system.

Groups of nodes should be initialized in this order for the correct functioning:
1. Acceptors
2. Learners
3. Proposers

You can find the scripts to start the system (or some specific group of nodes) in a "scripts" folder.