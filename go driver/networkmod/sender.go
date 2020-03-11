package networkmod

import (
	"fmt"
	//"time"
	"../config"
)


// We define some custom struct to send over the network.
// Note that all members we want to transmit must be public. Any private members
//  will be received as zero-values.
// type HelloMsg struct {
// 	Message string
// 	Iter    int
// 	Matix   [3][4]int
// }




// The example message. We just send one of these every second.


// sender <-chan config.ElevatorOrder

func SendData(id string, ch config.NetworkChannels, newOrder chan config.ElevatorOrder) {
	go func() {
		for {
			select{
			case msg := <- newOrder:
				ch.TransmitterCh <- msg
				//time.Sleep(1 * time.Second)//må endre tid, sikkert sende meldinger mye oftere
			}
		}
	}()

	fmt.Println("Started")
	for {
		select {
		case p := <-ch.PeerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)
			if len(p.Lost) !=0 { //hvis du mister netverk med en heis
				fmt.Printf("oh no, you lost a elevator")
			}

		//case a := <-helloRx:
		//	fmt.Printf("Received: %#v\n", a)
		}
}
}
