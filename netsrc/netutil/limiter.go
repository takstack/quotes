package netutil

import (
	"time"

	LOG "github.com/takstack/logger"
)

//Limiter limits the rate
func Limiter(lim <-chan interface{}, q1 <-chan interface{}) {
	c := time.Tick(100 * time.Millisecond)
loop:
	for next := range c {
		//fmt.Printf("%v %T\n", next, next)
		select {
		default:
			_ = next
		case <-lim:
			//receives from channel unblocks lim chan
		case <-q1:
			LOG.GL.Println("quit loop broken")
			break loop
		}
	}
	//return

}
