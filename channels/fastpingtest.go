package channels

import (
	"fmt"
	"github.com/tatsushid/go-fastping"
	"net"
	"os"
	"sync"
	"time"
)

type response struct {
	ipAddr *net.IPAddr
	active bool
	rtt    time.Duration
}

func ipAddrGenerator(done <-chan interface{}, ips ...string) <-chan string {
	// Generate a stream of IP addresses
	ipAddrStream := make(chan string)
	go func() {
		defer close(ipAddrStream)
		for _, ipAddr := range ips {
			select {
			case <-done:
				return
			case ipAddrStream <- ipAddr:
			}
		}
	}()
	return ipAddrStream
}

func fastPing(ipAddr string) *response {
	// create a new pinger
	p := fastping.NewPinger()

	// Resolve the input ip address
	ra, err := net.ResolveIPAddr("ip4:icmp", ipAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	resp := &response{ipAddr: ra, rtt: 0, active: false}

	// Create a map to store the results of a given ip addr
	//results := make(map[string]*response)
	//results[ra.String()] = nil
	p.AddIPAddr(ra)

	// Create two channels for received responses and idle respectively
	onRecv, onIdle := make(chan *response), make(chan bool)
	defer close(onRecv)
	defer close(onIdle)

	p.OnRecv = func(addr *net.IPAddr, t time.Duration) {
		//fmt.Println("received")
		*resp = response{ipAddr: addr, rtt: t, active: true}
	}
	//p.OnIdle = func() {
	//	fmt.Println("no response.")
	//	*resp = response{ipAddr: ra, rtt: 0, active: false}
	//}

	// Set a timeout of 1 sec for a ping test
	//p.MaxRTT = time.Second
	err = p.Run()
	p.Stop()
	return resp
}

func ping(done <-chan interface{}, inIPAddrStream <-chan string) <-chan *response {
	resultStream := make(chan *response)
	go func() {
		defer close(resultStream)

		//defer fmt.Println("Exit pingtest closure.")

		for ipAddr := range inIPAddrStream {
			select {
			case <-done:
				return
			case resultStream <- fastPing(ipAddr):
			}
		}

	}()

	return resultStream
}

func aggregateResults(done <-chan interface{}, resultChannels ...<-chan *response) <-chan *response {
	var wg sync.WaitGroup
	multiplexedStream := make(chan *response)

	multiplex := func(c <-chan *response, fanInWorker int) {
		defer wg.Done()
		//defer fmt.Printf("Fan-In %d closure exited.\n", fanInWorker)

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

func fastpingtest() {
	start := time.Now()
	ips := []string{}
	for i := 1; i < 255; i++ {
		ip := fmt.Sprintf("100.69.46.%d", i)
		ips = append(ips, ip)
	}

	done := make(chan interface{})
	defer close(done)

	ipStream := ipAddrGenerator(done, ips...)
	// Create multiple workers
	//workerNum := runtime.NumCPU()
	workerNum := 80
	workers := make([]<-chan *response, workerNum)
	for i := 0; i < workerNum; i++ {
		workers[i] = ping(done, ipStream)
	}

	// Aggregate the results
	resultStream := aggregateResults(done, workers...)

	// Collect results
	activeIPs := []*response{}
	//percent := 0.0
	respIPCnt := 0.0
	for res := range resultStream {
		//fmt.Printf("IP=> %s, rtt=> %v, active=> %v\n", res.ipAddr, res.rtt, res.active)
		respIPCnt += 1.0
		//percent = respIPCnt / 255.0 * 100
		//fmt.Printf("Percent: %.2f%\n", percent)
		if res.active {
			activeIPs = append(activeIPs, res)
		}
	}

	fmt.Println("Time elapsed: ", time.Since(start))
	fmt.Println("Active IPs:")
	for _, ip := range activeIPs {
		fmt.Println(ip.ipAddr)
	}
}
