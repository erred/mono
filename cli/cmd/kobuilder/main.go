// kobuilder builds the container image used in the cloudbuild ci process.
// The only dependencies are that authentication is available by default.
package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

var logV bool

func log(a ...interface{}) {
	if logV {
		fmt.Fprintln(os.Stderr, a...)
	}
}

func main() {
	var baseImgName, targetImgName string
	flag.StringVar(&baseImgName, "base", "gcr.io/google.com/cloudsdktool/cloud-sdk:alpine", "base image name")
	flag.StringVar(&targetImgName, "target", "us-central1-docker.pkg.dev/com-seankhliao/run/ko:latest", "target image name")
	flag.BoolVar(&logV, "v", false, "log progress")
	flag.Parse()

	err := run(baseImgName, targetImgName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run(baseImgName, targetImgName string) error {
	koURL, err := latestKo()
	if err != nil {
		return fmt.Errorf("get latest ko: %w", err)
	}
	log("got ko url", koURL)
	fmt.Println("Ko:", koURL)

	goURL, err := latestGo()
	if err != nil {
		return fmt.Errorf("get latest go: %w", err)
	}
	log("got go url", goURL)
	fmt.Println("Go:", goURL)

	baseRef, err := name.ParseReference(baseImgName)
	if err != nil {
		return fmt.Errorf("parse base img=%s: %w", baseImgName, err)
	}
	log("got base ref", baseRef)

	img, err := remote.Image(baseRef)
	if err != nil {
		return fmt.Errorf("get base: %w", err)
	}
	log("got base img")

	koLayer, err := tarToLayer(koURL, "/usr/local/bin", "ko")
	if err != nil {
		return fmt.Errorf("get ko layer: %w", err)
	}
	log("got ko layer")

	goLayer, err := tarToLayer(goURL, "/usr/local", "")
	if err != nil {
		return fmt.Errorf("get go layer: %w", err)
	}
	log("got go layer")

	img, err = mutate.AppendLayers(img, goLayer, koLayer)
	if err != nil {
		return fmt.Errorf("append layers: %w", err)
	}
	log("appended layers")

	configFile, err := img.ConfigFile()
	if err != nil {
		return fmt.Errorf("get config file")
	}
	log("got config file")
	configFile.Config.WorkingDir = "/worspace"
	for i, e := range configFile.Config.Env {
		if strings.HasPrefix(e, "PATH=") {
			configFile.Config.Env[i] = "PATH=/usr/local/go/bin:" + e[5:]
		}
	}
	configFile.Config.Env = append(
		configFile.Config.Env,
		"CGO_ENABLED=0",
		"GOFLAGS=-trimpath",
	)

	img, err = mutate.ConfigFile(img, configFile)
	if err != nil {
		return fmt.Errorf("set config file: %w", err)
	}
	log("mutated config file")

	targetRef, err := name.ParseReference(targetImgName)
	if err != nil {
		return fmt.Errorf("parse target img=%s: %w", targetImgName, err)
	}
	log("got target ref", targetRef)

	err = remote.Write(targetRef, img, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return fmt.Errorf("write target img=%s: %w", targetImgName, err)
	}
	log("wrote target img")
	dig, err := img.Digest()
	if err != nil {
		return fmt.Errorf("get img digest: %w", err)
	}
	fmt.Printf("%s@%s\n", targetRef, dig)

	return nil
}

// turns a remote tarball to an image layer
// download is the download url.
// root is the target root within the image for the extracted tarball.
// only limits the output to a single file matching the name.
func tarToLayer(download, root, only string) (v1.Layer, error) {
	res, err := http.Get(download)
	if err != nil {
		return nil, fmt.Errorf("GET %s: %w", download, err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("GET %s: %s", download, res.Status)
	}
	log("got tar response")

	gr, err := gzip.NewReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("gzip reader %s: %w", download, err)
	}

	tr := tar.NewReader(gr)
	pr, pw := io.Pipe()
	tw := tar.NewWriter(pw)

	go func() {
		log("started writing tar")
		defer log("finished writing tar")
		defer pw.Close()
		defer tw.Close()
		for th, err := tr.Next(); err == nil; th, err = tr.Next() {
			if only != "" && th.Name != only {
				continue
			}
			th.Name = path.Join(root, th.Name)
			err = tw.WriteHeader(th)
			if err != nil {
				panic("write header: " + err.Error())
			}
			_, err = io.Copy(tw, tr)
			if err != nil {
				panic("copy body: " + err.Error())
			}
		}
	}()

	return tarball.LayerFromReader(pr, tarball.WithEstargz)
}

func latestKo() (string, error) {
	res, err := http.Get("https://api.github.com/repos/google/ko/releases")
	if err != nil {
		return "", fmt.Errorf("GET google/ko: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", fmt.Errorf("GET google/ko: %s", res.Status)
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("read google/ko: %w", err)
	}
	var ghrel ghRel
	err = json.Unmarshal(b, &ghrel)
	if err != nil {
		return "", fmt.Errorf("unmarshal google/ko: %w", err)
	}
	for _, asset := range ghrel[0].Assets {
		if strings.Contains(asset.Name, "Linux") && strings.Contains(asset.Name, "x86_64") {
			return asset.Browser_download_url, nil
		}
	}
	return "", fmt.Errorf("no file found in google/ko")
}

type ghRel []struct {
	Tag_name string
	Assets   []struct {
		Name                 string
		Browser_download_url string
	}
}

func latestGo() (string, error) {
	res, err := http.Get("https://go.dev/dl/?mode=json&include=all")
	if err != nil {
		return "", fmt.Errorf("GET go.dev/dl: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", fmt.Errorf("GET go.dev/dl: %s", res.Status)
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("read go.dev/dl: %w", err)
	}
	var godl goDL
	err = json.Unmarshal(b, &godl)
	if err != nil {
		return "", fmt.Errorf("unamarshal go.dev/dl: %w", err)
	}
	for _, file := range godl[0].Files {
		if file.Os == "linux" && file.Arch == "amd64" && file.Kind == "archive" {
			return "https://go.dev/dl/" + file.Filename, nil
		}
	}
	return "", fmt.Errorf("no file found in go.dev/dl")
}

type goDL []struct {
	Version string
	Stable  bool
	Files   []struct {
		Filename string
		Os       string
		Arch     string
		Version  string
		Sha256   string
		Size     int
		Kind     string
	}
}
