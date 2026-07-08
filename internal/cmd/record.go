package cmd

import (
	"fmt"
	"strings"

	lessflags "github.com/xhd2015/less-flags"
	"github.com/xhd2015/cmdlog/internal/event"
	"github.com/xhd2015/cmdlog/internal/store"
)

func RunRecord(args []string) error {
	var cwd string
	remain, err := lessflags.String("--cwd", &cwd).StopOnFirstArg().Parse(args)
	if err != nil {
		return err
	}

	command := extractCommand(remain)
	if cwd == "" {
		return fmt.Errorf("record requires --cwd=<dir>")
	}
	if command == "" {
		return fmt.Errorf("record requires a command after --")
	}

	return store.Append(event.Event{
		CWD: cwd,
		CMD: command,
	})
}

func extractCommand(args []string) string {
	for i, arg := range args {
		if arg == "--" {
			return strings.Join(args[i+1:], " ")
		}
	}
	return strings.Join(args, " ")
}