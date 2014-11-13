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
	done      chan bool
}

func New() (*Peerpipe, error) {
	peerpipe := new(Peerpipe)
	peerpipe.done = make(chan bool)
	peerpipe.listen()
	peerpipe.generateHash(false)
	log.Println("Ready for connections on port", peerpipe.Port, "at", peerpipe.addresses)
	return peerpipe, nil
}

func (self *Peerpipe) Connect(peerHash string) {
	log.Println("Connecting to", peerHash)
	addresses, port := self.decodeHash(peerHash)
	var client net.Conn
	var err error
	for _, address := range addresses {
		log.Println("via", address+":"+strconv.Itoa(port))
		client, err = net.DialTimeout("tcp", address+":"+strconv.Itoa(port), time.Second*2)
		if err == nil {
			break
		}
		fmt.Println(err)
	}
	if client == nil {
		os.Exit(1)
	}
	defer self.ListenTCP.Close()
	client.Write([]byte("hi\n"))
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
	self.Port = r.Intn(65535-1024) + 1024
	addr := net.TCPAddr{
		Port: self.Port,
		IP:   net.ParseIP("0.0.0.0"),
	}
	self.ListenTCP, err = net.ListenTCP("tcp", &addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	go self.accept()
}

func (self *Peerpipe) accept() {
	log.Println("Listening now.")
	defer self.ListenTCP.Close()
	client, err := self.ListenTCP.Accept()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	client.Write([]byte("hi\n"))
	self.done <- true
}

func (self *Peerpipe) Wait() {
	<-self.done
}

func (self *Peerpipe) decodeHash(peerHash string) ([]string, int) {
	// [todo] - figure out how to determine if the hash is "short"
	// setup our string as a rune slice
	var addresses []string
	var address string
	sliceHash := strings.Split(peerHash, "")
	for len(sliceHash) > 8 {
		sliceHash, address = RemoveOneAddress(sliceHash, 4)
		addresses = append(addresses, address)
	}
	return addresses, CharToInt(strings.Join(sliceHash, ""))
}
