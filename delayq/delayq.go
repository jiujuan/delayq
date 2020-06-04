package delayq

import (
	"time"
)

func Start() {
	ticker := time.NewTicker(time.Second * 1)

	go func() {
		for {
			select {
			case t := <-ticker.C:
				run()
			}
		}
	}()
}

func run() {

}
