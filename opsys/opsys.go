package opsys

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/vrmiguel/shego/cli"
	"github.com/vrmiguel/shego/utils"
)

// UserData holds information about the current user
type UserData struct {
	Username  string // Username
	Hostname  string // Hostname
	HomeDir   string // Home directory
	CurrentWD string // Current working directory
	PrettyCWD string // Current working directory with /home/dir translated to ~
}

// GetUserData gathers info about the user and saves it into a UserData object
func GetUserData() UserData {
	var userdata UserData
	curUser, err := user.Current()
	utils.AssertNonNil(err)

	userdata.Username = curUser.Username
	cwd, err := filepath.Abs("./")
	utils.AssertNonNil(err)

	userdata.CurrentWD = cwd
	curHost, err := os.Hostname()
	utils.AssertNonNil(err)

	hmd, err := os.UserHomeDir()
	utils.AssertNonNil(err)

	userdata.HomeDir = hmd
	userdata.Hostname = curHost
	userdata.PrettyCWD = strings.Replace(cwd, hmd, "~", 1)
	return userdata
}

func updateCwd(ud *UserData) {
	cwd, err := filepath.Abs("./")
	utils.AssertNonNil(err)
	ud.CurrentWD = cwd
	ud.PrettyCWD = strings.Replace(cwd, ud.HomeDir, "~", 1)
}

// ShowPrompt returns the shell's prompt (username@hostname:cwd$)
func ShowPrompt(ud UserData) {
	fmt.Printf("%s%s@%s%s:%s%s$ ", cli.AnsiGreen, ud.Username, ud.Hostname, cli.AnsiBlue, ud.PrettyCWD, cli.AnsiReset)
}

func changeDirectory(tokens []string, ud *UserData) {
	if len(tokens) > 2 {
		fmt.Fprintf(os.Stderr, "shego: cd: too many arguments\n")
		return
	}
	if len(tokens) == 1 || tokens[1] == "~" || tokens[1] == "$HOME" {
		err := os.Chdir(ud.HomeDir)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		err := os.Chdir(tokens[1])
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func runSimpleCommand(tokens []string) error {
	var simpleCommand *exec.Cmd
	simpleCommand = exec.Command(tokens[0], tokens[1:]...)
	res, err := simpleCommand.Output()
	print(string(res))
	if err != nil {
		fmt.Fprintf(os.Stderr, "shego: problem running '%s'\n", simpleCommand.Args)
	}
	return nil
}

// ParseCommand TODO description
func ParseCommand(line string, ud *UserData) {
	line = line[:len(line)-1] // Remove newline
	tokens := strings.Split(line, " ")
	if tokens[0] == "cd" {
		changeDirectory(tokens, ud)
		updateCwd(ud)
	} else if tokens[0] == "exit" {
		os.Exit(0)
	} else {
		runSimpleCommand(tokens)
	}
}
