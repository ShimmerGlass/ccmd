package main

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

func ProviderDC(c *api.Client) Provider {
	return func() ([]map[string]string, error) {
		dcs, err := c.Catalog().Datacenters()
		if err != nil {
			return nil, fmt.Errorf("error fetching dcs: %s", err)
		}

		args := []map[string]string{}
		for _, d := range dcs {
			args = append(args, map[string]string{
				"dc": d,
			})
		}

		return args, nil
	}
}
