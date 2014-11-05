package main

import (
	"flag"
	"fmt"
	"github.com/brimstone/peerpipe/libpeerpipe"
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
		libpeerpipe.Connect(args[0])
	} else {
		libpeerpipe.Listen()
	}
}