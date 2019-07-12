package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/aestek/ccmd/tmpl"
)

type Options struct {
	Parralel int
}

func Run(cmd []string, instances []map[string]string, opts Options) error {
	if len(cmd) == 0 {
		return fmt.Errorf("command cannot be empty")
	}
	if opts.Parralel == 0 {
		opts.Parralel = 1
	}
	if opts.Parralel == -1 {
		opts.Parralel = len(instances)
	}

	c := make(chan map[string]string, len(instances))
	for _, i := range instances {
		c <- i
	}
	close(c)

	var wg sync.WaitGroup
	wg.Add(opts.Parralel)

	for i := 0; i < opts.Parralel; i++ {
		go func() {
			for i := range c {
				runOne(cmd, i)
			}
			wg.Done()
		}()
	}

	wg.Wait()

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
