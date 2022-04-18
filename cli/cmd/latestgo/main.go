package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type GoVersion struct {
	Version string `json:"version"`
}

func main() {
	res, err := http.Get("https://go.dev/dl/?mode=json")
	if err != nil {
		log.Fatalln("get available versions", err)
	}
	if res.StatusCode != 200 {
		log.Fatalln(res.Status)
	}
	defer res.Body.Close()
	var versions []GoVersion
	err = json.NewDecoder(res.Body).Decode(&versions)
	if err != nil {
		log.Fatalln("decode go version", err)
	}
	fmt.Print(versions[0].Version)
}
