package event

import "time"

// Event is one append-only JSONL record.
type Event struct {
	TS  time.Time `json:"ts"`
	CWD string    `json:"cwd"`
	CMD string    `json:"cmd"`
}