package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			log.Fatalf("Unable to read input command: %s", command)
		}

		fmt.Printf("%s: command not found\n", command[:len(command)-1])
	}
}
