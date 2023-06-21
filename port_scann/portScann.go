package main

import (
	"fmt"
	"net"
	"sync"
)

func main() {
}

func concurrentPortScanner() {

	var wg sync.WaitGroup
	for i := 1; i <= 1024; i++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			address := fmt.Sprintf("scanme.nmap.org:%d", j)
			conn, err := net.Dial("tcp", address)
			if err == nil {
				fmt.Printf("Port openeed: %d\n", j)
				conn.Close() // defer?
			}

		}(i)
	}
	wg.Wait()
}
func rangePortScanner() {
	for i := 1; i <= 1024; i++ {
		address := fmt.Sprintf("scanme.nmap.org:%d", i)
		conn, err := net.Dial("tcp", address)
		if err == nil {
			fmt.Printf("Port openeed: %d\n", i)
		}
		conn.Close() // defer?
	}
}

func singlePortScanner() {
	_, err := net.Dial("tcp", "scanme.nmap.org:80")
	if err == nil {
		fmt.Printf("Port openeed: %d\n", 80)
	}
}
