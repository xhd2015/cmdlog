# Scenario

**Feature**: timeline --from/--to returns only matching events

```
# events at 08:00, 09:00, 10:00 UTC today; filter 08:30–09:30 -> only 09:00 event
cmdlog timeline --from ... --to ... -> single matching command
```

## Steps

1. Seed three events at 08:00, 09:00, and 10:00 local today (stored as UTC).
2. Run timeline with `--from` and `--to` bracketing only the 09:00 event.

```go
func Setup(t *testing.T, req *Request) error {
	req.SeedEvents = []Event{
		{TS: utcTodayAt(8, 0), CWD: "/tmp", CMD: "before-window"},
		{TS: utcTodayAt(9, 0), CWD: "/tmp", CMD: "inside-window"},
		{TS: utcTodayAt(10, 0), CWD: "/tmp", CMD: "after-window"},
	}
	req.From = utcTodayAt(8, 30)
	req.To = utcTodayAt(9, 30)
	return nil
}
```