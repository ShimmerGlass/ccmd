package commands

import (
	"reflect"

	"github.com/urfave/cli"
)

type dcArgs struct {
	DC string `tmpl:"dc"`
}

func init() {
	addCommand(cli.Command{
		Name:      "dc",
		Usage:     "Runs the commmand passed as argument for each datacenter",
		UsageText: modelHelp(reflect.ValueOf(dcArgs{})),
		Action: func(c *cli.Context) error {
			consul := getConsul(c)

			dcs, err := consul.Catalog().Datacenters()
			if err != nil {
				return err
			}

			args := make(chan interface{}, len(dcs))
			for _, d := range dcs {
				args <- dcArgs{
					DC: d,
				}
			}

			return run(c, c.Args(), args)
		},
	})
}
