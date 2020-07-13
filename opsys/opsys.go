package opsys

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"../cli"
	"../utils"
)

// UserData holds information about the current user
type UserData struct {
	Username  string // Username
	Hostname  string // Hostname
	Hmd       string // Home directory
	Cwd       string // Current working directory
	PrettyCwd string // Current working directory with /home/dir translated to ~
}

// GetUserData gathers info about the user and saves it into a UserData object
func GetUserData() UserData {
	var userdata UserData
	curUser, err := user.Current()
	utils.AssertNonNil(err)

	userdata.Username = curUser.Username
	cwd, err := filepath.Abs("./")
	utils.AssertNonNil(err)

	userdata.Cwd = cwd
	curHost, err := os.Hostname()
	utils.AssertNonNil(err)

	hmd, err := os.UserHomeDir()
	utils.AssertNonNil(err)

	userdata.Hmd = hmd
	userdata.Hostname = curHost
	userdata.PrettyCwd = strings.Replace(cwd, hmd, "~", 1)
	return userdata
}

func updateCwd(ud *UserData) {
	cwd, err := filepath.Abs("./")
	utils.AssertNonNil(err)
	ud.Cwd = cwd
	ud.PrettyCwd = strings.Replace(cwd, ud.Hmd, "~", 1)
}

// ShowPrompt returns the shell's prompt (username@hostname:cwd$)
func ShowPrompt(ud UserData) {
	fmt.Printf("%s%s@%s%s:%s%s$ ", cli.AnsiGreen, ud.Username, ud.Hostname, cli.AnsiBlue, ud.PrettyCwd, cli.AnsiReset)
}

func changeDirectory(tokens []string, ud *UserData) {
	if len(tokens) > 2 {
		fmt.Fprintf(os.Stderr, "shego: cd: too many arguments\n")
		return
	}
	if len(tokens) == 1 || tokens[1] == "~" || tokens[1] == "$HOME" {
		err := os.Chdir(ud.Hmd)
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

// ParseCommand TODO description
func ParseCommand(line string, ud *UserData) {
	tokens := strings.Split(line, " ")
	re := regexp.MustCompile("\\n")
	for i := 0; i < len(tokens); i++ {
		tokens[i] = re.ReplaceAllString(tokens[i], "")
	}
	if tokens[0] == "cd" {
		changeDirectory(tokens, ud)
		updateCwd(ud)
	}

}
