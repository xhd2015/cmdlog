package cmd

import (
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	lessflags "github.com/xhd2015/less-flags"
	"github.com/xhd2015/cmdlog/internal/event"
	"github.com/xhd2015/cmdlog/internal/store"
)

func warnMalformed(w io.Writer, lineNum int, err error) {
	fmt.Fprintf(w, "warning: skipping malformed JSONL line %d: %v\n", lineNum, err)
}

func loadEvents() ([]event.Event, error) {
	return store.ReadAll(func(lineNum int, err error) {
		warnMalformed(os.Stderr, lineNum, err)
	})
}

func RunToday() error {
	events, err := loadEvents()
	if err != nil {
		return err
	}
	today := time.Now().In(time.Local)
	filtered := store.FilterByLocalDate(events, today)
	if len(filtered) == 0 {
		fmt.Println("No commands recorded today.")
		return nil
	}
	printTimeline(os.Stdout, filtered)
	return nil
}

func RunSummary(args []string) error {
	var dateFlag string
	_, err := lessflags.String("--date", &dateFlag).StopOnFirstArg().Parse(args)
	if err != nil {
		return err
	}

	date := time.Now().In(time.Local)
	if dateFlag != "" {
		date, err = store.ParseDate(dateFlag)
		if err != nil {
			return err
		}
	}

	events, err := loadEvents()
	if err != nil {
		return err
	}
	filtered := store.FilterByLocalDate(events, date)
	if len(filtered) == 0 {
		fmt.Printf("0 commands, 0 directories, top: (none)\n")
		return nil
	}

	dirCounts := map[string]int{}
	for _, ev := range filtered {
		dirCounts[ev.CWD]++
	}

	topDir := ""
	topCount := 0
	dirs := make([]string, 0, len(dirCounts))
	for dir, count := range dirCounts {
		dirs = append(dirs, dir)
		if count > topCount || (count == topCount && (topDir == "" || dir < topDir)) {
			topDir = dir
			topCount = count
		}
	}

	fmt.Printf("%d commands, %d directories, top: %s (%d)\n",
		len(filtered), len(dirs), topDir, topCount)
	return nil
}

func RunTimeline(args []string) error {
	var fromFlag, toFlag string
	_, err := lessflags.String("--from", &fromFlag).String("--to", &toFlag).StopOnFirstArg().Parse(args)
	if err != nil {
		return err
	}

	events, err := loadEvents()
	if err != nil {
		return err
	}

	var from, to time.Time
	if fromFlag != "" {
		from, err = time.Parse(time.RFC3339, fromFlag)
		if err != nil {
			return fmt.Errorf("invalid --from: %w", err)
		}
	}
	if toFlag != "" {
		to, err = time.Parse(time.RFC3339, toFlag)
		if err != nil {
			return fmt.Errorf("invalid --to: %w", err)
		}
	}

	filtered := filterByRange(events, from, to, fromFlag == "", toFlag == "")
	if len(filtered) == 0 {
		fmt.Println("No commands in the selected time range.")
		return nil
	}
	printTimeline(os.Stdout, filtered)
	return nil
}

func RunDirs(args []string) error {
	var dateFlag string
	_, err := lessflags.String("--date", &dateFlag).StopOnFirstArg().Parse(args)
	if err != nil {
		return err
	}

	date := time.Now().In(time.Local)
	if dateFlag != "" {
		date, err = store.ParseDate(dateFlag)
		if err != nil {
			return err
		}
	}

	events, err := loadEvents()
	if err != nil {
		return err
	}
	filtered := store.FilterByLocalDate(events, date)
	if len(filtered) == 0 {
		fmt.Println("No commands recorded for this date.")
		return nil
	}
	printDirs(os.Stdout, filtered)
	return nil
}

func filterByRange(events []event.Event, from, to time.Time, noFrom, noTo bool) []event.Event {
	var out []event.Event
	for _, ev := range events {
		if !noFrom && ev.TS.Before(from) {
			continue
		}
		if !noTo && ev.TS.After(to) {
			continue
		}
		out = append(out, ev)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].TS.Before(out[j].TS)
	})
	return out
}

func printTimeline(w io.Writer, events []event.Event) {
	sort.Slice(events, func(i, j int) bool {
		return events[i].TS.Before(events[j].TS)
	})

	loc := time.Local
	type cwdKey struct {
		hour string
		cwd  string
	}
	groups := make([]cwdKey, 0)
	commands := map[cwdKey][]string{}

	for _, ev := range events {
		local := ev.TS.In(loc)
		hour := fmt.Sprintf("%02d:00", local.Hour())
		key := cwdKey{hour: hour, cwd: ev.CWD}
		if _, ok := commands[key]; !ok {
			groups = append(groups, key)
		}
		commands[key] = append(commands[key], ev.CMD)
	}

	for _, key := range groups {
		fmt.Fprintf(w, "%s %s\n", key.hour, key.cwd)
		for _, cmd := range commands[key] {
			fmt.Fprintf(w, "  %s\n", cmd)
		}
	}
}

func printDirs(w io.Writer, events []event.Event) {
	sort.Slice(events, func(i, j int) bool {
		if events[i].CWD != events[j].CWD {
			return events[i].CWD < events[j].CWD
		}
		return events[i].TS.Before(events[j].TS)
	})

	type group struct {
		cwd  string
		cmds []string
	}
	var groups []group
	index := map[string]int{}

	for _, ev := range events {
		i, ok := index[ev.CWD]
		if !ok {
			groups = append(groups, group{cwd: ev.CWD})
			i = len(groups) - 1
			index[ev.CWD] = i
		}
		groups[i].cmds = append(groups[i].cmds, ev.CMD)
	}

	for _, g := range groups {
		fmt.Fprintf(w, "%s\n", g.cwd)
		for _, cmd := range g.cmds {
			fmt.Fprintf(w, "  %s\n", cmd)
		}
	}
}