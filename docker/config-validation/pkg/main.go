package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const filePath = "/etc/odyssey/odyssey-test.conf"
const newFilePath = "/etc/odyssey/odyssey-new-test.conf"

const configIsValid = "config is valid"

type testCase struct {
	input        string
	outputPrefix string
	errorMsg     string
}

func changeConfig(pattern string, stringToReplace string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		text := scanner.Text()
		if matched, _ := regexp.Match(pattern, []byte(text)); matched {
			text = stringToReplace
		}

		lines = append(lines, text)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	err = ioutil.WriteFile(newFilePath, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		return err
	}

	return nil
}

func check(tests []testCase, field string) error {
	ctx := context.TODO()

	for _, test := range tests {
		err := changeConfig(field+`*`, field+" "+test.input)

		if err != nil {
			return err
		}

		out, _ := exec.CommandContext(ctx, "/usr/bin/odyssey", "/etc/odyssey/odyssey-new-test.conf", "--test").Output()

		if strOut := string(out); !strings.Contains(strOut, test.outputPrefix) {
			return errors.New(test.errorMsg)
		}
	}

	return nil
}

func checkWorkers() error {
	const currentFieldIsNotPass = "current workers field is not pass"
	const noCurrentFieldIsPass = "no current workers field is pass"
	const badField = "bad workers number"

	var tests = []testCase{
		{"1", configIsValid, currentFieldIsNotPass},
		{"10", configIsValid, currentFieldIsNotPass},
		{"\"auto\"", configIsValid, currentFieldIsNotPass},
		{"-1", badField, noCurrentFieldIsPass},
		{"0", badField, noCurrentFieldIsPass},
		{"-10", badField, noCurrentFieldIsPass},
	}

	return check(tests, "workers")
}

func checkResolvers() error {
	const currentFieldIsNotPass = "current resolvers field is not pass"
	const noCurrentFieldIsPass = "no current resolvers field is pass"
	const badField = "bad resolvers number"

	var tests = []testCase{
		{"1", configIsValid, currentFieldIsNotPass},
		{"10", configIsValid, currentFieldIsNotPass},
		{"-1", badField, noCurrentFieldIsPass},
		{"0", badField, noCurrentFieldIsPass},
		{"-10", badField, noCurrentFieldIsPass},
	}

	return check(tests, "resolvers")
}

func checkCoroutineStackSize() error {
	const currentFieldIsNotPass = "current coroutine_stack_size field is not pass"
	const noCurrentFieldIsPass = "no current coroutine_stack_size field is pass"
	const badField = "bad coroutine_stack_size number"

	var tests = []testCase{
		{"4", configIsValid, currentFieldIsNotPass},
		{"10", configIsValid, currentFieldIsNotPass},
		{"-1", badField, noCurrentFieldIsPass},
		{"3", badField, noCurrentFieldIsPass},
		{"0", badField, noCurrentFieldIsPass},
		{"-10", badField, noCurrentFieldIsPass},
	}

	return check(tests, "coroutine_stack_size")
}

func checkLogFormat() error {
	const currentFieldIsNotPass = "current log_format field is not pass"
	const noCurrentFieldIsPass = "no current log_format field is pass"
	const badField = "log is not defined"

	var tests = []testCase{
		{`"%p %t %l [%i %s] (%c) %m\\n"`, configIsValid, currentFieldIsNotPass},
		{"", badField, noCurrentFieldIsPass},
	}

	return check(tests, "log_format")
}

func checkUnixSocketMode() error {
	const currentFieldIsNotPass = "current unix_socket_mode field is not pass"
	const noCurrentFieldIsPass = "no current unix_socket_mode field is pass"
	const badField = "unix_socket_mode is not set"

	var tests = []testCase{
		{`"0644"`, configIsValid, currentFieldIsNotPass},
		{"", badField, noCurrentFieldIsPass},
	}

	return check(tests, "unix_socket_mode")
}

func checkListen() error {
	const currentFieldIsNotPass = "config with listen is not pass"
	const noCurrentFieldIsPass = "config with missing listen is pass"
	const badField = "no listen servers defined"
	return nil
}

func main() {
	if err := checkWorkers(); err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("checkWorkers: Ok")
	}

	if err := checkResolvers(); err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("checkResolvers: Ok")
	}

	if err := checkCoroutineStackSize(); err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("checkCoroutineStackSize: Ok")
	}

	if err := checkLogFormat(); err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("checkLogFormat: Ok")
	}

	if err := checkUnixSocketMode(); err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("checkUnixSocketMode: Ok")
	}
}
