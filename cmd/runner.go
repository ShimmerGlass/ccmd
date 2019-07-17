package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"strings"

	"github.com/aestek/ccmd/tmpl"
)

type Options struct {
	Parralel int
	NoPrefix bool
}

func Run(cmd []string, instances []interface{}, opts Options) error {
	if len(cmd) == 0 {
		return fmt.Errorf("command cannot be empty")
	}
	if opts.Parralel == 0 {
		opts.Parralel = 1
	}
	if opts.Parralel == -1 {
		opts.Parralel = len(instances)
	}

	cmdt, err := compileCmd(cmd)
	if err != nil {
		return err
	}

	instanceArgs := make([]map[string]string, len(instances))
	prefixes := make([]string, len(instances))
	maxPrefixLength := 0

	for i, inst := range instances {
		args, err := getArgs(inst)
		if err != nil {
			return err
		}
		instanceArgs[i] = args
	}

	if !opts.NoPrefix {
		for i, inst := range instanceArgs {
			a := cmdArgs(cmdt, inst)
			p := getPrefix(a)
			if len(p) > maxPrefixLength {
				maxPrefixLength = len(p)
			}
			prefixes[i] = p
		}
	}

	consoleWriterOut := &consoleWriter{
		inner:           os.Stdout,
		maxPrefixLength: maxPrefixLength,
	}
	consoleWriterErr := &consoleWriter{
		inner:           os.Stderr,
		maxPrefixLength: maxPrefixLength,
	}

	var wg sync.WaitGroup
	wg.Add(opts.Parralel)

	type instance struct {
		args map[string]string
		idx  int
	}
	c := make(chan instance)

	for i := 0; i < opts.Parralel; i++ {
		go func() {
			for i := range c {
				stdOut := &writerAdapter{inner: consoleWriterOut, prefix: prefixes[i.idx]}
				stdErr := &writerAdapter{inner: consoleWriterErr, prefix: prefixes[i.idx]}
				runOne(cmdt, i.args, stdOut, stdErr)
			}
			wg.Done()
		}()
	}

	go func() {
		for i, args := range instanceArgs {
			c <- instance{
				args: args,
				idx:  i,
			}
		}
		close(c)
	}()

	wg.Wait()

	return nil
}

func compileCmd(cmd []string) ([]*tmpl.Template, error) {
	res := make([]*tmpl.Template, len(cmd))
	for i, p := range cmd {
		t, err := tmpl.Parse(p)
		if err != nil {
			return nil, err
		}
		res[i] = t
	}
	return res, nil
}

func cmdArgs(cmd []*tmpl.Template, allArgs map[string]string) map[string]string {
	res := map[string]string{}
	for _, p := range cmd {
		for _, k := range p.Vars {
			v, ok := allArgs[k]
			if !ok {
				continue
			}
			res[k] = v
		}
	}
	return res
}

func runOne(cmdt []*tmpl.Template, args map[string]string, stdOut, stdErr io.Writer) {
	tcmd := []string{}
	for _, i := range cmdt {
		tcmd = append(tcmd, i.Exec(args))
	}

	if len(tcmd) == 1 && strings.ContainsAny(tcmd[0], " \t\n") {
		tcmd = []string{os.Getenv("SHELL"), "-c", tcmd[0]}
	}

	cmdArgs := []string{}
	if len(tcmd) > 1 {
		cmdArgs = tcmd[1:]
	}

	c := exec.Command(tcmd[0], cmdArgs...)
	c.Stdout = stdOut
	c.Stderr = stdErr

	env := os.Environ()
	for k, v := range args {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	c.Env = env

	err := c.Run()
	if err != nil {
		stdErr.Write([]byte(err.Error()))
	}
}
