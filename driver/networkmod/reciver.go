package networkmod

import (
	"fmt"

	"../config"
<<<<<<< HEAD
	"../elevio"
	"../orderhandler"

	//"reflect"
	//"time"
	//"sync"
=======
	"time"
>>>>>>> origin/development
	"strconv"
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
<<<<<<< HEAD
			//ch.TransmittStateCh <- *elevatorList

			for _, peer := range p.New {
				peerId, _ := strconv.Atoi(peer)
				activeElevators[peerId] = true
=======
			
			for _, peer := range p.New{
				peerId, _ := strconv.Atoi(peer)
				go waitActive(activeElevators, peerId-1)
				if elevatorList[idIndex].Id == id{
					ch.TransmittStateCh <- map[string][config.NumElevators]config.ElevatorState{idAsString:*elevatorList}
				}
>>>>>>> origin/development
			}

			//If lost a peer, update the active elevator map
			if len(p.Lost) > 0 {
				for _, peer := range p.Lost {
					peerId, _ := strconv.Atoi(peer)
<<<<<<< HEAD
					activeElevators[peerId] = false
				}
			}

			ch.TransmittStateCh <- map[string][config.NumElevators]config.ElevatorState{idAsString: *elevatorList}

			//Update local elevator map with the state of the peers on the network
		case newState := <-ch.RecieveStateCh:
			for i, elevatorStateList := range newState {
				senderIdAsInt, _ := strconv.Atoi(i)
				elevatorList[senderIdAsInt] = elevatorStateList[senderIdAsInt]
				if checkCabQueue(elevatorStateList[id]) {
					elevatorList[id] = elevatorStateList[id]
				}
			}
			fmt.Println("Id: ", id)
			fmt.Println("State: ", elevatorList[id])
			fmt.Println("")

		case newOrder := <-ch.RecieveOrderCh:
			id := newOrder.ExecutingElevator
			elevatorList[id].Queue[newOrder.Floor][newOrder.Button] = !(newOrder.OrderStatus)
			if newOrder.OrderStatus {
				orderhandler.SwitchOffButtonLight(newOrder.Floor)
=======
					activeElevators[peerId-1] = false
					lostConnection <- elevatorList[peerId-1]

				}
>>>>>>> origin/development
			}
		}
	}
}

<<<<<<< HEAD
func checkCabQueue(elevatorState config.ElevatorState) bool {
	for floor := 0; floor < config.NumFloors; floor++ {
		if elevatorState.Queue[floor][elevio.BT_Cab] {
			return true
		}
	}
	return false
=======

func waitActive(activeElevators *[config.NumElevators]bool, id int){
	time.Sleep(2 * time.Second)
	activeElevators[id] = true
>>>>>>> origin/development
}
