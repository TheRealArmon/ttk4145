package networkmod

import (
	"fmt"
	"../config"
	//"reflect"
	//"time"
	"sync"
)

// We define some custom struct to send over the network.
// Note that all members we want to transmit must be public. Any private members
//  will be received as zero-values.



func RecieveData(id string, ch config.NetworkChannels, mutex *sync.RWMutex) {

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
				mutex.Lock()
				config.ActiveElevatorMap[peer] = true
				mutex.Unlock()
			}
			

			
			//If lost a peer, update the active elevator map
			if len(p.Lost) > 0{
				for _, peer := range p.Lost{
					mutex.Lock()
					config.ActiveElevatorMap[peer] = false
					mutex.Unlock()
				}
			}
			

		    //Update local elevator map with the state of the peers on the network
		case newState := <-ch.RecieveStateCh:
			mutex.Lock()
			for id, elevatorState := range newState{
				config.ElevatorMap[id] = elevatorState
			}
			mutex.Unlock()
			

		case newOrder := <-ch.RecieveOrderCh:
			mutex.Lock()
			id := newOrder.ExecutingElevator
			var status = config.ElevatorMap[id]
			status.Queue[newOrder.Floor][newOrder.Button] = !(newOrder.OrderStatus)
			config.ElevatorMap[id] = status
			fmt.Println(config.ElevatorMap)
			mutex.Unlock()
		}

	}
}
