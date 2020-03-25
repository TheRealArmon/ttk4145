package timer

import "time"

func SetTimer(reciever chan<- bool, seconds time.Duration) {
  ticker := time.NewTicker(seconds * 1000 * time.Millisecond)
  for {
    select {
    case <- ticker.C:
      reciever <- true
      return
    }
  }
}
