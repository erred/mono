package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

var (
	lineNumberRe = regexp.MustCompile(`^(\d+):(\d+):.*`)
	ansi         = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
)

func main() {
	cmd := exec.Command("rg", append([]string{"--heading", "--column", "--color=always"}, os.Args[1:]...)...)
	cmd.Stderr = os.Stderr
	rc, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	defer rc.Close()
	rd := bufio.NewReaderSize(rc, 64*1024*1024)

	f, err := os.Create("/tmp/t_aliases")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var (
		curPath string
		idx     int = 1
	)

	cmd.Start()
	defer cmd.Wait()

	for {
		coloredLine, err := rd.ReadBytes('\n')
		if err != nil {
			break
		}

		line := ansi.ReplaceAll(coloredLine, nil)
		if curPath == "" {
			line = line[:len(line)-1]
			curPath, err = filepath.Abs(string(line))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s", coloredLine)
		} else if groupIdxs := lineNumberRe.FindSubmatchIndex(line); len(groupIdxs) > 0 {
			// _, err := fmt.Fprintf(f, "alias e%d='nvim %s +%s'\n", idx, curPath, string(line[groupIdxs[2]:groupIdxs[3]]))
			_, err := fmt.Fprintf(f, `alias e%d='nvim -c "call cursor(%s, %s)" "%s"'`+"\n", idx, string(line[groupIdxs[2]:groupIdxs[3]]), string(line[groupIdxs[4]:groupIdxs[5]]), curPath)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("\x1b[34m[\x1b[31m%d\x1b[34m]\x1b[0m %s", idx, coloredLine)
			idx++
		} else {
			curPath = ""
			fmt.Println()
		}
	}
}
