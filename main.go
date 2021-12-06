package main

import (
	"bufio"
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
		repl.Start(os.Stdin, os.Stdout)

	case 1:
		fp, err := os.Open(filename[0])
		if err != nil {
			panic(err)
		}

		scanner := bufio.NewScanner(fp)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		fp.Close()

	default:
		fmt.Printf("Argument length must be 1. got=%d\n", len(filename))
		return
	}

}
