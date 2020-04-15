package orderhandler

import "../elevio"
import "../config"
import "strconv"
//import "fmt"



func OrderHandler(buttonCh <-chan elevio.ButtonEvent, sendOrder chan<- config.ElevatorOrder, sendState chan<- map[string][config.NumElevators]config.ElevatorState, 
	recievedStateUpdate <-chan map[string][config.NumElevators]config.ElevatorState, recievedOrder <-chan config.ElevatorOrder, 
	lostConnection <-chan config.ElevatorState, id int, elevatorList *[config.NumElevators]config.ElevatorState, activeElevators *[config.NumElevators]bool,
	){
		idIndex := id - 1
		for{
			select{
			case pressedButton := <- buttonCh:
				button_type := pressedButton.Button
				order_floor := pressedButton.Floor
				elevio.SetButtonLamp(button_type, order_floor, true)
				best_elevator := costCalculator(order_floor, button_type, elevatorList, activeElevators, id)
				sendOrder <- config.ElevatorOrder{button_type, order_floor, best_elevator, false}
			
			case newState := <- recievedStateUpdate:
				for i, elevatorStateList := range newState{
					senderIdAsInt,_ := strconv.Atoi(i)
					elevatorList[senderIdAsInt-1] = elevatorStateList[senderIdAsInt-1]
					if checkCabQueue(elevatorStateList[idIndex]) && senderIdAsInt != id && !activeElevators[senderIdAsInt-1]{
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

			case lostElevator := <- lostConnection:
				go transferHallOrders(lostElevator, elevatorList, activeElevators, sendOrder, sendState, id)
			}
		}
}
