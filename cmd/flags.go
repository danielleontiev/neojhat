package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/danielleontiev/neojhat/objects"
)

const (
	threadsSubCmdName = "threads"
	summarySubCmdName = "summary"
	objectsSubCmdName = "objects"
)

var (
	threadsSubCmd = flag.NewFlagSet(threadsSubCmdName, flag.ExitOnError)
	summarySubCmd = flag.NewFlagSet(summarySubCmdName, flag.ExitOnError)
	objectsSubCmd = flag.NewFlagSet(objectsSubCmdName, flag.ExitOnError)
)

func init() {
	threadsSubCmd.SetOutput(os.Stdout)
	summarySubCmd.SetOutput(os.Stdout)
	objectsSubCmd.SetOutput(os.Stdout)
}

func printHelp() {
	fmt.Printf("neojhat (%s|%s|%s)\n\n", threadsSubCmdName, summarySubCmdName, objectsSubCmdName)
	threadsSubCmd.Usage()
	fmt.Println()
	summarySubCmd.Usage()
	fmt.Println()
	objectsSubCmd.Usage()
	os.Exit(0)
}

func printUsage(subCommand *flag.FlagSet) {
	subCommand.Usage()
	os.Exit(0)
}

const (
	hprofName    = "hprof"
	hprofDefault = ""
	hprofDesc    = "path to .hprof file (required)"

	noColorName    = "no-color"
	noColorDefault = false
	noColorDesc    = "disable color output"

	nonInteractiveName    = "non-interactive"
	nonInteractiveDefault = false
	nonInteractiveDesc    = "disable interactive output"

	allPropsName    = "all-props"
	allPropsDefault = false
	allPropsDesc    = "print all available properties from java.lang.System"

	localVarsName    = "local-vars"
	localVarsDefault = false
	localVarsDesc    = "show local variables"

	sortByName = "sort-by"
	sortByDesc = "Sort output by 'size' or 'count' (default)"
)

type threadFlags struct {
	hprof          string
	noColor        bool
	nonInteractive bool
	localVars      bool
}

type summaryFlags struct {
	hprof          string
	noColor        bool
	nonInteractive bool
	allProps       bool
}

type objectsFlags struct {
	hprof          string
	noColor        bool
	nonInteractive bool
	sortBy         objects.SortBy
}

var (
	threadFlagsValues  threadFlags
	summaryFlagsValues summaryFlags
	objectsFlagsValues objectsFlags
)

func init() {
	threadsSubCmd.StringVar(&threadFlagsValues.hprof, hprofName, hprofDefault, hprofDesc)
	threadsSubCmd.BoolVar(&threadFlagsValues.noColor, noColorName, noColorDefault, noColorDesc)
	threadsSubCmd.BoolVar(&threadFlagsValues.nonInteractive, nonInteractiveName, nonInteractiveDefault, nonInteractiveDesc)
	threadsSubCmd.BoolVar(&threadFlagsValues.localVars, localVarsName, localVarsDefault, localVarsDesc)

	summarySubCmd.StringVar(&summaryFlagsValues.hprof, hprofName, hprofDefault, hprofDesc)
	summarySubCmd.BoolVar(&summaryFlagsValues.noColor, noColorName, noColorDefault, noColorDesc)
	summarySubCmd.BoolVar(&summaryFlagsValues.nonInteractive, nonInteractiveName, nonInteractiveDefault, nonInteractiveDesc)
	summarySubCmd.BoolVar(&summaryFlagsValues.allProps, allPropsName, allPropsDefault, allPropsDesc)

	objectsSubCmd.StringVar(&objectsFlagsValues.hprof, hprofName, hprofDefault, hprofDesc)
	objectsSubCmd.BoolVar(&objectsFlagsValues.noColor, noColorName, noColorDefault, noColorDesc)
	objectsSubCmd.BoolVar(&objectsFlagsValues.nonInteractive, nonInteractiveName, nonInteractiveDefault, nonInteractiveDesc)
	objectsSubCmd.Var(&objectsFlagsValues.sortBy, sortByName, sortByDesc)
}
