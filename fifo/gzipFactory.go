package fifo

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

func gizpFactory() {
	startTime := time.Now()
	var wg sync.WaitGroup
	//for _, fileName := range os.Args[1:] {
	//	err := compress(fileName)
	//	if err != nil {
	//		fmt.Printf("Failed to compress %s", fileName)
	//	}
	//}

	// process with goroutines
	for i := 0; i < 20; i++ {
		// spawn up goroutines
		wg.Add(1)
		go compressWorker(&wg)
	}
	wg.Wait()
	fmt.Printf("Process with goroutines took: %v \n", time.Since(startTime))

	// process sequentially
	start := time.Now()
	for i := 0; i < 20; i++ {
		compressSeqWorker()
	}
	fmt.Printf("Process sequentially took: %v \n", time.Since(start))

}

func compressWorker(wg *sync.WaitGroup) {
	defer wg.Done()
	// do some processing work
	time.Sleep(1 * time.Nanosecond)
}

func compressSeqWorker() {
	// do some processing work
	time.Sleep(1 * time.Nanosecond)
}

func compress(fn string) error {
	in, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer in.Close()

	// Create a destination file to write
	out, err := os.Create(fn + ".gz")
	if err != nil {
		return err
	}
	defer out.Close()

	// Compress the data and then writes it to the destination file
	gzout := gzip.NewWriter(out)
	_, err = io.Copy(gzout, in)
	gzout.Close()

	return err
}
