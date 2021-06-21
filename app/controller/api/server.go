package api

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/shiv3/slackube/app/config"
	v1 "github.com/shiv3/slackube/app/controller/api/v1"
	"github.com/shiv3/slackube/app/controller/api/v1/ping"
	"github.com/shiv3/slackube/app/controller/api/v1/slackevents"
)

type ServerImpl struct {
	Config     *config.Config
	EchoServer *echo.Echo
	V1Router   v1.RouterImpl
}

func NewServerImpl(config *config.Config) (ServerImpl, error) {
	v1Router := v1.NewV1Router(
		slackevents.NewHandlerImpl(""),
		ping.NewHandlerImpl("ok"),
	)
	return ServerImpl{
		Config:     config,
		EchoServer: echo.New(),
		V1Router:   v1Router,
	}, nil
}

func (s ServerImpl) Run() error {
	s.V1Router.Dispatch(s.EchoServer)
	go func() {
		if err := s.EchoServer.Start(fmt.Sprintf(":%d", s.Config.ServerConfig.ServerPort)); err != nil {
			// logger.Fatalf(ctx, "failed to start server. err -> [%s]", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal)
	// syscall.SIGINT: Ctrl-C
	// syscall.SIGTERM: from docker, k8s...
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		// TODO: move config
		timeout := s.Config.ServerConfig.GracefulPeriod
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := s.EchoServer.Shutdown(ctx); err != nil {
			// logger.Errorf(ctx, "failed to shutdown. err -> [%s]", err)
		}
	}
	return nil
}
