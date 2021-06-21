package config

import "time"

type (
	Config struct {
		EnvConfig    EnvConfig    `mapstructure:"env"`
		ServerConfig ServerConfig `mapstructure:"api"`
		SlackConfig  SlackConfig
	}

	EnvConfig struct {
		Env         string `validate:"required"`
		ServiceName string `mapstructure:"service_name" validate:"required"`
		ProjectID   string `mapstructure:"project_id" validate:"required"`
		LogLevel    string `mapstructure:"log_level" validate:"required"`
	}

	ServerConfig struct {
		ServerPort      int           `mapstructure:"server_port" validate:"required"`
		MetricsPort     uint32        `mapstructure:"metrics_port" validate:"required"`
		ServerKeepAlive time.Duration `mapstructure:"server_keep_alive_ms" validate:"required"`
		GracefulPeriod  time.Duration `mapstructure:"graceful_period" validate:"required"`
		RequestTimeout  time.Duration `mapstructure:"request_timeout" validate:"required"`
	}

	SlackConfig struct {
		Token         string
		SingingSecret string
	}
)
