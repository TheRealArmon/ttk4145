package networkmod

import (
	"fmt"
	"../config"
	"../orderhandler"
	//"reflect"
	//"time"
	//"sync"
	"strconv"
)

func RecieveData(id int, ch config.NetworkChannels, elevatorList *[config.NumElevators]config.ElevatorState,
	activeElevators *[config.NumElevators]bool) {
	idAsString := strconv.Itoa(id)

	fmt.Println("Started")
	for {
		select {
		case p := <-ch.PeerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)
			
			//If recievig from a new peer, upadate avtive elevator map
			//ch.TransmittStateCh <- *elevatorList
			
			for _, peer := range p.New{
				peerId, _ := strconv.Atoi(peer)
				activeElevators[peerId] = true
				//go func(){ch.TransmittStateCh <- map[string][config.NumElevators]config.ElevatorState{idAsString:*elevatorList}}()
			}
			
			//If lost a peer, update the active elevator map
			if len(p.Lost) > 0{
				for _, peer := range p.Lost{
					peerId, _ := strconv.Atoi(peer)
					activeElevators[peerId] = false
					//go func(){ch.TransmittStateCh <- map[string][config.NumElevators]config.ElevatorState{idAsString:*elevatorList}}()
				}
			}
			
			ch.TransmittStateCh <- map[string][config.NumElevators]config.ElevatorState{idAsString:*elevatorList}

		    //Update local elevator map with the state of the peers on the network
		case newState := <-ch.RecieveStateCh:
			for i, elevatorStateList := range newState{
				senderIdAsInt,_ := strconv.Atoi(i)
				elevatorList[senderIdAsInt] = elevatorStateList[senderIdAsInt]
			}
			

		case newOrder := <-ch.RecieveOrderCh:
			id := newOrder.ExecutingElevator
			elevatorList[id].Queue[newOrder.Floor][newOrder.Button] = !(newOrder.OrderStatus)
			if (newOrder.OrderStatus){
				orderhandler.SwitchOffButtonLight(newOrder.Floor)
			}

		}

	}
}
