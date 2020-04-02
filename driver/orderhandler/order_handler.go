package orderhandler

import (
	"strconv"

	"../config"
	"../elevio"
)

//import "sync"
//import "fmt"

func OrderHandler(buttonCh <-chan elevio.ButtonEvent, sendOrder chan<- config.ElevatorOrder, id int,
	elevatorList *[config.NumElevators]config.ElevatorState, activeElevators *[config.NumElevators]bool,
	newState chan<- map[string][config.NumElevators]config.ElevatorState) {
	idAsString := strconv.Itoa(id)
	for {
		select {
		case pressedButton := <-buttonCh:
			button_type := pressedButton.Button
			order_floor := pressedButton.Floor
			elevio.SetButtonLamp(button_type, order_floor, true)
			best_elevator := costCalculator(order_floor, button_type, elevatorList, activeElevators, id)
			isDone := false
			go func() { sendOrder <- config.ElevatorOrder{button_type, order_floor, best_elevator, isDone} }()
			go func() { newState <- map[string][config.NumElevators]config.ElevatorState{idAsString: *elevatorList} }()

		}
	}
}
