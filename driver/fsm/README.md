Finite State Machine
==============================

The state machine handles the execution of local orders. First it initilizes the elevator list that contains the state of all the elevators on the network. By sending this initilized list, all the other elevators on the network can update the state of the sender in their own list, such that the orderhandler can distribute the orders to the best choice of elevator. It then goes on to handle the door timer, door light, motor direction, floor indicator and checks for motor loss. 