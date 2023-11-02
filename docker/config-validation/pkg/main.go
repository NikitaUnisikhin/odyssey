package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const filePath = "/etc/odyssey/odyssey-test.conf"
const newFilePath = "/etc/odyssey/odyssey-new-test.conf"

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

func currentWorkers() error {
	ctx := context.TODO()

	var err error
	out, err := exec.CommandContext(ctx, "/usr/bin/odyssey /etc/odyssey/odyssey-test.conf --test").Output()

	if err != nil {
		return err
	}

	if strOut := string(out); !strings.Contains(strOut, "config is valid") {
		return err
	}

	return nil
}

func noCurrentWorkers() error {
	ctx := context.TODO()

	err := changeConfig("workers", "workers -1")

	if err != nil {
		return err
	}

	out, err := exec.CommandContext(ctx, "/usr/bin/odyssey /etc/odyssey/odyssey-new-test.conf --test").Output()

	if err != nil {
		return err
	}

	if strOut := string(out); !strings.Contains(strOut, "bad workers number") {
		return err
	}

	return nil
}

func main() {
	if err := currentWorkers(); err == nil {
		fmt.Println("error: current workers field is not pass")
	} else if err := noCurrentWorkers(); err == nil {
		fmt.Println("error: no current workers field is pass")
	} else {
		fmt.Println("workers: Ok")
	}
}
