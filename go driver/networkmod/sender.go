package networkmod

import (
	"fmt"
	//"time"
	"../config"
)

func SendData(id string, ch config.NetworkChannels, newOrder chan config.ElevatorOrder, updateState chan config.ElevatorState) {
	go func() {
		for {
			select{
			case orderMsg := <- newOrder:
				ch.TransmitterCh <- orderMsg
				//time.Sleep(1 * time.Second)//må endre tid, sikkert sende meldinger mye oftere
			case stateMsg := <- updateState:
				ch.TransmitterCh <- updateState
			}
		}
	}()
}
