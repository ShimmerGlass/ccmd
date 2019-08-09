package commands

import (
	"reflect"
	"time"

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
			dcFlag,
			filterFlag,
			watchFlag,
		},
		Action: func(c *cli.Context) error {
			consul := getConsul(c)

			var idx uint64
			for {
				nodes, meta, err := consul.Catalog().Nodes(&api.QueryOptions{
					WaitIndex:  idx,
					WaitTime:   10 * time.Minute,
					Datacenter: c.String("dc"),
				})
				if err != nil {
					return err
				}
				idx = meta.LastIndex

				args := []interface{}{}
				for _, d := range nodes {
					args = append(args, nodeArgs{
						Name: d.Node,
						Meta: d.Meta,
					})
				}

				args, err = filter(c, args, []nodeArgs{})
				if err != nil {
					return err
				}

				err = run(c, c.Args(), args)
				if err != nil {
					return err
				}
				if !c.Bool("watch") {
					break
				}
			}

			return nil
		},
	})

}
