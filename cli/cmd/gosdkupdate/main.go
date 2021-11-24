package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func main() {
	var minMinor int
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `%s [FLAGS]
        Updates the patch version of multiple Go minor versions.
`, os.Args[0])
		flag.CommandLine.PrintDefaults()
	}
	flag.IntVar(&minMinor, "min", 11, "minimum minor version to keep")
	flag.Parse()

	latestPatches, err := getLatest(minMinor)
	if err != nil {
		fmt.Fprintf(os.Stderr, "get latest versions: %v", err)
		os.Exit(1)
	}

	dir, err := cleanup()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cleanup: %v", err)
		os.Exit(1)
	}

	latestPatches = append(latestPatches, "gotip")

	var wg sync.WaitGroup
	for _, v := range latestPatches {
		wg.Add(1)
		go func(v string) {
			defer wg.Done()
			b, err := exec.Command("go", "install", "golang.org/dl/"+v+"@latest").CombinedOutput()
			if err != nil {
				fmt.Fprintf(os.Stderr, "get %s: %v: %s", v, err, b)
				return
			}
			b, err = exec.Command(v, "download").CombinedOutput()
			if err != nil {
				fmt.Fprintf(os.Stderr, "download %s: %v: %s", v, err, b)
				return
			}
			fmt.Println("updated", v)
		}(v)
	}

	wg.Wait()

	os.Link(filepath.Join(dir, "gotip"), filepath.Join(dir, "go"))
}

func cleanup() (string, error) {
	dir := os.Getenv("GOBIN")
	if dir == "" {
		dir = filepath.Join(os.Getenv("GOPATH"), "bin")
	}
	des, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("read dir %s: %w", dir, err)
	}
	os.Remove(filepath.Join(dir, "go"))
	for _, de := range des {
		if strings.HasPrefix(de.Name(), "go1.") {
			f := filepath.Join(dir, de.Name())
			err = os.Remove(f)
			if err != nil {
				return "", fmt.Errorf("rm %s: %w", f, err)
			}
		}
	}
	return dir, nil
}

func getLatest(minMinor int) ([]string, error) {
	u := "https://go.dev/dl/?mode=json&include=all"
	res, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf("GET %s: %w", u, err)
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code %s", res.Status)
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	var releases []Release
	err = json.Unmarshal(b, &releases)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	maxPatch := make(map[int]int, 30)  // minor: patch
	maxPre := make(map[int]string, 30) // minor: prerel
	for _, rel := range releases {
		var minor, patch int
		var pre string

		v := strings.TrimPrefix(rel.Version, "go")
		parts := strings.Split(v, ".")
		switch len(parts) {
		case 1:
			// go1
		case 2:
			// go1.1 go1.1beta1
			n, _ := fmt.Sscanf(parts[1], "%d", &minor)
			pre = parts[1][n:]
		case 3:
			// go1.1.1
			minor, _ = strconv.Atoi(parts[1])
			patch, _ = strconv.Atoi(parts[2])
		}
		if pre != "" {
			c := maxPre[minor]
			if c < pre {
				maxPre[minor] = pre
			}
		} else {
			c := maxPatch[minor]
			if c < patch {
				maxPatch[minor] = patch
			}
		}
	}

	var keep []string
	for minor, patch := range maxPatch {
		if minor < minMinor {
			continue
		}
		var v string
		if patch == 0 {
			v = fmt.Sprintf("go1.%d", minor)
		} else {
			v = fmt.Sprintf("go1.%d.%d", minor, patch)
		}
		keep = append(keep, v)
	}
	for minor, pre := range maxPre {
		if _, ok := maxPatch[minor]; minor < minMinor || ok {
			continue
		}
		keep = append(keep, fmt.Sprintf("go1.%d%s", minor, pre))
	}
	sort.Strings(keep)
	return keep, nil
}

type Release struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
	Files   []struct {
		Filename string `json:"filename"`
		Os       string `json:"os"`
		Arch     string `json:"arch"`
		Version  string `json:"version"`
		Sha256   string `json:"sha256"`
		Size     int    `json:"size"`
		Kind     string `json:"kind"`
	} `json:"files"`
}
