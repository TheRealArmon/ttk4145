package orderhandler

import "../elevio"
import "../config"
import "strconv"
import "fmt"
//import "sync"
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
				best_elevator := costCalculator(order_floor, button_type, elevatorList, activeElevators, id)
				if best_elevator == id && checkIfOthersAreActive(activeElevators, idIndex){
					elevatorList[idIndex].Queue[order_floor][button_type] = true
					elevio.SetButtonLamp(button_type, order_floor, true)
				}
				sendOrder <- config.ElevatorOrder{button_type, order_floor, best_elevator, false}

			case newState := <- recievedStateUpdate:
				for i, elevatorStateList := range newState{
					senderIdAsInt,_ := strconv.Atoi(i)
					elevatorList[senderIdAsInt-1] = elevatorStateList[senderIdAsInt-1]
					if senderIdAsInt != id && !activeElevators[senderIdAsInt-1]{
						activeElevators[senderIdAsInt-1] = true
						stateFromSender := elevatorStateList[idIndex]
						sendersElevatorQueue := elevatorStateList[senderIdAsInt-1].Queue
						go syncElev(idIndex, stateFromSender, elevatorList)
						go turnOnHallLightsWhenReconnectingToNetwork(sendersElevatorQueue)
					}
					if elevatorStateList[senderIdAsInt-1].State == config.SystemFailure{
						activeElevators[senderIdAsInt-1] = false
					}
				}

			case newOrder := <- recievedOrder:
				executingElevator := newOrder.ExecutingElevator
				fmt.Println(elevatorList[executingElevator-1])
				if elevatorList[executingElevator-1].Queue[newOrder.Floor][newOrder.Button] != !(newOrder.OrderStatus){
					elevatorList[executingElevator-1].Queue[newOrder.Floor][newOrder.Button] = !(newOrder.OrderStatus)
					if (newOrder.Button != elevio.BT_Cab || executingElevator == id){
						elevio.SetButtonLamp(newOrder.Button, newOrder.Floor, !(newOrder.OrderStatus))
					}
					if newOrder.OrderStatus{
						SwitchOffButtonLight(newOrder.Floor)
					}
			}
			case lostElevator := <- lostConnection:
				activeElevators[lostElevator.Id-1] = false
				go transferHallOrders(lostElevator, elevatorList, activeElevators, sendOrder, sendState, id)
			}
		}
	}

