package main

import (
	"os"

	"github.com/code-ready/admin-helper/cmd"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:     "admin-helper",
		Usage:    "manage your hosts file goodly",
		Commands: cmd.Commands(),
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
