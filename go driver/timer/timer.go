package timer

import "time"

func SetDoorTimer(reciever chan<- bool) {
  ticker := time.NewTicker(3000 * time.Millisecond)
  for {
    select {
    case <- ticker.C:
      reciever <- true
      return
    }
  }
}
