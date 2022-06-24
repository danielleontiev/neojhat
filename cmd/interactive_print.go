package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func interactive(progressFunc func() string, noop bool) func() {
	if noop {
		return func() {}
	}
	cancel := make(chan bool)
	go func() {
		for {
			select {
			case <-cancel:
				return
			default:
				progressStr := progressFunc()
				os.Stderr.WriteString(progressStr)
				time.Sleep(1 * time.Second)
				clearLine()
			}
		}
	}()
	return func() {
		cancel <- true
		clearLine()
	}
}

func clearLine() {
	os.Stderr.WriteString("\033[2K\r")
}

func progressBar(max int, positionFn func() int, label string) func() string {
	return func() string {
		progress := 100 * float64(positionFn()) / float64(max)
		blocksNumber := int(progress / 2)
		blocks := strings.Repeat("▮", blocksNumber) + strings.Repeat("▯", 50-blocksNumber)
		str := fmt.Sprintf("%s: %.2f%% [%s]", label, progress, blocks)
		return str
	}
}
