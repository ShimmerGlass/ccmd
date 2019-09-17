package commands

import (
	"fmt"
	"log"
	"reflect"

	"github.com/aestek/ccmd/cmd"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-bexpr"
	"github.com/urfave/cli"
)

var (
	filterFlag = cli.StringFlag{
		Name:  "filter, f",
		Usage: "Filter result",
	}
	dcsFlag = cli.StringSliceFlag{
		Name:  "dc",
		Usage: "Datacenters to target (default the agent's datacenter)",
	}
	allDcsFlag = cli.BoolFlag{
		Name:  "all-dcs",
		Usage: "Target all datacenters (default the agent's datacenter)",
	}
)

var Commands []cli.Command

func addCommand(c cli.Command) {
	Commands = append(Commands, c)
}

func getConsul(c *cli.Context) *api.Client {
	consul, err := api.NewClient(&api.Config{
		Address: c.GlobalString("consul"),
	})
	if err != nil {
		log.Fatal(err)
	}

	return consul
}

func getRunOpts(c *cli.Context) cmd.Options {
	return cmd.Options{
		Parallel: c.GlobalInt("parallel"),
		NoPrefix: c.GlobalBool("no-prefix"),
	}
}

func run(c *cli.Context, command []string, args chan interface{}) error {
	return cmd.Run(command, args, getRunOpts(c))
}

func filter(c *cli.Context, model interface{}) func(interface{}) bool {
	if c.String("filter") == "" {
		return func(interface{}) bool {
			return true
		}
	}

	evaluator, err := bexpr.CreateEvaluatorForType(c.String("filter"), &bexpr.EvaluatorConfig{}, model)
	if err != nil {
		log.Fatal(err)
	}

	return func(d interface{}) bool {
		r, err := evaluator.Evaluate(d)
		if err != nil {
			log.Fatal(err)
		}
		return r
	}
}

func modelHelp(typ reflect.Value) string {
	res := "Variables : \n"
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Type().Field(i)
		name := f.Tag.Get("tmpl")
		kind := f.Type.Kind()
		help := f.Tag.Get("help")

		res += fmt.Sprintf("      {%s} (%s)", name, kind)
		if help != "" {
			res += ": " + help
		}
		res += "\n"
	}
	return res
}

func getDcs(c *cli.Context, consul *api.Client) []string {
	if c.Bool("all-dcs") {
		dcs, err := consul.Catalog().Datacenters()
		if err != nil {
			log.Fatal(err)
		}

		return dcs
	}
	if len(c.StringSlice("dc")) > 0 {
		return c.StringSlice("dc")
	}

	return []string{""}
}
