package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

const pathPrefix = "/etc/odyssey/configs"

const configIsValid = "config is valid"
const configWithInvalidValuePass = "config with invalid value pass"

func makeTest(pathToConfig string, isValidConfig bool) error {
	ctx := context.TODO()

	out, _ := exec.CommandContext(ctx, "/usr/bin/odyssey", pathToConfig, "--test").Output()
	strOut := string(out)

	if isValidConfig && !strings.Contains(strOut, configIsValid) {
		return errors.New(strOut)
	}

	if !isValidConfig && strings.Contains(strOut, configIsValid) {
		return errors.New(configWithInvalidValuePass)
	}

	return nil
}

func makeTests(field string) error {
	pathToDir := pathPrefix + "/" + field + "/valid"
	configs, _ := ioutil.ReadDir(pathToDir)

	for _, config := range configs {
		pathToConfig := pathToDir + "/" + config.Name()
		if err := makeTest(pathToConfig, true); err != nil {
			return err
		}
	}

	pathToDir = pathPrefix + "/" + field + "/invalid"
	configs, _ = ioutil.ReadDir(pathToDir)

	for _, config := range configs {
		pathToConfig := pathToDir + "/" + config.Name()
		if err := makeTest(pathToConfig, false); err != nil {
			return err
		}
	}

	return nil
}

func printTestsResult(field string) {
	if err := makeTests(field); err != nil {
		fmt.Println(field+"_test (ERROR) :", err)
	} else {
		fmt.Println(field + "_test: Ok")
	}
}

func runTests() {
	tests := []string{
		"workers",
		"resolvers",
		"coroutine_stack_size",
		"log_format",
		"unix_socket_mode",
		"listen_empty",
		"listen_tls",
		"storage_type",
		"storage_tls",
		"storage_name",
		"pool_type",
		"pool_reserve_prepared_statement",
		"pool_routing",
		"authentication",
		"auth_query",
	}

	for _, test := range tests {
		printTestsResult(test)
	}
}

func main() {
	runTests()
}
