package main

import (
	"fmt"
	"os"

	lessflags "github.com/xhd2015/less-flags"
	"github.com/xhd2015/cmdlog/internal/cmd"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Print(cmd.HelpText)
		os.Exit(0)
	}

	args, err := lessflags.Help("-h,--help", cmd.HelpText).
		StopOnFirstArg().
		Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if len(args) == 0 {
		fmt.Print(cmd.HelpText)
		os.Exit(0)
	}

	subcommand := args[0]
	remain := args[1:]

	var runErr error
	switch subcommand {
	case "record":
		runErr = cmd.RunRecord(remain)
	case "today":
		runErr = cmd.RunToday(remain)
	case "summary":
		runErr = cmd.RunSummary(remain)
	case "timeline":
		runErr = cmd.RunTimeline(remain)
	case "dirs":
		runErr = cmd.RunDirs(remain)
	case "integration":
		runErr = cmd.RunIntegration(remain)
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand %q\n", subcommand)
		fmt.Fprint(os.Stderr, cmd.HelpText)
		os.Exit(1)
	}

	if runErr != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", runErr)
		os.Exit(1)
	}
}