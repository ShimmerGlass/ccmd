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
		Action: func(c *cli.Context) error {
			consul := getConsul(c)

			svcs, err := consul.Agent().Services()
			if err != nil {
				return err
			}

			args := []interface{}{}
			for _, d := range svcs {
				args = append(args, agentServiceArgs{
					ServiceName: d.Service,
					ServiceID:   d.ID,
					ServicePort: d.Port,
				})
			}

			return run(c, args)
		},
	})

}
