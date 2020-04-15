package orderhandler

import "../elevio"
import "../config"
import "strconv"
import "fmt"



func OrderHandler(buttonCh <-chan elevio.ButtonEvent, sendOrder chan<- config.ElevatorOrder, recievedStateUpdate <-chan map[string][config.NumElevators]config.ElevatorState,
	recievedOrder <-chan config.ElevatorOrder, lostConnection <-chan config.ElevatorState, id int, elevatorList *[config.NumElevators]config.ElevatorState,
	activeElevators *[config.NumElevators]bool,
	){
		idIndex := id - 1
		for{
			select{
			case pressedButton := <- buttonCh:
				button_type := pressedButton.Button
				order_floor := pressedButton.Floor
				elevio.SetButtonLamp(button_type, order_floor, true)
				best_elevator := costCalculator(order_floor, button_type, elevatorList, activeElevators, id)
				isDone := false
				sendOrder <- config.ElevatorOrder{button_type, order_floor, best_elevator, isDone}

			case newState := <- recievedStateUpdate:
				for i, elevatorStateList := range newState{
					senderIdAsInt,_ := strconv.Atoi(i)
					elevatorList[senderIdAsInt-1] = elevatorStateList[senderIdAsInt-1]
					if checkCabQueue(elevatorStateList[idIndex]) && senderIdAsInt != id && !activeElevators[senderIdAsInt-1]{
						fmt.Println("Temp")
						tempElev := elevatorStateList[idIndex]
						go syncElev(idIndex, tempElev, elevatorList)
					}
				}

			case newOrder := <- recievedOrder:
				executingElevator := newOrder.ExecutingElevator
				elevatorList[executingElevator-1].Queue[newOrder.Floor][newOrder.Button] = !(newOrder.OrderStatus)
				if (newOrder.Button != elevio.BT_Cab || executingElevator == id){
					elevio.SetButtonLamp(newOrder.Button, newOrder.Floor, !(newOrder.OrderStatus))
				}

			// case lostElevator := <- lostConnection:
			// 	lostElevatorIndex := lostElevator.Id-1
			// 	for floor := 0; floor < config.NumFloors; floor++{
			// 		for button := elevio.BT_HallUp; button < elevio.BT_Cab; button++{
			// 			if lostElevator.Queue[floor][button]{
			// 				elevatorList[lostElevatorIndex].Queue[floor][button] = false
			// 				//newExecutingElevator := costCalculator(floor, button, elevatorList, activeElevators, id)
			//
			// 			}
			// 		}
			// 	}
			}
		}
}
