# chanx
Break through fix sized golang chan and let writing never blocked.

## idea
Write : create new channel as current write channel if current write channel is full.
Read : read next channel if current read channel is empty.

## example:
```go
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
}
time.Sleep(1 * time.Second)
```