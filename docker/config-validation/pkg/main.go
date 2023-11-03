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

func makeTest(pathToConfig string, errorTriggerMsg string, isValidConfig bool) error {
	ctx := context.TODO()

	out, _ := exec.CommandContext(ctx, "/usr/bin/odyssey", pathToConfig, "--test").Output()

	if strOut := string(out); strings.Contains(strOut, errorTriggerMsg) {
		if isValidConfig {
			return errors.New(errorTriggerMsg)
		} else {
			return errors.New(configWithInvalidValuePass)
		}
	}

	return nil
}

func makeTests(field string, errorTriggerMsg string) error {
	pathToDir := pathPrefix + "/" + field + "/valid"
	configs, _ := ioutil.ReadDir(pathToDir)

	for _, config := range configs {
		pathToConfig := pathToDir + "/" + config.Name()
		if err := makeTest(pathToConfig, errorTriggerMsg, true); err != nil {
			return err
		}
	}

	pathToDir = pathPrefix + "/" + field + "/invalid"
	configs, _ = ioutil.ReadDir(pathToDir)

	for _, config := range configs {
		pathToConfig := pathToDir + "/" + config.Name()
		if err := makeTest(pathToConfig, configIsValid, false); err != nil {
			return err
		}
	}

	return nil
}

func printTestsResult(field string, errorTriggerMsg string) {
	if err := makeTests(field, errorTriggerMsg); err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println(field + "_test: Ok")
	}
}

func main() {
	printTestsResult("workers", "bad workers number")
	printTestsResult("resolvers", "bad resolvers number")
	printTestsResult("coroutine_stack_size", "bad coroutine_stack_size number")
	printTestsResult("log_format", "log is not defined")
	printTestsResult("unix_socket_mode", "unix_socket_mode is not set")
	printTestsResult("listen_empty", "no listen servers defined")
	printTestsResult("unix_socket_dir", "listen host is not set and no unix_socket_dir is specified")
	printTestsResult("listen_tls", "unknown tls_opts->tls mode")
	printTestsResult("storage_type_null", "no type is specified")
	printTestsResult("storage_type", "unknown storage type")
	printTestsResult("storage_type_host_unix_socket_dir", "no host specified and unix_socket_dir is not set")
	printTestsResult("storage_tls", "unknown storage tls_opts->tls mode")
}
