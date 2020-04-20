package orderhandler


import (
	cf "../config"
 	"../elevio"
	"strconv"
)

func OrderHandler(buttonCh <-chan elevio.ButtonEvent, orderCh cf.OrderChannels,
	recievedStateUpdate <-chan map[string][cf.NumElevators]cf.ElevatorState, recievedOrder <-chan cf.ElevatorOrder,
	lostConnection <-chan cf.ElevatorState, id int, elevatorList *[cf.NumElevators]cf.ElevatorState, activeElevators *[cf.NumElevators]bool,
	){
		idIndex := id - 1
		for{
			select{
			case pressedButton := <- buttonCh:
				button := pressedButton.Button
				floor := pressedButton.Floor
				bestElevator := costCalculator(floor, button, elevatorList, activeElevators, id)
				if bestElevator == id && !checkIfOthersAreActive(activeElevators, idIndex){
					elevatorList[idIndex].Queue[floor][button] = true
					elevio.SetButtonLamp(button, floor, true)
				}
				go func(){orderCh.SendOrder <- cf.ElevatorOrder{button, floor, bestElevator, false}}()

			case newState := <- recievedStateUpdate:
				for i, elevatorStateList := range newState{
					senderIdAsInt,_ := strconv.Atoi(i)
					elevatorList[senderIdAsInt-1] = elevatorStateList[senderIdAsInt-1]
					if senderIdAsInt != id && !activeElevators[senderIdAsInt-1]{
						activeElevators[senderIdAsInt-1] = true
						stateFromSender := elevatorStateList[idIndex]
						sendersElevatorQueue := elevatorStateList[senderIdAsInt-1].Queue
						go syncReconnectedElev(idIndex, stateFromSender, elevatorList)
						go turnOnHallLightsWhenReconnectingToNetwork(sendersElevatorQueue)
					}
					if elevatorStateList[senderIdAsInt-1].State == cf.SystemFailure{
						activeElevators[senderIdAsInt-1] = false
					}
				}

			case newOrder := <- recievedOrder:
				executingElevator := newOrder.ExecutingElevator
				elevatorList[executingElevator-1].Queue[newOrder.Floor][newOrder.Button] = !(newOrder.OrderStatus)
				if (newOrder.Button != elevio.BT_Cab || executingElevator == id){
					elevio.SetButtonLamp(newOrder.Button, newOrder.Floor, !(newOrder.OrderStatus))
				}
				if newOrder.OrderStatus{
					SwitchOffButtonLight(newOrder.Floor)
				}

			case lostElevator := <- lostConnection:
				activeElevators[lostElevator.Id-1] = false
				go transferHallOrders(lostElevator, elevatorList, activeElevators, orderCh, id)
			}
		}
	}
