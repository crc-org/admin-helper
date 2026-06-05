package main

import (
	"fmt"
	"os"
	"os/user"
)

func cliCaller() string {
	currentUser, err := user.Current()
	if err != nil {
		return fmt.Sprintf("cli:uid:%d", os.Getuid())
	}
	return "cli:" + currentUser.Username
}
