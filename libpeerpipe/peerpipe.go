package libpeerpipe

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Peerpipe struct {
	Port      int
	peerhash  string
	ListenUDP *net.UDPConn
	ListenTCP *net.TCPListener
}

func New() (*Peerpipe, error) {
	peerpipe = new(Peerpipe)
	peerpipe.listen()
	peerpipe.generateHash()
	return peerpipe
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
	self.Port = r.Intn(65535 - 1024)
	_, err = net.Listen("tcp", ":"+strconv.Itoa(self.Port))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	log.Println("Ready for connections")
}
