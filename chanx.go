package chanx

import (
	"container/list"
	"fmt"
	"sync"
)

//Chanx Extends native chan,
type Chanx struct {
	ReadChan  chan interface{}
	WriteChan chan interface{}
	bufferLen int
	chans     *list.List
	chansLock sync.Locker
	StopChan  chan struct{}
}

//New create a chanx
func New(bufferLen int) *Chanx {
	chanx := &Chanx{
		chans:     list.New(),
		bufferLen: bufferLen,
		chansLock: &sync.Mutex{},
		StopChan:  make(chan struct{}),
	}
	chanx.ReadChan = chanx.createChan()
	chanx.WriteChan = chanx.ReadChan
	chanx.pushChan(chanx.ReadChan)
	return chanx
}

//Write a message , never block
func (c *Chanx) Write(msg interface{}) {
	select {
	case c.WriteChan <- msg:
	default: //WriteChan full
		c.addWriteBuffer()
		c.WriteChan <- msg
	}
}

//WriteFrom write message from out chan
func (c *Chanx) WriteFrom(outReadChan chan interface{}) {
	for {
		select {
		case m := <-outReadChan:
			select {
			case c.WriteChan <- m:
			default: //WriteChan full
				c.addWriteBuffer()
				c.WriteChan <- m
			}
		}
	}

}

//ReadTo Read message from Chanx to out chan
func (c *Chanx) ReadTo(outWriteChan chan interface{}) {
head:
	for {
		select {
		case m := <-c.ReadChan:
			outWriteChan <- m
		default: //ReadChan is empty
			if !c.hasMessage() { // no message, continue listen
				select {
				case m := <-c.ReadChan:
					outWriteChan <- m
					goto head
				}
			} else {
				c.ReadTo(outWriteChan)
			}
		}
	}
}

//Close close chanx
func (c *Chanx) Close() {
	// c.StopChan <- struct{}{}
	c.chansLock.Lock()
	fmt.Println("chans len", c.chans.Len(), c.chans.Front())
	for c.chans.Len() > 0 {
		e := c.chans.Front()
		close(e.Value.(chan interface{}))
		c.chans.Remove(e)
	}
	c.chansLock.Unlock()
}

func (c *Chanx) addWriteBuffer() {
	wchan := make(chan interface{}, c.bufferLen)
	c.pushChan(wchan)
	c.WriteChan = wchan
}
func (c *Chanx) createChan() chan interface{} {
	return make(chan interface{}, c.bufferLen)
}

//HasMessage return true if no message
func (c *Chanx) hasMessage() bool {
	if c.ReadChan == c.WriteChan {
		return false
	}
	c.ReadChan = c.popChan()
	return true
}
func (c *Chanx) pushChan(ch chan interface{}) {
	c.chansLock.Lock()
	c.chans.PushBack(ch)
	c.chansLock.Unlock()
}
func (c *Chanx) popChan() chan interface{} {
	c.chansLock.Lock()
	first := c.chans.Front()
	c.chans.Remove(first)
	close(first.Value.(chan interface{}))
	r := c.chans.Front().Value.(chan interface{})
	c.chansLock.Unlock()
	return r
}
func (c *Chanx) clearChans() {

}
