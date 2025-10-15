package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Env string

const (
	EnvLocal Env = "local"
	EnvProd  Env = "prod"
)

type Log struct {
	Level string `mapstructure:"level"` // debug|info|warn|error
}

type DynamoDB struct {
	Region               string `mapstructure:"region"`
	EndpointOverride     string `mapstructure:"endpointOverride"` // http://localhost:8000 for local
	ShowsTable           string `mapstructure:"showsTable"`
	CreateTableIfMissing bool   `mapstructure:"createTableIfMissing"`
}

type Config struct {
	Env      Env      `mapstructure:"env"`
	Log      Log      `mapstructure:"log"`
	DynamoDB DynamoDB `mapstructure:"dynamodb"`
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetDefault("env", string(EnvLocal))
	v.SetDefault("log.level", "info")
	v.SetDefault("dynamodb.region", "ap-southeast-2")
	v.SetDefault("dynamodb.endpointOverride", "")
	v.SetDefault("dynamodb.showsTable", "shows")
	v.SetDefault("dynamodb.createTableIfMissing", false)

	env := strings.ToLower(strings.TrimSpace(getEnv("APP_ENV", string(EnvLocal))))
	if env != string(EnvLocal) && env != string(EnvProd) {
		env = string(EnvLocal)
	}

	v.SetConfigFile(fmt.Sprintf("configs/config.%s.yaml", env))
	_ = v.ReadInConfig() // optional

	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	cfg.Env = Env(env)
	return &cfg, nil
}

func getEnv(key, def string) string {
	tmp := viper.New()
	tmp.AutomaticEnv()
	if s := tmp.GetString(key); s != "" {
		return s
	}
	return def
}
