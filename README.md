# chanx
Break through fix sized golang chan and let writing never blocked.

## idea
*Write* : create new channel as current write channel if current write channel is full.  
*Read* : read next channel if current read channel is empty.

## example:
```go
func TestNormalProc(t *testing.T) {
	cx := New(1)
	outWrite := make(chan interface{})
	go cx.ReadTo(outWrite)
	var sum int32 = 0
	for j := 0; j < 10; j++ {
		go func() {
			for {
				m := <-outWrite
				atomic.AddInt32(&sum, int32(m.(int)))
				fmt.Println("hello ", m)
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

	if sum != 4500 {
		t.Error("sum not match", sum)
	}
}
```