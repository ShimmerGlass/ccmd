package commands

import (
	"fmt"
	"reflect"
	"time"

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
		ArgsUsage: "[SERVICE] [CMD...]",
		Flags: []cli.Flag{
			dcFlag,
			filterFlag,
			watchFlag,
		},
		Action: func(c *cli.Context) error {
			consul := getConsul(c)

			if len(c.Args()) == 0 {
				return fmt.Errorf("please provide the service as first argument")
			}

			service := c.Args()[0]

			var idx uint64
			for {
				svcs, meta, err := consul.Health().Service(service, "", false, &api.QueryOptions{
					WaitIndex:  idx,
					WaitTime:   10 * time.Minute,
					Datacenter: c.String("dc"),
				})
				if err != nil {
					return err
				}
				idx = meta.LastIndex

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

				err = run(c, c.Args()[1:], args)
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
