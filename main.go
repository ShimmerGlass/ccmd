package main

import (
	"log"
	"os"

	"github.com/aestek/ccmd/commands"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "consul, c",
			Value:  "http://127.0.0.1:8500",
			Usage:  "Consul agent address",
			EnvVar: "CONSUL_ADDR",
		},
		cli.IntFlag{
			Name:  "parallel, p",
			Value: 1,
			Usage: "Run N commands in parallel. -1 for all commands in parallel.",
		},
		cli.BoolFlag{
			Name:  "no-prefix, n",
			Usage: "DO not prefix output by used variables",
		},
	}

	app.Commands = commands.Commands

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
