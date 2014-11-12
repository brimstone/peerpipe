package libpeerpipe

import (
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var charMapping = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	"A", "B", "C", "D", "E", "F"}

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

func index(haystack []string, needle string) int {
	for i, straw := range haystack {
		if straw == needle {
			return i
		}
	}
	return -1
}

func CharToInt(input string) int {
	inputBytes := strings.Split(input, "")
	value := 0
	for i, r := range inputBytes {
		value += index(charMapping, r) * int(math.Pow(float64(len(charMapping)), float64(len(input)-1-i)))
	}
	return value
}

func MakeReadable(input []byte) string {
	readable := ""
	for i := 0; i < len(input); i++ {
		readable += IntToChar(int(input[i]))
	}
	return readable
}

func RemoveOneAddress(input []string, length int) ([]string, string) {
	// pop off 4 octets worth
	address := ""
	for i := 0; i < length-1; i++ {
		address += strconv.Itoa(CharToInt(strings.Join(input[0:2], ""))) + "."
		input = input[2:]
	}
	address += strconv.Itoa(CharToInt(strings.Join(input[0:2], "")))
	input = input[2:]
	return input, address
}

func Fetch(url string) (string, error) {
	var err error
	client := &http.Client{
		Timeout: time.Second * 5,
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
