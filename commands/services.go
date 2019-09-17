package commands

import (
	"log"
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
			dcsFlag,
			allDcsFlag,
			filterFlag,
		},
		Action: func(c *cli.Context) error {
			consul := getConsul(c)
			args := make(chan interface{})
			match := filter(c, serviceArgs{})

			go func() {
				for _, dc := range getDcs(c, consul) {
					svcs, _, err := consul.Catalog().Services(&api.QueryOptions{
						Datacenter: dc,
					})
					if err != nil {
						log.Fatalf("error fetching catalog: %s", err)
					}

					for s, tags := range svcs {
						a := serviceArgs{
							ServiceName: s,
							ServiceTags: tags,
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
