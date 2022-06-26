package main

import (
	"fmt"
	"os"

	"github.com/danielleontiev/neojhat/cmd"
)

func main() {
	if len(os.Args) < 2 {
		cmd.PrintHelp()
	}
	args := os.Args[2:]
	switch os.Args[1] {
	case cmd.Threads:
		cmd.ThreadsCommand.Parse(args)
		threads()
	case cmd.Summary:
		cmd.SummaryCommand.Parse(args)
		summary()
	case cmd.Objects:
		cmd.ObjectsCommand.Parse(args)
		objects()
	default:
		cmd.PrintHelp()
	}
}

func threads() {
	if cmd.ThreadFlags.Hprof == "" {
		cmd.PrintUsage(cmd.ThreadsCommand)
	}
	flags := cmd.ThreadFlags
	if err := cmd.ParseHprof(flags.Hprof, flags.NonInteractive); err != nil {
		onError(err)
	}
	if err := cmd.GetThreads(flags.Hprof, flags.NoColor, flags.LocalVars); err != nil {
		onError(err)
	}
}

func summary() {
	if cmd.SummaryFlags.Hprof == "" {
		cmd.PrintUsage(cmd.SummaryCommand)
	}
	flags := cmd.SummaryFlags
	if err := cmd.ParseHprof(flags.Hprof, flags.NonInteractive); err != nil {
		onError(err)
	}
	if err := cmd.GetSummary(flags.Hprof, flags.NoColor, flags.AllProps); err != nil {
		onError(err)
	}
}

func objects() {
	if cmd.ObjectsFlags.Hprof == "" {
		cmd.PrintUsage(cmd.ObjectsCommand)
	}
	flags := cmd.ObjectsFlags
	if err := cmd.ParseHprof(flags.Hprof, flags.NonInteractive); err != nil {
		onError(err)
	}
	if err := cmd.GetObjects(flags.Hprof, flags.NoColor, flags.SortBy); err != nil {
		onError(err)
	}
}

func onError(err error) {
	fmt.Printf("Error occurred: %v", err)
	os.Exit(1)
}
