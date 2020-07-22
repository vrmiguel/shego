package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/vrmiguel/shego/cli"
	"github.com/vrmiguel/shego/opsys"
)

func signalHandler(q chan bool, cfg cli.Args) {
	var quit bool

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	if cfg.IsVerbose {
		fmt.Fprintf(os.Stderr, "shego-verbose: signal handler channel running.\n")
	}

	for signal := range c {
		switch signal {
		case syscall.SIGINT:
			fmt.Fprintf(os.Stderr, "SIGINT received.\n")
			quit = true
		case syscall.SIGTERM:
			fmt.Fprintf(os.Stderr, "SIGTERM received.\n")
			quit = true
		case syscall.SIGHUP:
			fmt.Fprintf(os.Stderr, "SIGHUP received.\n")
			if !cfg.IgnoreHangups {
				quit = true
			}
		}

		if quit {
			quit = false
			os.Exit(0)
		}
		q <- quit
	}
}

func main() {
	// init two channels, one for the signals, one for the main loop
	sig := make(chan bool)   // Signal handler channel
	loop := make(chan error) // Main loop (REPL) channel

	cfg := cli.ParseArgs(os.Args)
	if cfg.IsVerbose {
		fmt.Printf("shego-verbose: running version %s\n", cli.Version)
	}

	go signalHandler(sig, cfg) // Start signal handler channel

	fmt.Println("shego --- github.com/vrmiguel/shego")
	for quit := false; !quit; {
		//var line string
		ud := opsys.GetUserData()
		in := bufio.NewReader(os.Stdin)
		go func() {
			opsys.ShowPrompt(ud)
			line, err := in.ReadString('\n')
			// fmt.Scanln(&line)
			if err == io.EOF {
				fmt.Fprintf(os.Stdout, "^D received.\n")
				quit = true
			}
			if line != "" {
				opsys.ParseCommand(line, &ud)
			}
			loop <- nil
		}()

		select {
		case quit = <-sig:
		case <-loop:
		}
	}
}
