package commands

import (
	"reflect"

	"github.com/urfave/cli"
)

type agentServiceArgs struct {
	ServiceName string `tmpl:"service_name"`
	ServiceID   string `tmpl:"service_id"`
	ServicePort int    `tmpl:"service_port"`
}

func init() {
	addCommand(cli.Command{
		Name:      "agent-services",
		Usage:     "Runs the commmand passed as argument for agent services",
		UsageText: modelHelp(reflect.ValueOf(agentServiceArgs{})),
		ArgsUsage: "Command to run",
		Action: func(c *cli.Context) error {
			consul := getConsul(c)

			svcs, err := consul.Agent().Services()
			if err != nil {
				return err
			}

			args := make(chan interface{}, len(svcs))
			for _, d := range svcs {
				args <- agentServiceArgs{
					ServiceName: d.Service,
					ServiceID:   d.ID,
					ServicePort: d.Port,
				}
			}

			return run(c, c.Args(), args)
		},
	})

}
