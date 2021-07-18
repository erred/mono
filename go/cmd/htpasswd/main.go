package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	var user, pass, file string
	flag.StringVar(&user, "user", "", "username")
	flag.StringVar(&pass, "pass", "", "password")
	flag.StringVar(&file, "file", "htpasswd", "htpasswd file")
	var add, del, stdin bool
	flag.BoolVar(&add, "add", false, "add a user/pass")
	flag.BoolVar(&del, "del", false, "delete a user")
	flag.BoolVar(&stdin, "pass-stdin", false, "read password form stdin")
	flag.Parse()

	if (add && del) || (!add && !del) {
		fmt.Fprintln(os.Stderr, "please specify one of -add / -del")
		os.Exit(1)
	}

	entries, err := readFile(file)
	if err != nil && !add {
		fmt.Fprintf(os.Stderr, "read %s: %v\n", file, err)
		os.Exit(1)
	}

	if add {
		if stdin {
			b, err := io.ReadAll(os.Stdin)
			if err != nil {
				fmt.Fprintf(os.Stderr, "read passwd from stdin: %v\n", err)
				os.Exit(1)
			}
			pass = string(b)
		}

		hashedPass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
		if err != nil {
			fmt.Fprintf(os.Stderr, "hash passed: %v\n", err)
			os.Exit(1)
		}
		entries = append(entries, []string{user, string(hashedPass)})
	} else { // del
		for i, entry := range entries {
			if entry[0] == user {
				entries = append(entries[:i], entries[i+1:]...)
				break
			}
		}
	}

	err = writeFile(file, entries)
	if err != nil {
		fmt.Fprintf(os.Stderr, "write file: %v\n", err)
		os.Exit(1)
	}
}

func writeFile(fn string, entries [][]string) error {
	f, err := os.Create(fn + ".tmp")
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()
	for i, entry := range entries {
		_, err = fmt.Fprintf(f, "%s:%s\n", entry[0], entry[1])
		if err != nil {
			return fmt.Errorf("write entru %d: %w", i, err)
		}
	}
	f.Close()
	err = os.Rename(fn+".tmp", fn)
	if err != nil {
		return fmt.Errorf("rename: %w", err)
	}
	return nil
}

func readFile(fn string) ([][]string, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	var entries [][]string
	users := make(map[string]struct{})
	sc := bufio.NewScanner(f)
	for i := 0; sc.Scan(); i++ {

		bb := bytes.SplitN(sc.Bytes(), []byte(":"), 2)
		if len(bb) != 2 {
			return nil, fmt.Errorf("invalid entry: %d", i)
		}
		u := bytes.TrimSpace(bb[0])
		if _, ok := users[string(u)]; ok {
			return nil, fmt.Errorf("duplicate entry for %s", string(u))
		}
		entries = append(entries, []string{string(u), string(bytes.TrimSpace(bb[1]))})
	}
	return entries, nil
}
