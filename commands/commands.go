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
	dcFlag = cli.StringFlag{
		Name:  "dc",
		Usage: "Datacenter to target (default the agent's datacenter)",
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
		Parralel: c.GlobalInt("parallel"),
		NoPrefix: c.GlobalBool("no-prefix"),
	}
}

func run(c *cli.Context, args []interface{}) error {
	return cmd.Run([]string(c.Args()), args, getRunOpts(c))
}

func filter(c *cli.Context, args []interface{}, model interface{}) ([]interface{}, error) {
	if c.String("filter") == "" {
		return args, nil
	}
	f, err := bexpr.CreateFilter(c.String("filter"), &bexpr.EvaluatorConfig{}, model)
	if err != nil {
		return nil, err
	}

	argsif, err := f.Execute(args)
	if err != nil {
		return nil, err
	}

	return argsif.([]interface{}), nil
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
