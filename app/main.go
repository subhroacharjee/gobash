package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/codecrafters-io/shell-starter-go/internal/commands"
)

func main() {
	run()
}

func run() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input

		raw, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v: error occured\n", err)
		}

		_, err = commands.RunCommand(raw)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
	}
}
