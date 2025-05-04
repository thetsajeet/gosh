package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const EXIT = "exit"

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		inputCommand, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Printf("Unable to read input command: %s\n", inputCommand)
		}

		inputCommandParts := strings.Fields(inputCommand)
		command, args := inputCommandParts[0], inputCommandParts[1:]

		switch command {
		case EXIT:
			if len(args) < 1 {
				fmt.Printf("Exit code required as an arg: %v\n", args)
				break
			}
			exitCode, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Printf("Error converting string to number: %v\n", err)
				break
			}
			os.Exit(exitCode)
		default:
			fmt.Printf("%s: command not found\n", command)
		}
	}
}
