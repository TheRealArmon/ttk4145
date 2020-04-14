package networkmod

import (
	"fmt"
	"../config"
	"time"
	"strconv"
)

func RecieveData(id int, ch config.NetworkChannels, elevatorList *[config.NumElevators]config.ElevatorState, activeElevators *[config.NumElevators]bool) {
	idAsString := strconv.Itoa(id)
	idIndex := id - 1

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
				peerId, _ := strconv.Atoi(peer)
				go waitActive(activeElevators, peerId-1)
				if elevatorList[idIndex].Id == id{
					//fmt.Println("This is the list being sent when new elev is connected: ", elevatorList[peerId-1])
					ch.TransmittStateCh <- map[string][config.NumElevators]config.ElevatorState{idAsString:*elevatorList}
				}
			}
			
			//If lost a peer, update the active elevator map
			if len(p.Lost) > 0{
				for _, peer := range p.Lost{
					peerId, _ := strconv.Atoi(peer)
					activeElevators[peerId-1] = false
					//fmt.Println("This is the list being sent when new elev is lost: ", elevatorList[peerId-1])
					ch.TransmittStateCh <- map[string][config.NumElevators]config.ElevatorState{idAsString:*elevatorList}

				}
			}
		}
	}
}


func waitActive(activeElevators *[config.NumElevators]bool, id int){
	time.Sleep(2 * time.Second)
	activeElevators[id] = true
}
