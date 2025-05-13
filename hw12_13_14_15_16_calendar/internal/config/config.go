package config

import "github.com/BurntSushi/toml"

func Read(fpath string) (c Config, err error) {
	_, err = toml.DecodeFile(fpath, &c)
	return
}

type Config struct {
	PSQL   PSQLConfig
	Logger LoggerConfig
	HTTP   HTTPConfig
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
