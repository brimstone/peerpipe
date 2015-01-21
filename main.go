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
	peerpipe, err := libpeerpipe.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	log.Println("Peerhash:", peerpipe.GetHash())

	// fi, _ := os.Stdin.Stat() // get the FileInfo struct describing the standard input.
	// if fi.Mode() & os.ModeCharDevice {
	// pipe
	// else
	// interactive

	if len(args) == 1 {
		// We're the client
		go peerpipe.Connect(args[0])
	} else {
		// We're the server
		log.Println("Already listening")
	}
	peerpipe.Wait()
	log.Println("Done waiting.")
}
