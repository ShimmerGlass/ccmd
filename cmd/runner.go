package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"strings"

	"github.com/prometheus/common/log"
	"github.com/shimmerglass/ccmd/tmpl"
)

type Options struct {
	Parallel int
	NoPrefix bool
}

func Run(cmd []string, instances chan interface{}, opts Options) error {
	if len(cmd) == 0 {
		return runPrint(cmd, instances, opts)
	}

	if opts.Parallel == 0 {
		opts.Parallel = 1
	}

	cmdt, err := compileCmd(cmd)
	if err != nil {
		return err
	}

	writer := &consoleWriter{}

	var wg sync.WaitGroup

	var sem chan struct{}
	if opts.Parallel > 0 {
		sem = make(chan struct{}, opts.Parallel)
		for i := 0; i < opts.Parallel; i++ {
			sem <- struct{}{}
		}
	}

	for instance := range instances {
		if opts.Parallel > 0 {
			<-sem
		}
		wg.Add(1)
		go func(instance interface{}) {
			defer func() {
				wg.Done()
				sem <- struct{}{}
			}()

			args, err := getArgs(instance)
			if err != nil {
				log.Error(err)
				return
			}

			a := cmdArgs(cmdt, args)
			p := getVarsDesc(a, true)

			if len(p) > writer.maxPrefixLength {
				writer.maxPrefixLength = len(p)
			}

			var stdOut io.Writer = os.Stdout
			var stdErr io.Writer = os.Stderr

			if !opts.NoPrefix {
				stdOut = &writerAdapter{target: os.Stdout, inner: writer, prefix: p}
				stdErr = &writerAdapter{target: os.Stderr, inner: writer, prefix: p}
			}

			runOne(cmdt, args, stdOut, stdErr)
		}(instance)
	}

	wg.Wait()

	return nil
}

func runPrint(cmd []string, instances chan interface{}, opts Options) error {
	for inst := range instances {
		args, err := getArgs(inst)
		if err != nil {
			return err
		}
		fmt.Println(getVarsDesc(args, false))
	}

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
