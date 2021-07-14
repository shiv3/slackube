package config

import "time"

type (
	Config struct {
		EnvConfig    EnvConfig    `mapstructure:"env" validate:"required"`
		ServerConfig ServerConfig `mapstructure:"api" validate:"required"`
		SlackConfig  SlackConfig  `mapstructure:"slack" validate:"required"`
	}

	EnvConfig struct {
		Env         string `validate:"required"`
		ServiceName string `mapstructure:"service_name" validate:"required"`
		ProjectID   string `mapstructure:"project_id" validate:"required"`
		LogLevel    string `mapstructure:"log_level" validate:"required"`
	}

	ServerConfig struct {
		PingPath        string `mapstructure:"ping_path" default:"/eventsapi"`
		SlackEventsPath string `mapstructure:"slack_events_path" default:"/eventsapi"`
		SlacActionsPath string `mapstructure:"slack_actions_path" default:"/interactive-endpoint"`

		ServerPort      int           `mapstructure:"server_port" validate:"required"`
		MetricsPort     uint32        `mapstructure:"metrics_port" validate:"required"`
		ServerKeepAlive time.Duration `mapstructure:"server_keep_alive_ms" validate:"required"`
		GracefulPeriod  time.Duration `mapstructure:"graceful_period" validate:"required"`
		RequestTimeout  time.Duration `mapstructure:"request_timeout" validate:"required"`
	}

	SlackConfig struct {
		Token         string `mapstructure:"token" validate:"required"`
		SigningSecret string `mapstructure:"signing_secret" validate:"required"`
	}
)
