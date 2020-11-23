package main

import (
	"os"

	"github.com/goodhosts/cli/cmd"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:     "goodhosts",
		Usage:    "manage your hosts file goodly",
		Commands: cmd.Commands(),
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
