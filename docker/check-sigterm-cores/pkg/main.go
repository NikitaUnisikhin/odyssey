package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os/exec"
	"syscall"
	"time"
)

const benchTimeSec = 10
const odProcName = "odyssey"

func bunchProcess(ctx context.Context) {
	_, err := exec.CommandContext(ctx, "pgbench", "--builtin select-only",
		"-c 40", fmt.Sprintf("-T %d", benchTimeSec), "-j 20", "-n", "-h localhost", "-p 6432", "-U postgres", "db1", "-P 1").Output()

	files, _ := ioutil.ReadDir("/var/cores")
	fmt.Printf("COUNT CORES: %d", len(files))

	if err != nil {
		fmt.Printf(err.Error())
	}
}

func testProcess(ctx context.Context) {
	if err := ensureOdysseyRunning(ctx); err != nil {
		fmt.Println(err.Error())
	}

	go bunchProcess(ctx)

	rand.Seed(time.Now().UnixNano())
	timeSleepMs := rand.Float32() * benchTimeSec

	time.Sleep(time.Duration(timeSleepMs) * time.Millisecond)

	if _, err := signalToProc(syscall.SIGINT, odProcName); err != nil {
		fmt.Println(err.Error())
	}
}

func main() {
	fmt.Println("Start check-sigterm-cores")
	defer fmt.Println("End check-sigterm-cores")

	ctx := context.TODO()
	for i := 0; i < 1000; i++ {
		testProcess(ctx)
	}
}
