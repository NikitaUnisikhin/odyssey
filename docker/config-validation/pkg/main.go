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

func makeTest(pathToConfig string, prefix string) error {
	ctx := context.TODO()

	out, _ := exec.CommandContext(ctx, "/usr/bin/odyssey", pathToConfig, "--test").Output()

	if strOut := string(out); !strings.Contains(strOut, prefix) {
		return errors.New(strOut)
	}

	return nil
}

func makeTests(field string, errorTriggerMsg string) error {
	pathToDir := pathPrefix + "/" + field + "/valid"
	configs, _ := ioutil.ReadDir(pathToDir)

	for _, config := range configs {
		pathToConfig := pathToDir + "/" + config.Name()
		if err := makeTest(pathToConfig, configIsValid); err != nil {
			return err
		}
	}

	pathToDir = pathPrefix + "/" + field + "/invalid"
	configs, _ = ioutil.ReadDir(pathToDir)

	for _, config := range configs {
		pathToConfig := pathToDir + "/" + config.Name()
		if err := makeTest(pathToConfig, errorTriggerMsg); err == nil {
			return err
		}
	}

	return nil
}

func main() {
	if err := makeTests("workers", "bad workers number"); err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("workersTest: Ok")
	}
}
