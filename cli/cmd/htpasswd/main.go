package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	fl := flag.NewFlagSet("htpasswd", flag.ExitOnError)
	fl.Usage = func() {
		fmt.Fprintf(os.Stderr,
			`htpasswd sets/deletes username/password entries from an htpasswd file.
The only supported algorithm is bcrypt.

Usage:
        htpasswd [-file htpasswdfile] set   -user username -pass password
        htpasswd [-file htpasswdfile] set   -user username -pass-stdin
        htpasswd [-file htpasswdfile] check -user username -pass password
        htpasswd [-file htpasswdfile] check -user username -pass-stdin
        htpasswd [-file htpasswdfile] del   -user username
`)
	}

	var filename string
	fl.StringVar(&filename, "file", "htpasswd", "htpasswd file to update")
	fl.Parse(os.Args[1:])
	action := fl.Arg(0)

	f2 := flag.NewFlagSet(fl.Arg(0), flag.ExitOnError)
	f2.Usage = fl.Usage
	var username, pass string
	var passStdin bool
	f2.StringVar(&username, "user", "", "username")
	switch action {
	case "set", "check":
		f2.StringVar(&pass, "pass", "", "password")
		f2.BoolVar(&passStdin, "pass-stdin", false, "read password from stdin")
	case "del":
	default:
		fl.Usage()
		os.Exit(1)
	}
	f2.Parse(fl.Args()[1:])
	if len(f2.Args()) > 0 {
		f2.Usage()
		os.Exit(1)
	}

	var newEntry []byte
	switch action {
	case "set", "check":
		if passStdin {
			sc := bufio.NewScanner(os.Stdin)
			if !sc.Scan() {
				checkErr("reading passwd", sc.Err())
			}
			pass = sc.Text()
		}
		hashed, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
		checkErr("hashing passwd", err)
		newEntry = append([]byte(username+":"), hashed...)
	}

	var out bytes.Buffer
	old, err := os.ReadFile(filename)
	if !errors.Is(err, fs.ErrNotExist) { // ignore not exist errors
		checkErr("reading file", err)
	}

	sc := bufio.NewScanner(bytes.NewReader(old))
	var found bool
	for sc.Scan() {
		b := sc.Bytes()
		if pre := append([]byte(username), ':'); bytes.HasPrefix(b, pre) {
			found = true
			switch action {
			case "add":
				b = newEntry
			case "check":
				err = bcrypt.CompareHashAndPassword(bytes.TrimPrefix(b, pre), []byte(pass))
				checkErr("mismatched passwd", err)
			case "del":
				continue
			}
		}
		out.Write(b)
		out.WriteRune('\n')
	}
	if !found {
		out.Write(newEntry)
		out.WriteRune('\n')
	}

	err = os.WriteFile(filename, out.Bytes(), 0o644)
	checkErr("writing file", err)
}

func checkErr(msg string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, msg+": %v\n", err)
		os.Exit(1)
	}
}
