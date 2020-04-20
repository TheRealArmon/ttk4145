package networkmod


import (
	//"fmt"
	"time"
	"../config"
)

//Sends data 10 times with a frequency of 20 per second
func SendData(ch config.NetworkChannels, newOrder <-chan config.ElevatorOrder, newState <-chan map[string][config.NumElevators]config.ElevatorState) {
	interval := 15 * time.Millisecond
	for {
		select{
		case orderMsg := <- newOrder:
			for i := 0; i < 10; i++{
				ch.TransmittOrderCh <- orderMsg
				time.Sleep(interval)
			}

		case stateMsg := <- newState:
			for i := 0; i < 10; i++{
				ch.TransmittStateCh <- stateMsg
				time.Sleep(interval)
			}
		}	
	}	
}