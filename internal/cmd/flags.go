package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/danielleontiev/neojhat/internal/objects"
)

const (
	version = "v0.2.0"
)

const (
	Threads = "threads"
	Summary = "summary"
	Objects = "objects"
)

var (
	ThreadsCommand = flag.NewFlagSet(Threads, flag.ExitOnError)
	SummaryCommand = flag.NewFlagSet(Summary, flag.ExitOnError)
	ObjectsCommand = flag.NewFlagSet(Objects, flag.ExitOnError)
)

func init() {
	ThreadsCommand.SetOutput(os.Stdout)
	SummaryCommand.SetOutput(os.Stdout)
	ObjectsCommand.SetOutput(os.Stdout)

	ThreadsCommand.StringVar(&ThreadFlags.Hprof, hprofName, hprofDefault, hprofDesc)
	ThreadsCommand.BoolVar(&ThreadFlags.NoColor, noColorName, noColorDefault, noColorDesc)
	ThreadsCommand.BoolVar(&ThreadFlags.NonInteractive, nonInteractiveName, nonInteractiveDefault, nonInteractiveDesc)
	ThreadsCommand.BoolVar(&ThreadFlags.LocalVars, localVarsName, localVarsDefault, localVarsDesc)
	ThreadsCommand.Var(&ThreadFlags.Output, outputName, outputDesc)

	SummaryCommand.StringVar(&SummaryFlags.Hprof, hprofName, hprofDefault, hprofDesc)
	SummaryCommand.BoolVar(&SummaryFlags.NoColor, noColorName, noColorDefault, noColorDesc)
	SummaryCommand.BoolVar(&SummaryFlags.NonInteractive, nonInteractiveName, nonInteractiveDefault, nonInteractiveDesc)
	SummaryCommand.BoolVar(&SummaryFlags.AllProps, allPropsName, allPropsDefault, allPropsDesc)
	SummaryCommand.Var(&SummaryFlags.Output, outputName, outputDesc)

	ObjectsCommand.StringVar(&ObjectsFlags.Hprof, hprofName, hprofDefault, hprofDesc)
	ObjectsCommand.BoolVar(&ObjectsFlags.NoColor, noColorName, noColorDefault, noColorDesc)
	ObjectsCommand.BoolVar(&ObjectsFlags.NonInteractive, nonInteractiveName, nonInteractiveDefault, nonInteractiveDesc)
	ObjectsCommand.Var(&ObjectsFlags.SortBy, sortByName, sortByDesc)
	ObjectsCommand.Var(&ObjectsFlags.Output, outputName, outputDesc)
}

func PrintHelp() {
	fmt.Printf("neojhat %s\n", version)
	fmt.Printf("neojhat (%s|%s|%s)\n\n", Threads, Summary, Objects)
	ThreadsCommand.Usage()
	fmt.Println()
	SummaryCommand.Usage()
	fmt.Println()
	ObjectsCommand.Usage()
	os.Exit(0)
}

func PrintUsage(subCommand *flag.FlagSet) {
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

	outputName = "output"
	outputDesc = "Output type. 'plain' (default) or 'html'"
)

type OutputType int

const (
	Plain OutputType = iota
	Html
)

func (o *OutputType) String() string {
	switch *o {
	case Plain:
		return "plain"
	case Html:
		return "html"
	}
	return "unknown"
}

func (o *OutputType) Set(value string) error {
	switch value {
	case "plain":
		*o = Plain
		return nil
	case "html":
		*o = Html
		return nil
	case "":
		*o = Plain
		return nil
	}
	return fmt.Errorf("Use \"plain\" or \"html\" instead")
}

type threadFlags struct {
	Hprof          string
	NoColor        bool
	NonInteractive bool
	LocalVars      bool
	Output         OutputType
}

type summaryFlags struct {
	Hprof          string
	NoColor        bool
	NonInteractive bool
	AllProps       bool
	Output         OutputType
}

type objectsFlags struct {
	Hprof          string
	NoColor        bool
	NonInteractive bool
	SortBy         objects.SortBy
	Output         OutputType
}

var (
	ThreadFlags  threadFlags
	SummaryFlags summaryFlags
	ObjectsFlags objectsFlags
)
