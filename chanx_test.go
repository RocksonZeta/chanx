package chanx

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func TestNormalProc(t *testing.T) {
	cx := NewChanx(1)
	outWrite := make(chan interface{})
	go cx.ReadTo(outWrite)
	var count int32 = 0
	for j := 0; j < 10; j++ {
		go func() {
			for {
				m := <-outWrite
				fmt.Println("hello ", m)
				atomic.AddInt32(&count, 1)
			}
		}()
	}
	for j := 0; j < 100; j++ {
		go func() {
			for i := 0; i < 10; i++ {
				cx.Write(i)

			}
		}()
	}
	time.Sleep(1 * time.Second)

	if count != 1000 {
		t.Error("count not match", count)
	}
	time.Sleep(2 * time.Second)
}
