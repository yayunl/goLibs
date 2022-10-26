package fifo

import (
	"fmt"
	"os"
	//"sync"
	"time"
)

func readStdin2(done <-chan interface{}, byteStream chan<- []byte) {

	// This function writes values to the writer channel `byteStream` and present in Stdout
	defer close(byteStream)
	defer fmt.Println("Close byteStream channel.")
	defer fmt.Println("goroutine `readStdin` closure closed.")

	reader := make([]byte, 1024)
	//defaultBytes := []byte{'a', 'b', 'c', 'd', 'e'}
	for {
		select {
		case <-done:
			//fmt.Println("goroutine received timeout")
			return
		default:
			fmt.Println("Blocking here until user typed.")
			l, err := os.Stdin.Read(reader)
			if err != nil { // Maybe log non io.EOF errors, if you want
				return
			}
			if l > 0 {
				byteStream <- reader
			}
			//_, ok := <-done
			//if !ok {
			//
			//	close(byteStream)
			//}
		}

	}

	//for _, s := range []byte{'a', 'b', 'c', 'd', 'e'} {
	//	select {
	//	case <-done:
	//		fmt.Println("goroutine received timeout")
	//		return
	//	case byteStream <- s:
	//		time.Sleep(1 * time.Second)
	//
	//		//default:
	//		//	data := make([]byte, 1024)
	//		//	l, _ := os.Stdin.Read(data)
	//		//	if l > 0 {
	//		//		byteStream <- data
	//		//	}
	//	}
	//}
}

func readstd() {
	//var wg sync.WaitGroup
	//Create a ticker
	totalTime := 5

	interval := time.Duration(1) * time.Second
	tk := time.NewTicker(interval)

	// timeout channel
	//done := time.After(time.Duration(totalTime) * time.Second)
	done := make(chan interface{})

	byteStream := make(chan []byte)
	// Create a goroutine for reading the user's inputs
	//wg.Add(1)
	go readStdin2(done, byteStream)
	// a channel of []type for facilitating the communication between the user's inputs and output console
loop:
	for {
		select {
		case <-done:
			fmt.Println("Main Timed out.")
			//os.Exit(0)
			break loop

		case buf, ok := <-byteStream:
			if !ok { //closed channel

				fmt.Print("Timed out")
				break loop
			}
			os.Stdout.Write(buf)

		case <-tk.C:
			totalTime -= 1
			if totalTime == 2 {
				fmt.Printf("timeout... %ds\n", totalTime)
				close(done)
			}
		}
	}

	//wg.Wait()
	fmt.Println("Exited.")
	// attempted to read from the byteStream ch to check if it is closed
	_, ok := (<-byteStream)
	if !ok {
		fmt.Println("byteStream is closed.")
	} else {
		fmt.Println("byteStream is NOT closed.")
	}

	time.Sleep(1 * time.Second)

}
