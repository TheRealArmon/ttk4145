package networkmod

import (
	"fmt"
	"../config"
	//"reflect"
	//"time"
)

// We define some custom struct to send over the network.
// Note that all members we want to transmit must be public. Any private members
//  will be received as zero-values.



func RecieveData(id string, ch config.NetworkChannels) {

	fmt.Println("Started")
	for {
		select {
		case p := <-ch.PeerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)
			
			//If recievig from a new peer, upadate avtive elevator map
			 
			for _, peer := range p.New{
				config.ActiveElevatorMap[peer] = true
			}
			
			//If lost a peer, update the active elevator map
			if len(p.Lost) > 0{
				for _, peer := range p.Lost{
					config.ActiveElevatorMap[peer] = false
				}
			}
			

		    //Update local elevator map with the state of the peers on the network
		case newState := <-ch.RecieveStateCh:
			for id, elevatorState := range newState{
				config.ElevatorMap[id] = elevatorState
			}
			

		case newOrder := <-ch.RecieveOrderCh:
			id := newOrder.ExecutingElevator
			var status = config.ElevatorMap[id]
			status.Queue[newOrder.Floor][newOrder.Button] = !(newOrder.OrderStatus)
			config.ElevatorMap[id] = status
		}

	}
}
