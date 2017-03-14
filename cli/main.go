package cli

import (
	"fmt"
	"os"

	"github.com/micromdm/squirrel/version"
)

func Main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}
	switch os.Args[1] {
	case "serve":
		Serve()
		return
	case "help", "-h", "--help":
		printHelp()
		return
	case "version":
		version.Print()
		return
	default:
		fmt.Printf("no such command")
		return
	}
}

const usage = `
Usage: 
	squirrel <COMMAND>

Available Commands:
	help
	serve
	version

Use squirel <command> -h for additional usage of each command.
Example: squirrel serve -h
`

func printHelp() {
	fmt.Println(usage)
}
