package appconfig

import (
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// NewEnvironment returns an environment from a string
func NewEnvironment(config string) (Environment, error) {
	switch config {
	case localEnv.String():
		return localEnv, nil
	case testEnv.String():
		return testEnv, nil
	case devEnv.String():
		return devEnv, nil
	case implEnv.String():
		return implEnv, nil
	case prodEnv.String():
		return prodEnv, nil
	default:
		return "", fmt.Errorf("unknown environment: %s", config)
	}
}

// EnvironmentKey is used to access the environment from a config
const EnvironmentKey = "APP_ENV"

// Environment represents an environment
type Environment string

const (
	// localEnv is the local environment
	localEnv Environment = "local"
	// testEnv is the environment for running tests
	testEnv Environment = "test"
	// devEnv is the environment for the dev deployed env
	devEnv Environment = "dev"
	// implEnv is the environment for the impl deployed env
	implEnv Environment = "impl"
	// prodEnv is the environment for the impl deployed env
	prodEnv Environment = "prod"
)

// String gets the environment as a string
func (e Environment) String() string {
	switch e {
	case localEnv:
		return "local"
	case testEnv:
		return "test"
	case devEnv:
		return "dev"
	case implEnv:
		return "impl"
	case prodEnv:
		return "prod"
	default:
		return ""
	}
}

// Local returns true if the environment is local
func (e Environment) Local() bool {
	return e == localEnv
}

// Test returns true if the environment is local
func (e Environment) Test() bool {
	return e == testEnv
}

// Dev returns true if the environment is local
func (e Environment) Dev() bool {
	return e == devEnv
}

// Impl returns true if the environment is local
func (e Environment) Impl() bool {
	return e == implEnv
}

// Prod returns true if the environment is local
func (e Environment) Prod() bool {
	return e == prodEnv
}

// Deployed returns true if in a deployed environment
func (e Environment) Deployed() bool {
	switch e {
	case devEnv:
		return true
	case implEnv:
		return true
	case prodEnv:
		return true
	default:
		return false
	}
}

// APIProtocol is the key the API protocol
const APIProtocolKey = "API_PROTOCOL"

// APIHostKey is the key the API hostname
const APIHostKey = "API_HOST"

// APIPortKey is the key the API port
const APIPortKey = "API_PORT"

// DBHostConfigKey is the Postgres hostname config key
const DBHostConfigKey = "PGHOST"

// DBPortConfigKey is the Postgres port config key
const DBPortConfigKey = "PGPORT"

// DBNameConfigKey is the Postgres database name config key
const DBNameConfigKey = "PGDATABASE"

// DBUsernameConfigKey is the Postgres username config key
const DBUsernameConfigKey = "PGUSER"

// DBPasswordConfigKey is the Postgres password config key
const DBPasswordConfigKey = "PGPASS"

// DBSSLModeConfigKey is the Postgres SSL mode config key
const DBSSLModeConfigKey = "PGSSLMODE"

// AppConfig is open struct so tests can create as necessary
type AppConfig struct {
	Env Environment

	Logger *zap.Logger

	APIProtocol string
	APIHost string
	APIPort string

	DBHost string
	DBPort string
	DBName string
	DBUsername string
	DBPassword string
	DBSSLMode string

	AppVersion string
	AppDatetime string
	AppTimestamp string
}

const configMissingMessage = "Must set config: %v"

func getRequiredConfig(viperConfig *viper.Viper, configKey string, logger *zap.Logger) string {
	val := viperConfig.GetString(configKey)
	if val == "" {
		logger.Fatal(fmt.Sprintf(configMissingMessage, configKey))
	}
	return val
}

func getDefaultConfigString(viperConfig *viper.Viper, configKey string, defaultVal string) string {
	val := viperConfig.GetString(configKey)
	if val == "" {
		return defaultVal
	}
	return val
}

func NewLogger(env Environment) (*zap.Logger, error) {
	var zapLogger *zap.Logger
	var err error
	if env.Prod() {
		zapLogger, err = zap.NewProduction()
	} else {
		zapLogger, err = zap.NewDevelopment()
	}
	if err != nil {
		return nil, err
	}
	return zapLogger, nil
}

func NewAppConfig(viperConfig *viper.Viper) (AppConfig, error) {
	// Set environment from config, shared for all
	environment, err := NewEnvironment(viperConfig.GetString(EnvironmentKey))
	if err != nil {
		return AppConfig{}, err
	}

	// Build logger once, shared for all
	zapLogger, err := NewLogger(environment)
	if err != nil {
		return AppConfig{}, err
	}

	appConfig := AppConfig{
		Env: environment,

		Logger: zapLogger,

		APIProtocol: viperConfig.GetString(APIProtocolKey),
		APIHost: getRequiredConfig(viperConfig, APIHostKey, zapLogger),
		APIPort: viperConfig.GetString(APIPortKey),

		DBHost: getRequiredConfig(viperConfig, DBHostConfigKey, zapLogger),
		DBPort: getRequiredConfig(viperConfig, DBPortConfigKey, zapLogger),
		DBName: getRequiredConfig(viperConfig, DBNameConfigKey, zapLogger),
		DBUsername: getRequiredConfig(viperConfig, DBUsernameConfigKey, zapLogger),
		DBPassword: getRequiredConfig(viperConfig, DBPasswordConfigKey, zapLogger),
		DBSSLMode: getRequiredConfig(viperConfig, DBSSLModeConfigKey, zapLogger),

		AppVersion: viperConfig.GetString("APPLICATION_VERSION"),
		AppDatetime: viperConfig.GetString("APPLICATION_DATETIME"),
		AppTimestamp: viperConfig.GetString("APPLICATION_TS"),
	}

	return appConfig, nil
}
