package networkmod

import (
	"fmt"

	"strconv"
	"time"

	"../config"
)

func RecieveData(id int, ch config.NetworkChannels, elevatorList *[config.NumElevators]config.ElevatorState, activeElevators *[config.NumElevators]bool,
	lostConnection chan<- config.ElevatorState) {
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
				go waitActive(activeElevators, peerId-1, true)
				if elevatorList[idIndex].Id == id{
					ch.TransmittStateCh <- map[string][config.NumElevators]config.ElevatorState{idAsString:*elevatorList}
				}
			}

			//If lost a peer, update the active elevator map
			if len(p.Lost) > 0{
				for _, peer := range p.Lost{
					peerId, _ := strconv.Atoi(peer)
					if peerId != id{
						activeElevators[peerId-1] = false
						//go waitActive(activeElevators, peerId-1, false)
						lostConnection <- elevatorList[peerId-1]
					}
				}
			}
		}
	}
}


func waitActive(activeElevators *[config.NumElevators]bool, id int, state bool){
	time.Sleep(1 * time.Second)
	activeElevators[id] = state
}
