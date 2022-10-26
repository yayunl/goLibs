package fifo

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// fanOut := func(done <-chan interface{},
func fifo() {
	startTime := time.Now()

	done := make(chan interface{})
	defer close(done)

	// Find the prime of a given integer from the reader channel
	primeFinder := func(done <-chan interface{}, inDataStream <-chan int, fanOutWorker int) <-chan interface{} {
		valueDataStream := make(chan interface{})

		go func() {
			defer close(valueDataStream)
			defer fmt.Printf("primeFinder %d closure exited.\n", fanOutWorker)
			for v := range inDataStream {
				b := make([]bool, v)
				for i := 2; i < v; i++ {
					if b[i] == true {
						continue
					}

					select {
					case <-done:
						return
					case valueDataStream <- i:

					}

					for k := i * i; k < v; k += i {
						b[k] = true
					}
				}
			}
		}()

		return valueDataStream
	}

	generator := func(done <-chan interface{}, i ...int) <-chan int {
		valueStream := make(chan int)
		go func() {
			defer close(valueStream)
			defer fmt.Println("generator closure exited.")
			for _, v := range i {
				select {
				case <-done:
					return
				case valueStream <- v:

				}
			}
		}()
		return valueStream
	}

	// Fan in function
	fanIn := func(done <-chan interface{}, resultChannels ...<-chan interface{}) <-chan interface{} {
		var wg sync.WaitGroup
		multiplexedStream := make(chan interface{})

		multiplex := func(c <-chan interface{}, fanInWorker int) {
			defer wg.Done()
			defer fmt.Printf("Fan-In %d closure exited.\n", fanInWorker)

			for i := range c {
				select {
				case <-done:
					return
				case multiplexedStream <- i:
				}
			}
		}

		// Select from all the channels
		wg.Add(len(resultChannels))
		for i, c := range resultChannels {
			go multiplex(c, i)
		}

		// close channels and goroutines
		go func() {
			wg.Wait()
			close(multiplexedStream)
		}()
		return multiplexedStream
	}

	// data stream generator
	dataStream := generator(done, 1000, 4000, 9547, 20000, 78912, 123924)

	//***********************************************************
	// Fan out the stage of primeFinder to speed up computation
	numFinders := runtime.NumCPU()
	primeFinders := make([]<-chan interface{}, numFinders)
	for i := 0; i < numFinders; i++ {
		primeFinders[i] = primeFinder(done, dataStream, i)
	}
	//********* Fan-out code block ends *************************

	//***********  Fan-in code block starts***************************
	// Combine the results from the spawn channels in the fan-out stage
	primeStream := fanIn(done, primeFinders...)
	for v := range primeStream {
		if v == 101 {
			close(done)
		}
		fmt.Println(v)
	}
	//*********** Fan-in code block ends *****************************
	fmt.Printf("Spinning up %d prime finders.\n", numFinders)
	fmt.Printf("Time elapsed: %v\n", time.Since(startTime))

	time.Sleep(2 * time.Second)

}
