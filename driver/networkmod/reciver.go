package networkmod

import (
	"fmt"
	"strconv"
	"time"
	cf "../config"
)

func UpdatePeers(id int, ch cf.NetworkChannels, lostConnection chan<- cf.ElevatorState, elevatorList *[cf.NumElevators]cf.ElevatorState, 
	activeElevators *[cf.NumElevators]bool) {
	
	idAsString := strconv.Itoa(id)
	idIndex := id - 1
	for {
		select {
		case p := <-ch.PeerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)
			
			//When recieving new peer send state if initilized with correct id. Also update activeElevators list after 1 second
			for _, peer := range p.New{
				peerId, _ := strconv.Atoi(peer)
				go func(){time.Sleep(1 * time.Second); activeElevators[peerId-1] = true}()
				if elevatorList[idIndex].Id == id{
					ch.TransmittStateCh <- map[string][cf.NumElevators]cf.ElevatorState{idAsString:*elevatorList}
				}
			}

			//If lost a peer, update the active elevator list and start order redistribution 
			if len(p.Lost) > 0{
				for _, peer := range p.Lost{
					peerId, _ := strconv.Atoi(peer)
					if peerId != id{
						activeElevators[peerId-1] = false
						lostConnection <- elevatorList[peerId-1]
					}
				}
			}
		}
	}
}


