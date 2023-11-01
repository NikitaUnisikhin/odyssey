package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

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

func main() {
	if err := currentWorkers(); err != nil {
		fmt.Println("error: current workers field is not pass")
	} else {
		log.Println("workers: Ok")
	}
}
