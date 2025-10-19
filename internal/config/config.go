package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Env string

const (
	EnvLocal Env = "local"
	EnvDev   Env = "dev"
)

type Log struct {
	Level string `mapstructure:"level"` // debug|info|warn|error
}

type DynamoDB struct {
	Region           string `mapstructure:"region"`
	EndpointOverride string `mapstructure:"endpointOverride"` // http://localhost:8000 for local
	ShowsTable       string `mapstructure:"showsTable"`
}

type Cognito struct {
	UserPoolID  string   `mapstructure:"userPoolId"`
	ClientID    string   `mapstructure:"clientId"`
	Region      string   `mapstructure:"region"`
	JWKSURL     string   `mapstructure:"jwksUrl"`
	ValidScopes []string `mapstructure:"validScopes"`
}

type Config struct {
	Env      Env      `mapstructure:"env"`
	Log      Log      `mapstructure:"log"`
	DynamoDB DynamoDB `mapstructure:"dynamodb"`
	Cognito  Cognito  `mapstructure:"cognito"`
}

func Load() (*Config, error) {
	v := viper.New()

	// Set defaults
	v.SetDefault("env", string(EnvLocal))
	v.SetDefault("log.level", "info")
	v.SetDefault("dynamodb.region", "ap-southeast-2")
	v.SetDefault("dynamodb.endpointOverride", "")
	v.SetDefault("dynamodb.createTableIfMissing", false)

	env := determineEnvironment()

	v.SetDefault("dynamodb.showsTable", "shows-"+env)

	if env == string(EnvLocal) {
		v.SetConfigFile("configs/config.local.yaml")
		_ = v.ReadInConfig()
	}

	if env == string(EnvDev) {
		v.SetConfigFile("configs/config.dev.yaml")
		_ = v.ReadInConfig()
	}

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

// determineEnvironment detects if we're running in ECS/production
func determineEnvironment() string {

	if appEnv := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV"))); appEnv != "" {
		if appEnv == string(EnvDev) {
			return string(EnvDev)
		}
		if appEnv == string(EnvLocal) {
			return string(EnvLocal)
		}
	}

	if os.Getenv("ECS_CONTAINER_METADATA_URI") != "" ||
		os.Getenv("AWS_EXECUTION_ENV") != "" {
		return string(EnvDev)
	}

	return string(EnvLocal)
}
