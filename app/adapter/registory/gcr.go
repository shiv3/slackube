package main

import (
	"fmt"

	"github.com/heroku/docker-registry-client/registry"
)

func main() {
	url := "https://asia.gcr.io/"
	username := "" // anonymous
	password := "" // anonymous
	r, err := registry.New(url, username, password)
	if err != nil {
		panic(err)
	}
	fmt.Println(r)

	tags, err := r.Tags("mira-280805/mira-bid")
	for _, tag := range tags {
		fmt.Sprintln(tag)
	}
}
