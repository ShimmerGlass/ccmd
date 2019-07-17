package commands

import (
	"reflect"

	"github.com/hashicorp/consul/api"
	"github.com/urfave/cli"
)

type serviceInstanceArgs struct {
	NodeName       string   `tmpl:"node_name" bexpr:"node_name"`
	ServiceName    string   `tmpl:"service_name" bexpr:"service_name"`
	ServiceTags    []string `tmpl:"service_tags" bexpr:"service_tags"`
	InstanceID     string   `tmpl:"instance_id" bexpr:"instance_id"`
	InstanceHealth string   `tmpl:"instance_health" bexpr:"instance_health"`
}

func init() {
	addCommand(cli.Command{
		Name:      "service",
		Usage:     "Runs the commmand passed as argument for each service instances",
		UsageText: modelHelp(reflect.ValueOf(serviceInstanceArgs{})),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "service, s",
				Usage: "Service to get instances from",
			},
			dcFlag,
			filterFlag,
		},
		Action: func(c *cli.Context) error {
			consul := getConsul(c)

			svcs, _, err := consul.Health().Service(c.String("service"), "", false, &api.QueryOptions{
				Datacenter: c.String("dc"),
			})
			if err != nil {
				return err
			}

			args := []interface{}{}
			for _, d := range svcs {
				args = append(args, serviceInstanceArgs{
					NodeName:       d.Node.Node,
					ServiceName:    d.Service.Service,
					ServiceTags:    d.Service.Tags,
					InstanceID:     d.Service.ID,
					InstanceHealth: d.Checks.AggregatedStatus(),
				})
			}

			args, err = filter(c, args, []serviceInstanceArgs{})
			if err != nil {
				return err
			}

			return run(c, args)
		},
	})

}
