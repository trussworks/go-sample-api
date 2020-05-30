package appconfig

import "fmt"

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
