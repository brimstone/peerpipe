package libpeerpipe

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type Peerpipe struct {
	Port      int
	peerHash  string
	ListenUDP *net.UDPConn
	ListenTCP *net.TCPListener
	addresses string
}

func New() (*Peerpipe, error) {
	peerpipe := new(Peerpipe)
	peerpipe.listen()
	peerpipe.generateHash(false)
	log.Println("Ready for connections on port", peerpipe.Port, "at", peerpipe.addresses)
	return peerpipe, nil
}

func (self *Peerpipe) Connect(peerhash string) {
	log.Println("Connecting to", peerhash)
}

func (self *Peerpipe) generateHash(shortHash bool) string {
	peerHash := ""

	externalIP, err := GetExternalIP()
	if err != nil {
		fmt.Println(err)
	} else {
		peerHash = MakeReadable(externalIP.To4())
	}

	addresses, err := net.InterfaceAddrs()
	for _, addr := range addresses {
		// Parse out our IP
		address := strings.Split(addr.String(), "/")
		ip := net.ParseIP(address[0])
		// Skip local and loopbacks
		if !ip.IsGlobalUnicast() {
			continue
		}
		self.addresses += address[0] + " "
		// Convert our IP to a hash
		ip = ip.To4()
		peerHash += MakeReadable(ip)
	}
	self.peerHash = peerHash + IntToChar(self.Port)
	return peerHash
}

func (self *Peerpipe) GetHash() string {
	return self.peerHash
}

func (self *Peerpipe) listen() {
	var err error
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	self.Port = r.Intn(65535 - 1024) + 1024
	_, err = net.Listen("tcp", ":"+strconv.Itoa(self.Port))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
