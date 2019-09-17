package commands

import (
	"log"
	"reflect"

	"github.com/hashicorp/consul/api"
	"github.com/urfave/cli"
)

type nodeArgs struct {
	Name string            `tmpl:"name" bexpr:"name"`
	Meta map[string]string `tmpl:"meta" bexpr:"meta"`
}

func init() {
	addCommand(cli.Command{
		Name:      "nodes",
		Usage:     "Runs the commmand passed as argument for each node",
		UsageText: modelHelp(reflect.ValueOf(serviceInstanceArgs{})),
		ArgsUsage: "[CMD...]",
		Flags: []cli.Flag{
			dcsFlag,
			allDcsFlag,
			filterFlag,
		},
		Action: func(c *cli.Context) error {
			consul := getConsul(c)
			args := make(chan interface{})
			match := filter(c, nodeArgs{})

			dcs := getDcs(c, consul)

			go func() {
				for _, dc := range dcs {
					nodes, _, err := consul.Catalog().Nodes(&api.QueryOptions{
						Datacenter: dc,
					})
					if err != nil {
						log.Fatal(err)
					}

					for _, d := range nodes {
						a := nodeArgs{
							Name: d.Node,
							Meta: d.Meta,
						}
						if match(a) {
							args <- a
						}
					}
				}

				close(args)
			}()

			return run(c, c.Args(), args)
		},
	})

}
