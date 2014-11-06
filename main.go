package main

import (
	"flag"
	"fmt"
	"github.com/brimstone/peerpipe/libpeerpipe"
	"log"
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

// [todo] - pass shorthash in via a config map
	peerpipe := libpeerpipe.New()
	log.Println("Peerhash:", peerpipe.GetHash())

	if len(args) == 1 {
		peerpipe.Connect(args[0])
	} else {
		log.Println("Already listening")
	}
}
