package cli

import (
	"fmt"
	"os"
)

const usage = "Usage: ./shego [-n, --nohup] [-V, --verbose] [-v, --version]\n"

// AnsiReset is the ANSI escape code that resets the terminal text back to its default color
const AnsiReset = "\033[0m"

// AnsiGreen is the ANSI escape code that produces text in green
const AnsiGreen = "\033[32m"

// AnsiBlue is the ANSI escape code that produces text in blue
const AnsiBlue = "\033[34m"

// Version stores the current Shego version
const Version = "v.0.1-alpha"

// Args holds the command-line options given by the user
type Args struct {
	IgnoreHangups bool
	IsVerbose     bool
}

func printHelp() {
	fmt.Printf(usage)
	fmt.Printf("%-16s\tShow this message and exit.\n", "-h, --help")
	fmt.Printf("%-16s\tIgnore SIGHUPs.\n", "-n, --nohup")
	fmt.Printf("%-16s\tRun the shell in verbose mode.\n", "-V, --verbose")
	fmt.Printf("%-16s\tDisplay shego version and exit.\n", "-v, --version")
}

// ParseArgs reads through the given  arg. list and builds a Args struct
func ParseArgs(args []string) Args {
	if len(args) == 1 {
		return Args{false, false}
	}
	cfg := Args{false, false}
	for i := 1; i < len(args); i++ {
		arg := args[i]
		if arg == "-n" || arg == "--nohup" {
			cfg.IgnoreHangups = true
		} else if arg == "-h" || arg == "--help" {
			printHelp()
			os.Exit(0)
		} else if arg == "-v" || arg == "--version" {
			fmt.Printf("shego %s\n", Version)
			os.Exit(0)
		} else if arg == "-V" || arg == "--verbose" {
			cfg.IsVerbose = true
		} else {
			fmt.Fprintf(os.Stderr, "error: unknown option \"%s\"\n", arg)
			os.Exit(0)
		}
	}
	return cfg
}
