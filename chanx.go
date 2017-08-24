package chanx

import (
	"container/list"
	"sync"
)

type Chanx struct {
	ReadChan  chan interface{}
	WriteChan chan interface{}
	bufferLen int
	chans     *list.List
	chansLock sync.Locker
}

func NewChanx(bufferLen int) *Chanx {
	chanx := &Chanx{
		chans:     list.New(),
		bufferLen: bufferLen,
		chansLock: &sync.Mutex{},
	}
	chanx.ReadChan = chanx.createChan()
	chanx.WriteChan = chanx.ReadChan
	chanx.chans.PushBack(chanx.ReadChan)
	return chanx
}

func (c *Chanx) createChan() chan interface{} {
	return make(chan interface{}, c.bufferLen)
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
	r := c.chans.Front().Value.(chan interface{})
	c.chansLock.Unlock()
	return r
}
func (c *Chanx) AddWriteBuffer() {
	wchan := make(chan interface{}, c.bufferLen)
	c.pushChan(wchan)
	c.WriteChan = wchan
}
func (c *Chanx) Write(m interface{}) {
	select {
	case c.WriteChan <- m:
	default: //WriteChan full
		c.AddWriteBuffer()
		c.WriteChan <- m
	}
}
func (c *Chanx) WriteFrom(outReadChan chan interface{}) {
	for {
		select {
		case m := <-outReadChan:
			select {
			case c.WriteChan <- m:
			default: //WriteChan full
				c.AddWriteBuffer()
				c.WriteChan <- m
			}
		}
	}

}

//return true if no message
func (c *Chanx) HasMessage() bool {

	if c.ReadChan == c.WriteChan {
		return false
	}
	c.ReadChan = c.popChan()
	return true
}
func (c *Chanx) ReadTo(outWriteChan chan interface{}) {
head:
	for {
		select {
		case m := <-c.ReadChan:
			outWriteChan <- m
		default: //ReadChan is empty
			if !c.HasMessage() { // no message, continue listen
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
