package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aestek/ccmd/cmd"
	"github.com/aestek/ccmd/tmpl"
	"github.com/hashicorp/consul/api"
	"github.com/urfave/cli"
)

func getConsul(c *cli.Context) *api.Client {
	addr := c.GlobalString("consul")

	if c.GlobalString("consul-addr-template") != "" {
		vars := map[string]string{}
		for _, e := range os.Environ() {
			pair := strings.Split(e, "=")
			vars[pair[0]] = pair[1]
		}

		var err error
		addr, err = tmpl.Exec(c.GlobalString("consul-addr-template"), vars)
		if err != nil {
			log.Fatalf("error templating consul-addr-template: %s", err)
		}
	}

	consul, err := api.NewClient(&api.Config{
		Address: addr,
	})
	if err != nil {
		log.Fatal(err)
	}

	return consul
}

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "consul, c",
			Value:  "http://127.0.0.1:8500",
			Usage:  "Consul agent address",
			EnvVar: "CONSUL_ADDR",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "dc",
			Usage: "For each datacenter",
			Action: func(c *cli.Context) error {
				consul := getConsul(c)

				args, err := ProviderDC(consul)()
				if err != nil {
					return err
				}

				return cmd.Run([]string(c.Args()), args)
			},
		},
		{
			Name:  "service",
			Usage: "For each service instance",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "service, s",
					Usage: "Service name",
				},
				cli.StringFlag{
					Name:  "dc",
					Usage: "Datacenter",
				},
			},
			Action: func(c *cli.Context) error {
				consul := getConsul(c)

				service := c.String("service")
				if service == "" {
					return fmt.Errorf("must provide service")
				}

				args, err := ProviderService(consul, service, c.String("dc"))()
				if err != nil {
					return err
				}

				return cmd.Run([]string(c.Args()), args)
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
