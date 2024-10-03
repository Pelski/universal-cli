package main

import (
	"fmt"
	"os"
	"strings"
)

var debug bool

func main() {
	args := os.Args[1:]
	configPath := ""
	var newArgs []string

	// Parse command-line arguments
	i := 0
	for i < len(args) {
		arg := args[i]
		switch {
		case arg == "--config":
			if i+1 < len(args) {
				configPath = args[i+1]
				i += 2
			} else {
				fmt.Println("Missing value for --config")
				os.Exit(1)
			}
		case strings.HasPrefix(arg, "--config="):
			configPath = strings.TrimPrefix(arg, "--config=")
			i++
		case arg == "--ucli-debug":
			debug = true
			i++
		case strings.HasPrefix(arg, "--ucli-debug="):
			value := strings.TrimPrefix(arg, "--ucli-debug=")
			debug = value == "true" || value == "1"
			i++
		default:
			newArgs = append(newArgs, arg)
			i++
		}
	}
	args = newArgs

	InitConfig(configPath)

	if len(args) < 1 {
		fmt.Println("You need to provide an operation, e.g., get, create, update, delete")
		os.Exit(1)
	}

	operation := args[0]
	var resources []string
	var flagsArgs []string

	// Separate resources and flags
	for i := 1; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			flagsArgs = args[i:]
			break
		} else {
			resources = append(resources, arg)
		}
	}

	flags := ParseDynamicFlags(flagsArgs)

	if debug {
		fmt.Println("> Operation:", operation)
		fmt.Println("> Resources:", resources)
		fmt.Println("> Flags:", flags)
	}

	HandleOperation(operation, resources, flags)
}
