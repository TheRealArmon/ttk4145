package timer

import "time"
import "../config"


func SetTimer(timerCh config.TimerChannels, timerCase config.TimerCase){
    switch timerCase{
    case config.Door:
      go func(){time.Sleep(3 * time.Second); timerCh.Open_door <- true}()
    }
  }

  