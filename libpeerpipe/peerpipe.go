package libpeerpipe

import (
	"bufio"
	"fmt"
	"io"
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
	defer client.Close()
	self.readwrite(os.Stdin, client)
	client.Close()
	log.Println("Done with Connect")
	self.done <- true
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
	//var err error
	client, err := self.ListenTCP.Accept()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	self.readwrite(client, os.Stdout)
	log.Println("Done with accept")
	self.done <- true
}

func (self *Peerpipe) Wait() {
	fmt.Println("Waiting for session to end")
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

func (self *Peerpipe) readwrite(reader io.Reader, writer io.Writer) {
	nBytes, nChunks := int64(0), int64(0)
	buf := make([]byte, 0, 1024)
	r := bufio.NewReader(reader)
	for {
		n, err := r.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		writer.Write(buf)
		nChunks++
		nBytes += int64(len(buf))
		// process buf
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
	}
}
