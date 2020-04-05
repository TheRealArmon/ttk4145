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
					fmt.Println("This is the list being sent when new elev is connected: ", elevatorList[peerId-1])
					ch.TransmittStateCh <- map[string][config.NumElevators]config.ElevatorState{idAsString:*elevatorList}
				}
			}
			
			//If lost a peer, update the active elevator map
			if len(p.Lost) > 0{
				for _, peer := range p.Lost{
					peerId, _ := strconv.Atoi(peer)
					activeElevators[peerId-1] = false
					fmt.Println("This is the list being sent when new elev is lost: ", elevatorList[peerId-1])
					ch.TransmittStateCh <- map[string][config.NumElevators]config.ElevatorState{idAsString:*elevatorList}
				}
			}
			

		    //Update local elevator map with the state of the peers on the network
		case newState := <-ch.RecieveStateCh:
			for i, elevatorStateList := range newState{
				senderIdAsInt,_ := strconv.Atoi(i)
				elevatorList[senderIdAsInt-1] = elevatorStateList[senderIdAsInt-1]
				if checkCabQueue(elevatorStateList[idIndex]) && senderIdAsInt != id && !activeElevators[senderIdAsInt-1]{
					fmt.Println("Temp")
					tempElev := elevatorStateList[idIndex]
					go syncElev(idIndex, tempElev, elevatorList)
				}
			}

		case newOrder := <-ch.RecieveOrderCh:
			executingElevator := newOrder.ExecutingElevator
			elevatorList[executingElevator-1].Queue[newOrder.Floor][newOrder.Button] = !(newOrder.OrderStatus)
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
	if tempElev.Dir == -1 && tempElev.ElevState != config.ArrivedAtFloor{
		tempElev.Floor -= 1
	}
	for floor := 0; floor < config.NumFloors; floor++{
		elevio.SetButtonLamp(elevio.BT_Cab, floor, tempElev.Queue[floor][elevio.BT_Cab])
	}
	tempElev.ElevState = config.Idle
	tempElev.Dir = config.Stop
	time.Sleep(3 * time.Second)
	elevatorList[id] = tempElev
	fmt.Println("After setting temp elev the state is: ", elevatorList[id])
	fmt.Println("")
	return
}

func waitActive(activeElevators *[config.NumElevators]bool, id int){
	time.Sleep(2 * time.Second)
	activeElevators[id] = true
}