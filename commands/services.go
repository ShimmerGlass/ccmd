package commands

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/consul/api"
	"github.com/urfave/cli"
)

type serviceArgs struct {
	ServiceName string   `tmpl:"service_name" bexpr:"service_name"`
	ServiceTags []string `tmpl:"service_tags" bexpr:"service_tags"`
}

func init() {
	addCommand(cli.Command{
		Name:      "services",
		Usage:     "Runs the commmand passed as argument for catalog service",
		UsageText: modelHelp(reflect.ValueOf(serviceArgs{})),
		Flags: []cli.Flag{
			dcFlag,
			filterFlag,
		},
		Action: func(c *cli.Context) error {
			consul := getConsul(c)

			svcs, _, err := consul.Catalog().Services(&api.QueryOptions{
				Datacenter: c.String("dc"),
			})
			if err != nil {
				return fmt.Errorf("error fetching dcs: %s", err)
			}

			args := []interface{}{}
			for s, tags := range svcs {
				args = append(args, serviceArgs{
					ServiceName: s,
					ServiceTags: tags,
				})
			}

			args, err = filter(c, args, []serviceArgs{})
			if err != nil {
				return err
			}

			return run(c, args)
		},
	})

}
