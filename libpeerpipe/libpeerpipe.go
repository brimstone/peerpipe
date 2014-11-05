package libpeerpipe

import (
	"fmt"
	"log"
	"net"
	"os"
)

func Connect(peerhash string) {
	log.Println("Connecting to", peerhash)
}

func Listen(shortHash bool) {
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
