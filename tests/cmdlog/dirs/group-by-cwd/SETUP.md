# Scenario

**Feature**: dirs groups today's commands under each cwd

```
# events in /alpha and /beta
cmdlog dirs -> sections per cwd with command lists
```

## Steps

1. Seed events for local today in `/alpha` (2 commands) and `/beta` (1 command).
2. Run `cmdlog dirs` (default today).

```go
func Setup(t *testing.T, req *Request) error {
	req.SeedEvents = []Event{
		{TS: utcTodayAt(9, 0), CWD: "/alpha", CMD: "npm install"},
		{TS: utcTodayAt(9, 15), CWD: "/alpha", CMD: "npm test"},
		{TS: utcTodayAt(10, 0), CWD: "/beta", CMD: "cargo build"},
	}
	return nil
}
```