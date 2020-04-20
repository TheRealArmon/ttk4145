OrderHandler
==============================

The orderhandler takes new orders and send them over the network such that the best elevator executes the order. It finds the best elevator by running the cost function. It also recieves new orders and states from its peers and updates the state in the elevatorList corrasponding to the elevator sending the messages. When an elevator is lost, the orderhandler transfers the pending orders from the lost elevator to the best elevator left on the network.