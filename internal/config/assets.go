package config

import (
	"os"
	"strconv"
)

func EnvIsString(envVar string, existCallback func(value string)) {
	value := os.Getenv(envVar)
	if len(value) == 0 {
		return
	}

	existCallback(value)
}

func EnvIsInt(envVar string, existCallback func(value int)) {
	value := os.Getenv(envVar)
	if len(value) == 0 {
		return
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return
	}

	existCallback(intValue)
}

func EnvIsBool(envVar string, existCallback func(value bool)) {
	value := os.Getenv(envVar)
	if len(value) == 0 {
		return
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return
	}

	existCallback(boolValue)
}
