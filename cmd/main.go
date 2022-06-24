package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
	}
	args := os.Args[2:]
	switch os.Args[1] {
	case threadsSubCmdName:
		threadsSubCmd.Parse(args)
		runThreadsSubCmd()
	case summarySubCmdName:
		summarySubCmd.Parse(args)
		runSummarySubCmd()
	case objectsSubCmdName:
		objectsSubCmd.Parse(args)
		runObjectsSubCmd()
	default:
		printHelp()
	}
}

func runThreadsSubCmd() {
	if threadFlagsValues.hprof == "" {
		printUsage(threadsSubCmd)
	}
	flags := threadFlagsValues
	if err := parseHprof(flags.hprof, flags.nonInteractive); err != nil {
		onError(err)
	}
	if err := getThreads(flags.hprof, flags.noColor, flags.localVars); err != nil {
		onError(err)
	}
}

func runSummarySubCmd() {
	if summaryFlagsValues.hprof == "" {
		printUsage(summarySubCmd)
	}
	flags := summaryFlagsValues
	if err := parseHprof(flags.hprof, flags.nonInteractive); err != nil {
		onError(err)
	}
	if err := getSummary(flags.hprof, flags.noColor, flags.allProps); err != nil {
		onError(err)
	}
}

func runObjectsSubCmd() {
	if objectsFlagsValues.hprof == "" {
		printUsage(objectsSubCmd)
	}
	flags := objectsFlagsValues
	if err := parseHprof(flags.hprof, flags.nonInteractive); err != nil {
		onError(err)
	}
	if err := getObjects(flags.hprof, flags.noColor, flags.sortBy); err != nil {
		onError(err)
	}
}

func onError(err error) {
	fmt.Printf("Error occurred: %v", err)
	os.Exit(1)
}
