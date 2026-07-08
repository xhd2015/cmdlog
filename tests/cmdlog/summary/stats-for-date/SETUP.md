# Scenario

**Feature**: summary shows counts and top directory for a date

```
# 4 commands: 3 in /work/main, 1 in /work/other
cmdlog summary --date <today> -> "4 commands, 2 directories, top: /work/main (3)"
```

## Steps

1. Seed four events for local today across two directories.
2. Run `cmdlog summary --date <local today>`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Date = localTodayDate()
	req.SeedEvents = []Event{
		{TS: utcTodayAt(8, 0), CWD: "/work/main", CMD: "git pull"},
		{TS: utcTodayAt(9, 0), CWD: "/work/main", CMD: "go build"},
		{TS: utcTodayAt(10, 0), CWD: "/work/main", CMD: "go test"},
		{TS: utcTodayAt(11, 0), CWD: "/work/other", CMD: "ls"},
	}
	return nil
}
```