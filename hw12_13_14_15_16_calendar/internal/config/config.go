package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

func Read(path string) (*Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	// Подставим переменные окружения в виде ${VAR}
	expanded := os.ExpandEnv(string(raw))
	// fmt.Println(expanded)
	var cfg Config
	if _, err := toml.Decode(expanded, &cfg); err != nil {
		return nil, fmt.Errorf("decode toml: %w", err)
	}

	return &cfg, nil
}

type Config struct {
	PSQL    PSQLConfig
	Logger  LoggerConfig
	HTTP    HTTPConfig
	GRPC    GRPCConfig
	Storage StorageConfig // "memory" или "sql"

	Queue     QueueConfig
	Scheduler SchedulerConfig
	Sender    SenderConfig
}

type QueueConfig struct {
	URL  string
	Name string
}

type SchedulerConfig struct {
	IntervalSeconds      int // how often to run (in seconds)
	CleanupEnabled       bool
	CleanupOlderThanDays int
}

type SenderConfig struct {
	LogLevel string
}

type StorageConfig struct {
	Type string
}

type PSQLConfig struct {
	DSN       string
	Migration string
}

type LoggerConfig struct {
	Level string
	Path  string
}

type HTTPConfig struct {
	Host string
	Port string
}

type GRPCConfig struct {
	Host string
	Port string
}

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
// type Config struct {
// 	Logger LoggerConf
// 	// TODO
// }

// type LoggerConf struct {
// 	Level string
// 	// TODO
// }

// func NewConfig() Config {
// 	return Config{}
// }

// TODO
