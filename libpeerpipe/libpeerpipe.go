package libpeerpipe

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

var charMapping = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	"A", "B", "C", "D", "E", "F", "G", "H", "J", "K", "M", "N", "P", "Q", "R", "S",
	"T", "U", "V", "X", "Y", "Z"}

// helper functions

func IntToChar(input int) string {
	currentByte := ""
	for input > 0 {
		newbase := input % len(charMapping)
		currentByte = charMapping[newbase] + currentByte
		input = input / len(charMapping)
	}
	if len(currentByte) == 0 {
		currentByte = "0"
	}
	if len(currentByte) == 1 {
		currentByte = "0" + currentByte
	}
	return currentByte
}

func MakeReadable(input []byte) string {
	readable := ""
	for i := 0; i < len(input); i++ {
		readable += IntToChar(int(input[i]))
	}
	return readable
}

func Fetch(url string) (string, error) {
	var err error
	client := &http.Client{
		Timeout: time.Second * 1,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("User-Agent", "curl/7.38.0")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	clean := func(r rune) rune {
		switch {
		case r >= '0' && r <= '9':
			return r
		case r == '.':
			return r
		}
		return -1
	}

	return strings.Map(clean, string(body)), nil

}

func GetExternalIP() (net.IP, error) {
	var err error
	var body string

	body, err = Fetch("http://ifconfig.me")
	if err == nil {
		return net.ParseIP(body), nil
	}

	body, err = Fetch("http://ip.appspot.com")
	if err == nil {
		return net.ParseIP(body), nil
	}

	return nil, fmt.Errorf("No available external IP lookup service.")
}

// specific functions

func Connect(peerhash string) {
	log.Println("Connecting to", peerhash)
}

func GenerateHash(shortHash bool) string {
	peerHash := ""

	externalIP, err := GetExternalIP()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	log.Println("External", MakeReadable(externalIP.To4()))
	peerHash = MakeReadable(externalIP.To4())

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
		log.Printf("Found %v\n", MakeReadable(ip))
		peerHash += MakeReadable(ip)
	}
	return peerHash
}

func Listen() {
	_, err = net.Listen("tcp", ":6000")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	log.Println("Ready for connections")
}
