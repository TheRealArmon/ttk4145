package networkmod

import (
	"fmt"
	"../config"
	//"../orderhandler"
	"../elevio"
	//"reflect"
	"time"
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
			
			for _, peer := range p.New{
				peerId, _ := strconv.Atoi(peer)
				activeElevators[peerId] = true
				if elevatorList[id].Id == id{
					fmt.Println("This is the list being sent when new elev is connected: ", elevatorList[peerId])
					ch.TransmittStateCh <- map[string][config.NumElevators]config.ElevatorState{idAsString:*elevatorList}
				}
			}
			
			//If lost a peer, update the active elevator map
			if len(p.Lost) > 0{
				for _, peer := range p.Lost{
					peerId, _ := strconv.Atoi(peer)
					activeElevators[peerId] = false
					fmt.Println("This is the list being sent when new elev is lost: ", elevatorList[peerId])
					ch.TransmittStateCh <- map[string][config.NumElevators]config.ElevatorState{idAsString:*elevatorList}
				}
			}
			

		    //Update local elevator map with the state of the peers on the network
		case newState := <-ch.RecieveStateCh:
			for i, elevatorStateList := range newState{
				senderIdAsInt,_ := strconv.Atoi(i)
				elevatorList[senderIdAsInt] = elevatorStateList[senderIdAsInt]
				if checkCabQueue(elevatorStateList[id]) && senderIdAsInt != id{
					tempElev := elevatorStateList[id]
					go syncElev(id, tempElev, elevatorList)
				}
			}

		case newOrder := <-ch.RecieveOrderCh:
			executingElevator := newOrder.ExecutingElevator
			elevatorList[executingElevator].Queue[newOrder.Floor][newOrder.Button] = !(newOrder.OrderStatus)
			if (newOrder.Button != elevio.BT_Cab || executingElevator == id){
				elevio.SetButtonLamp(newOrder.Button, newOrder.Floor, !(newOrder.OrderStatus))
			}
		}

	}
}

func checkCabQueue(elevatorState config.ElevatorState) bool {
	for floor := 0; floor < config.NumFloors; floor++ {
		if elevatorState.Queue[floor][elevio.BT_Cab] {
			return true
		}
	}
	return false
}




//Making sure that the reconnecting elevator has the right state so that it can execute pre existing cab orders, as well as
//turning on lights
func syncElev(id int, tempElev config.ElevatorState, elevatorList *[config.NumElevators]config.ElevatorState){
	time.Sleep(3 * time.Second)
	if tempElev.Dir == -1 && tempElev.Queue[tempElev.Floor][elevio.BT_Cab]{
		tempElev.Floor -= 1
	}
	for floor := 0; floor < config.NumFloors; floor++{
		elevio.SetButtonLamp(elevio.BT_Cab, floor, tempElev.Queue[floor][elevio.BT_Cab])
	}
	tempElev.ElevState = config.Idle
	elevatorList[id] = tempElev
	return
}