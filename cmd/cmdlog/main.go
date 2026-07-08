package main

import (
	"fmt"
	"os"

	"github.com/xhd2015/cmdlog/internal/cmd"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: cmdlog <record|today|summary|timeline|dirs|integration> ...")
		os.Exit(1)
	}

	subcommand := os.Args[1]
	args := os.Args[2:]

	var err error
	switch subcommand {
	case "record":
		err = cmd.RunRecord(args)
	case "today":
		err = cmd.RunToday()
	case "summary":
		err = cmd.RunSummary(args)
	case "timeline":
		err = cmd.RunTimeline(args)
	case "dirs":
		err = cmd.RunDirs(args)
	case "integration":
		err = cmd.RunIntegration(args)
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand %q\n", subcommand)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}