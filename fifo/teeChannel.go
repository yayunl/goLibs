package fifo

import (
	"fmt"
)

func tee(done <-chan interface{}, in <-chan interface{}) (_, _, _ <-chan interface{}) {
	//outs := make([]chan interface{}, outChanNum)
	out1 := make(chan interface{})
	out2 := make(chan interface{})
	out3 := make(chan interface{})

	go func() {
		defer close(out1)
		defer close(out2)
		defer close(out3)
		defer fmt.Println("tee closure exited.")

		for val := range in {
			var out1, out2, out3 = out1, out2, out3
			fmt.Println("read in val: ", val)
			for i := 0; i < 3; i++ {
				fmt.Printf("checking deadlock... iteration: %d\n", i)
				select {
				case <-done:
					// Write the value to the duplicated channel
				case out1 <- val:
					fmt.Println("Set out1 to nil")
					out1 = nil
				case out2 <- val:
					fmt.Println("Set out2 to nil")
					out2 = nil
				case out3 <- val:
					fmt.Println("Set out3 to nil")
					out3 = nil
				}
			}

		}
	}()

	return out1, out2, out3
}

func teeN(done <-chan interface{}, in <-chan interface{}, repChanNum int) []chan interface{} {
	//fmt.Printf("Creating a slice of %d rep channels.\n", repChanNum)
	repChannels := make([]chan interface{}, repChanNum)
	for i := 0; i < repChanNum; i++ {
		repChannels[i] = make(chan interface{})
	}

	go func() {
		for i := 0; i < repChanNum; i++ {
			defer close(repChannels[i])
		}
		//defer fmt.Println("teeN closure exited.")

		localRepChannels := make([]chan interface{}, repChanNum)
		for val := range in {
			for i := 0; i < repChanNum; i++ {
				localRepChannels[i] = repChannels[i]
			}

			for i := 0; i < repChanNum; i++ {
				select {
				case <-done:
				case localRepChannels[i] <- val:
					localRepChannels[i] = nil
				}
			}

		}

	}()

	return repChannels
}
func teeChannel() {

	repeat := func(
		done <-chan interface{},
		values ...interface{},
	) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			defer fmt.Println("repeat closure exited.")
			for {
				for _, v := range values {
					select {
					case <-done:
						return
					case valueStream <- v:
					}
				}
			}
		}()
		return valueStream
	}

	take := func(
		done <-chan interface{},
		valueStream <-chan interface{},
		num int,
	) <-chan interface{} {
		takeStream := make(chan interface{})
		go func() {
			defer close(takeStream)
			defer fmt.Println("take closure exited.")
			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeStream <- <-valueStream:
				}
			}
		}()
		return takeStream
	}

	done := make(chan interface{})
	defer close(done)
	// Create three duplicated channels from the orignal channel

	//N := 2
	//dataStreams1, dataStreams2, ds3 := tee(done, take(done, repeat(done, 1, 20), 3))
	//for val1 := range dataStreams1 {
	//	fmt.Printf("out1: %v, out2: %v, out3: %v\n", val1, <-dataStreams2, <-ds3)
	//}

	ds := teeN(done, take(done, repeat(done, 1, 20, 100), 3), 4)

	for val1 := range ds[0] {
		fmt.Printf("out1: %v, out2: %v, out3: %d, out4: %d\n", val1, <-ds[1], <-ds[2], <-ds[3])
	}

	//dataStream1 := take(done, repeat(done, 1, 20, 3), 5)
	//for val := range dataStream1 {
	//	fmt.Printf("val: %v\n", val)
	//}

	//time.Sleep(2 * time.Second)
}
