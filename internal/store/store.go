package store

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xhd2015/cmdlog/internal/event"
)

const eventsFileName = "events.jsonl"

// EventsPath returns the default events.jsonl path under HOME.
func EventsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".cmdlog", eventsFileName), nil
}

// Append writes one event to events.jsonl, creating ~/.cmdlog if needed.
func Append(ev event.Event) error {
	path, err := EventsPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	if ev.TS.IsZero() {
		ev.TS = time.Now().UTC()
	}
	b, err := json.Marshal(struct {
		TS  string `json:"ts"`
		CWD string `json:"cwd"`
		CMD string `json:"cmd"`
	}{
		TS:  ev.TS.UTC().Format(time.RFC3339),
		CWD: ev.CWD,
		CMD: ev.CMD,
	})
	if err != nil {
		return err
	}
	if _, err := f.Write(append(b, '\n')); err != nil {
		return err
	}
	return nil
}

// ReadAll loads valid events and invokes warn for each malformed line.
func ReadAll(warn func(lineNum int, err error)) ([]event.Event, error) {
	path, err := EventsPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var events []event.Event
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var raw struct {
			TS  string `json:"ts"`
			CWD string `json:"cwd"`
			CMD string `json:"cmd"`
		}
		if err := json.Unmarshal([]byte(line), &raw); err != nil {
			if warn != nil {
				warn(lineNum, err)
			}
			continue
		}
		ts, err := time.Parse(time.RFC3339, raw.TS)
		if err != nil {
			if warn != nil {
				warn(lineNum, err)
			}
			continue
		}
		events = append(events, event.Event{
			TS:  ts,
			CWD: raw.CWD,
			CMD: raw.CMD,
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

// FilterByLocalDate returns events whose timestamp falls on date in local TZ.
func FilterByLocalDate(events []event.Event, date time.Time) []event.Event {
	loc := date.Location()
	y, m, d := date.Date()
	start := time.Date(y, m, d, 0, 0, 0, 0, loc)
	end := start.Add(24 * time.Hour)
	var out []event.Event
	for _, ev := range events {
		local := ev.TS.In(loc)
		if !local.Before(start) && local.Before(end) {
			out = append(out, ev)
		}
	}
	return out
}

// ParseDate parses YYYY-MM-DD in local timezone.
func ParseDate(s string) (time.Time, error) {
	t, err := time.ParseInLocation("2006-01-02", s, time.Local)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date %q: %w", s, err)
	}
	return t, nil
}