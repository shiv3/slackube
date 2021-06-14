package main

import (
	"github.com/shiv3/slackube/app/config"
	"github.com/shiv3/slackube/app/controller/api"
)

func main() {
	conf := config.Config{
		EnvConfig: config.EnvConfig{
			Env:         "",
			ServiceName: "",
			ProjectID:   "",
			LogLevel:    "",
		},
		ServerConfig: config.ServerConfig{
			ServerPort:      8080,
			MetricsPort:     0,
			ServerKeepAlive: 0,
			GracefulPeriod:  0,
			RequestTimeout:  0,
		},
	}
	api, err := api.NewServerImpl(&conf)
	if err != nil {
		panic(err)
	}
	api.Run()
}
