package main

import (
	"github.com/shiv3/configmapper"
	"github.com/shiv3/slackube/app/config"
	"github.com/shiv3/slackube/app/controller/api"
)

func main() {
	c, err := configmapper.Initialize(config.Config{})
	if err != nil {
		panic(err)
	}
	conf := c.(config.Config)
	api, err := api.NewServerImpl(&conf)
	if err != nil {
		panic(err)
	}
	api.Run()
}
