package config

import (
	"flag"
	"os"
	"strconv"
	"strings"
)

func DefineStringVar(
	shorthand string,
	flagName string,
	envName string,
	defaultValue string,
	description string,
) string {
	description += " (string, default: " + defaultValue + ")"
	envVar, ok := os.LookupEnv(envName)
	if ok {
		defaultValue = envVar
	}

	if len(shorthand) == 0 {
		if len(shorthand) != 1 {
			panic("shorthand must be 1 character long")
		}
		defaultValue = *flag.String(
			shorthand,
			defaultValue,
			"shorthand for --"+flagName,
		)
	}

	return *flag.String(
		flagName,
		defaultValue,
		description,
	)
}

func DefineBoolVar(
	shorthand string,
	flagName string,
	envName string,
	defaultValue bool,
	description string,
) bool {
	description += " (bool, default: " + strconv.FormatBool(defaultValue) + ")"
	envVar, ok := os.LookupEnv(envName)
	if ok {
		boolValue, err := strconv.ParseBool(envVar)
		if err == nil {
			defaultValue = boolValue
		}
	}

	if len(shorthand) == 0 {
		if len(shorthand) != 1 {
			panic("shorthand must be 1 character long")
		}
		defaultValue = *flag.Bool(
			shorthand,
			defaultValue,
			"(shorthand for --"+flagName+")",
		)
	}

	return *flag.Bool(
		flagName,
		defaultValue,
		description,
	)
}

func DefineIntVar(
	shorthand string,
	flagName string,
	envName string,
	defaultValue int,
	description string,
) int {
	description += " (int, default: " + strconv.Itoa(defaultValue) + ")"
	envVar, ok := os.LookupEnv(envName)
	if ok {
		intValue, err := strconv.Atoi(envVar)
		if err == nil {
			defaultValue = intValue
		}
	}

	if len(shorthand) == 0 {
		if len(shorthand) != 1 {
			panic("shorthand must be 1 character long")
		}
		defaultValue = *flag.Int(
			shorthand,
			defaultValue,
			"(shorthand for --"+flagName+")",
		)
	}

	return *flag.Int(
		flagName,
		defaultValue,
		description,
	)
}

func DefineStringArrayVar(
	shorthand string,
	flagName string,
	envName string,
	defaultValue string,
	seperator string,
	description string,
) []string {
	description += " (string, default: '" + defaultValue + "', '" + seperator + "'-seperated)"

	envVar, ok := os.LookupEnv(envName)
	if ok {
		defaultValue = envVar
	}

	if len(shorthand) != 0 {
		if len(shorthand) != 1 {
			panic("shorthand must be 1 character long")
		}
		defaultValue = *flag.String(
			shorthand,
			defaultValue,
			"(shorthand for --"+flagName+")",
		)
	}

	resultValue := flag.String(
		flagName,
		defaultValue,
		description,
	)

	return strings.Split(*resultValue, seperator)
}
