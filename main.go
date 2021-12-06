package main

import (
	"flag"
	"fmt"
	"monkey/repl"
	"os"
	"os/user"
)

func main() {
	flag.Parse()
	filename := flag.Args()
	switch len(filename) {
	case 0:
		user, err := user.Current()
		if err != nil {
			panic(err)
		}
		fmt.Printf("hello %s! This is the Monkey programming language!\n", user.Username)
		fmt.Printf("Feel free to type in commands\n")
		repl.ReplFromCommandLine(os.Stdin, os.Stdout)

	case 1:
		fp, err := os.Open(filename[0])
		if err != nil {
			panic(err)
		}
		repl.ReplFromFile(fp, os.Stdin, os.Stdout)
		fp.Close()

	default:
		fmt.Printf("Argument length must be 1. got=%d\n", len(filename))
		return
	}

}
