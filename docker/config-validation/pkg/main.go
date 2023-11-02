package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
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

func changeConfig(prefix string, stringToReplace string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(scanner.Text(), prefix) {
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

func checkWorkers() error {
	const currentFieldIsNotPass = "current workers field is not pass"
	const noCurrentFieldIsPass = "no current workers field is pass"
	const badWorkersNumber = "bad workers number"

	var tests = []testCase{
		{"1", configIsValid, currentFieldIsNotPass},
		{"10", configIsValid, currentFieldIsNotPass},
		{"\"auto\"", configIsValid, currentFieldIsNotPass},
		{"-1", badWorkersNumber, noCurrentFieldIsPass},
		{"0", badWorkersNumber, noCurrentFieldIsPass},
		{"-10", badWorkersNumber, noCurrentFieldIsPass},
	}

	ctx := context.TODO()

	for _, test := range tests {
		err := changeConfig("workers", "workers "+test.input)

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

func checkResolvers() error {
	const currentFieldIsNotPass = "current resolvers field is not pass"
	const noCurrentFieldIsPass = "no current resolvers field is pass"
	const badResolversNumber = "bad resolvers number"

	var tests = []testCase{
		{"1", configIsValid, currentFieldIsNotPass},
		{"10", configIsValid, currentFieldIsNotPass},
		{"-1", badResolversNumber, noCurrentFieldIsPass},
		{"0", badResolversNumber, noCurrentFieldIsPass},
		{"-10", badResolversNumber, noCurrentFieldIsPass},
	}

	ctx := context.TODO()

	for _, test := range tests {
		err := changeConfig("resolvers", "resolvers "+test.input)

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

func main() {
	if err := checkWorkers(); err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("checkWorkers: Ok")
	}

	if err := checkResolvers(); err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("checkWorkers: Ok")
	}
}
