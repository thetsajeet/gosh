package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	EXIT = "exit"
	ECHO = "echo"
	TYPE = "type"
)

var BUILT_IN_TYPES = map[string]struct{}{
	ECHO: {},
	EXIT: {},
	TYPE: {},
} 

var PATH = os.Getenv("PATH")

func findExecutableFile(path, executableType string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, f := range files {
		if executableType == f.Name() {
			return nil
		}
	}

	return errors.New("unable to find the file") 
}

func runExecutableFile(path string, args []string) {
	cmd := exec.Command(path, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Unable to execute the command: %v", err)
	}
	fmt.Println(output)
}

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

		case ECHO:
			fmt.Printf("%s\n", strings.Join(args, " "))

		case TYPE:
			if len(args) < 1 {
				fmt.Printf("Require an arg to check type: %v\n", args)
				break
			} 
			if _, exists := BUILT_IN_TYPES[args[0]]; exists {
				fmt.Printf("%s is a shell builtin\n", args[0])
			} else {
				foundFile := false
				for p := range strings.SplitSeq(PATH, ":") {
					if findExecutableFile(p, args[0]) == nil {
						fmt.Printf("%s is %s\n", args[0], p + "/" + args[0])	
						foundFile = true
						break
					}
				}
				if !foundFile {
					fmt.Printf("%s: not found\n", args[0])
				}

			} 

		default:
			foundFile := false
			for p := range strings.SplitSeq(PATH, ":") {
					if findExecutableFile(p, command) == nil {
						runExecutableFile(p, args)	
						foundFile = true
						break
					}
				}
			if !foundFile {
				fmt.Printf("%s: command not found\n", command)
			}
		}
	}
}
