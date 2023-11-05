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

func makeTests(field string) {
	pathToDir := pathPrefix + "/" + field + "/valid"
	configs, _ := ioutil.ReadDir(pathToDir)

	for ind, config := range configs {
		pathToConfig := pathToDir + "/" + config.Name()
		if err := makeTest(pathToConfig, true); err != nil {
			fmt.Println(field+"_test_valid_"+string(ind)+"(ERROR) :", err)
		} else {
			fmt.Println(field + "_test_valid_" + string(ind) + ": Ok")
		}
	}

	pathToDir = pathPrefix + "/" + field + "/invalid"
	configs, _ = ioutil.ReadDir(pathToDir)

	for ind, config := range configs {
		pathToConfig := pathToDir + "/" + config.Name()
		if err := makeTest(pathToConfig, false); err != nil {
			fmt.Println(field+"_test_invalid_"+string(ind)+"(ERROR) :", err)
		} else {
			fmt.Println(field + "_test_invalid_" + string(ind) + ": Ok")
		}
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
		makeTests(test)
	}
}

func main() {
	runTests()
}
