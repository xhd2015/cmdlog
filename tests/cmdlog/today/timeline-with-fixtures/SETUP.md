# Scenario

**Feature**: today prints hour+cwd grouped timeline for fixture events

```
# seed today's events in two hours and two cwds
cmdlog today -> timeline grouped by hour and cwd
```

## Steps

1. Seed three events for local today at 09:00 and 10:00 in `/proj/a` and `/proj/b`.
2. Run `cmdlog today`.

## Context

- Fixtures use `utcTodayAt()` so events always fall on local today regardless of test run date.

```go
func Setup(t *testing.T, req *Request) error {
	req.SeedEvents = []Event{
		{TS: utcTodayAt(9, 5), CWD: "/proj/a", CMD: "git status"},
		{TS: utcTodayAt(9, 30), CWD: "/proj/a", CMD: "go test ./..."},
		{TS: utcTodayAt(10, 15), CWD: "/proj/b", CMD: "make build"},
	}
	return nil
}
```