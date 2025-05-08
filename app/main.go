package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"unicode"
)

const (
	EXIT = "exit"
	ECHO = "echo"
	TYPE = "type"
	PWD = "pwd"
	CD = "cd"
)

var BUILT_IN_TYPES = map[string]struct{}{
	ECHO: {},
	EXIT: {},
	TYPE: {},
	PWD: {}, 
	CD: {}, 
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

func runExecutableFile(command string,  args []string) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Unable to execute the command: %v", err)
	}
	fmt.Print(string(output))
}

func extractCmdArgs(inputString string) (string, []string) {
	var current strings.Builder
	args := []string{}
	inSingleQuote := false
	inDoubleQuote := false
	escaped := false

	for _, c := range inputString {

		switch {
		case escaped:
			current.WriteRune(c)
			escaped = false

		case c == '\\' && !inDoubleQuote && !inSingleQuote:
			escaped = true

		case c == '\'' && !inDoubleQuote:
			inSingleQuote = !inSingleQuote

		case c == '"' && !inSingleQuote:
			inDoubleQuote = !inDoubleQuote
		
		case unicode.IsSpace(c) && !inSingleQuote && !inDoubleQuote: 
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		
		default:
			current.WriteRune(c)
		
		}
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	if len(args) == 1 {
		return args[0], []string{}
	}

	return args[0], args[1:]
}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		inputCommand, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Printf("Unable to read input command: %s\n", inputCommand)
		}

		command, args := extractCmdArgs(inputCommand)

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
		
		case PWD:
			pwd, err := os.Getwd()
			if err != nil {
				fmt.Printf("Something went wrong: %v\n", err)
				break
			}
			fmt.Println(pwd)

		case CD:
			if len(args) < 1 {
				fmt.Printf("Missing args: %v\n", err)
				break
			}

			dir := args[0]
			if args[0] == "~" {
				dir, err = os.UserHomeDir()
				if err != nil {
					fmt.Printf("Something went wrong: %v", err)
					break
				}
			}
			
			if err := os.Chdir(dir); err != nil {
				fmt.Printf("cd: %s: No such file or directory\n", dir)
				break
			}

		default:
			foundFile := false
			for p := range strings.SplitSeq(PATH, ":") {
					if findExecutableFile(p, command) == nil {
						runExecutableFile(command, args)	
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
