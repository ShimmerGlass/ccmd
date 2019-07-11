package main

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

func ProviderService(c *api.Client, service, dc string) Provider {
	return func() ([]map[string]string, error) {
		svcs, _, err := c.Health().Service(service, "", false, &api.QueryOptions{
			Datacenter: dc,
		})
		if err != nil {
			return nil, fmt.Errorf("error fetching dcs: %s", err)
		}

		args := []map[string]string{}
		for _, d := range svcs {
			args = append(args, map[string]string{
				"node": d.Node.Node,
			})
		}

		return args, nil
	}
}
