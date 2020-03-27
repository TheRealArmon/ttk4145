package orderhandler

import "../elevio"
import "../config"
//import "sync"
//import "fmt"



func OrderHandler(buttonCh <-chan elevio.ButtonEvent, sendOrder chan<- config.ElevatorOrder, id int,
	elevatorList *[config.NumElevators]config.ElevatorState, activeElevators *[config.NumElevators]bool){
		for{
			select{
			case pressedButton := <- buttonCh:
				button_type := pressedButton.Button
				order_floor := pressedButton.Floor
				elevio.SetButtonLamp(button_type, order_floor, true)
				best_elevator := costCalculator(order_floor, button_type, elevatorList, activeElevators, id)
				isDone := false
				sendOrder <- config.ElevatorOrder{button_type, order_floor, best_elevator, isDone}
			}
		}
}


