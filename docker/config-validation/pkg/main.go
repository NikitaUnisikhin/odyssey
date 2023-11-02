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

func workers() error {
	var tests = []testCase{
		{"1", "config is valid", "current workers field is not pass"},
		{"10", "config is valid", "current workers field is not pass"},
		{"auto", "config is valid", "current workers field is not pass"},
		{"-1", "bad workers number", "no current workers field is pass"},
		{"-10", "bad workers number", "no current workers field is pass"},
	}

	ctx := context.TODO()

	for _, test := range tests {
		err := changeConfig("workers", "workers "+test.input)

		if err != nil {
			return err
		}

		out, err := exec.CommandContext(ctx, "/usr/bin/odyssey", "/etc/odyssey/odyssey-new-test.conf", "--test").Output()

		if strOut := string(out); !strings.Contains(strOut, test.outputPrefix) {
			return errors.New(test.errorMsg)
		}
	}

	return nil
}

func main() {
	if err := workers(); err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("workers: Ok")
	}
}
