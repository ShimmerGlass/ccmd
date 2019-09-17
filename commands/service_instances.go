package commands

import (
	"log"
	"reflect"

	"github.com/hashicorp/consul/api"
	"github.com/urfave/cli"
)

type serviceInstanceArgs struct {
	NodeName       string            `tmpl:"node_name" bexpr:"node_name"`
	ServiceName    string            `tmpl:"service_name" bexpr:"service_name"`
	ServiceTags    []string          `tmpl:"service_tags" bexpr:"service_tags"`
	NodeMeta       map[string]string `tmpl:"node_meta" bexpr:"node_meta"`
	InstanceID     string            `tmpl:"instance_id" bexpr:"instance_id"`
	InstanceHealth string            `tmpl:"instance_health" bexpr:"instance_health"`
}

func init() {
	addCommand(cli.Command{
		Name:      "service",
		Usage:     "Runs the commmand passed as argument for each service instances",
		UsageText: modelHelp(reflect.ValueOf(serviceInstanceArgs{})),
		ArgsUsage: "[CMD...]",
		Flags: []cli.Flag{
			dcsFlag,
			allDcsFlag,
			filterFlag,
			cli.StringSliceFlag{
				Name:  "service",
				Usage: "Service to target. Can be specified multiple times",
			},
			cli.BoolFlag{
				Name:  "all-services",
				Usage: "Target all services",
			},
		},
		Action: func(c *cli.Context) error {
			consul := getConsul(c)
			args := make(chan interface{})
			match := filter(c, serviceInstanceArgs{})

			go func() {

				dcs := getDcs(c, consul)

				for _, dc := range dcs {
					var services []string
					if c.Bool("all-services") {
						svcs, _, err := consul.Catalog().Services(&api.QueryOptions{
							Datacenter: dc,
						})
						if err != nil {
							log.Fatal(err)
						}
						for s := range svcs {
							services = append(services, s)
						}
					} else if len(c.StringSlice("service")) > 0 {
						services = c.StringSlice("service")
					} else {
						log.Fatal("No service specified, please use --service or --all-services")
					}

					for _, service := range services {
						svcs, _, err := consul.Health().Service(service, "", false, &api.QueryOptions{
							Datacenter: dc,
						})
						if err != nil {
							log.Fatal(err)
						}
						for _, d := range svcs {
							a := serviceInstanceArgs{
								NodeName:       d.Node.Node,
								NodeMeta:       d.Node.Meta,
								ServiceName:    d.Service.Service,
								ServiceTags:    d.Service.Tags,
								InstanceID:     d.Service.ID,
								InstanceHealth: d.Checks.AggregatedStatus(),
							}
							if match(a) {
								args <- a
							}
						}
					}
				}

				close(args)
			}()

			return run(c, c.Args(), args)
		},
	})

}
