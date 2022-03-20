package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	"golang.org/x/sys/unix"
	"sigs.k8s.io/yaml"
)

type Devloop struct {
	Template []DevloopTemplate
	Loop     []DevloopLoop
}

func (d Devloop) getLoop(name string) ([]DevloopAction, error) {
	loopIdx := -1
	for i, loop := range d.Loop {
		if loop.Name == name {
			loopIdx = i
			break
		}
	}
	if loopIdx == -1 {
		return nil, fmt.Errorf("loop %s not found", name)
	}

	templateName := d.Loop[loopIdx].TemplateRef.Name

	templateIdx := -1
	for i, template := range d.Template {
		if template.Name == templateName {
			templateIdx = i
			break
		}
	}
	if templateIdx == -1 {
		return nil, fmt.Errorf("template %s not found", templateName)
	}

	var replacements []string
	for k, v := range d.Loop[loopIdx].TemplateRef.Arg {
		replacements = append(replacements, "$"+k, v)
	}
	replacer := strings.NewReplacer(replacements...)

	var actions []DevloopAction
	for _, action := range d.Template[templateIdx].Action {
		a := DevloopAction{
			Cmd: make([]string, 0, len(action.Cmd)),
			Env: make(map[string]string),
		}
		for _, cmd := range action.Cmd {
			a.Cmd = append(a.Cmd, replacer.Replace(cmd))
		}
		for k, v := range action.Env {
			a.Env[k] = v
		}
		actions = append(actions, a)
	}

	return actions, nil
}

type DevloopTemplate struct {
	Name   string
	Action []DevloopAction
}

type DevloopAction struct {
	Cmd []string
	Env map[string]string
}

type DevloopLoop struct {
	Name        string
	TemplateRef DevloopTemplateRef
}

type DevloopTemplateRef struct {
	Name string
	Arg  map[string]string
}

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Fatalln("expected single arg")
	}

	loopName := flag.Arg(0)

	devloopRaw, err := os.ReadFile("devloop.yaml")
	if err != nil {
		log.Fatalln("read devloop.yaml", err)
	}

	var devloop Devloop
	err = yaml.Unmarshal(devloopRaw, &devloop)
	if err != nil {
		log.Fatalln("unmarshal devloop")
	}

	actions, err := devloop.getLoop(loopName)
	if err != nil {
		log.Fatalln("get loop", err)
	}

	errc := make(chan error, 1)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, unix.SIGINT, unix.SIGTERM)

	stdinc := make(chan string, 1)
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := os.Stdin.Read(buf)
			if err != nil {
				errc <- err
				return
			}
			stdinc <- string(buf[:n])
		}
	}()

	for {
		select {
		case <-errc:
			log.Fatalln(err)
		case <-sigc:
			log.Fatalln("got signal")
		case <-stdinc:
		}

		for _, action := range actions {
			c := exec.Command(action.Cmd[0], action.Cmd[1:]...)
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			err := c.Run()
			if err != nil {
				log.Fatalln("running cmd", action.Cmd[0])
			}
		}
	}
}
