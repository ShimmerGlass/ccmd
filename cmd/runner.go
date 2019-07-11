package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/aestek/ccmd/tmpl"
)

func Run(cmd []string, instances []map[string]string) error {
	if len(cmd) == 0 {
		return fmt.Errorf("command cannot be empty")
	}

	for _, i := range instances {
		err := runOne(cmd, i)
		if err != nil {
			return err
		}
	}

	return nil
}

func runOne(cmd []string, args map[string]string) error {
	usedArgs := []string{}

	tcmd := []string{}
	for _, i := range cmd {
		t, err := tmpl.Parse(i)
		if err != nil {
			return err
		}
		usedArgs = append(usedArgs, t.Vars...)
		tcmd = append(tcmd, t.Exec(args))
	}

	cmdArgs := []string{}
	if len(tcmd) > 1 {
		cmdArgs = tcmd[1:]
	}

	c := exec.Command(tcmd[0], cmdArgs...)
	c.Stdout = newIOWrapper(args, usedArgs, 1, os.Stdout)
	c.Stderr = newIOWrapper(args, usedArgs, 1, os.Stderr)

	env := os.Environ()
	for k, v := range args {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	c.Env = env

	return c.Run()
}
