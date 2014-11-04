package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

var shortHash = flag.Bool("s", true, "Short hash, only good for local networks.")

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) > 1 {
		fmt.Println("usage")
		os.Exit(1)
	}
	if len(args) == 1 {
		log.Println("Connecting to someone")

	} else {
		addresses, err := net.InterfaceAddrs()
		for _, addr := range addresses {
			log.Println("Found", addr.String())
		}
		_, err = net.Listen("tcp", ":6000")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		log.Println("Ready for connections")
	}
}
