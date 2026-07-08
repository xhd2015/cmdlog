# Scenario

**Feature**: cmdlog summary reports command and directory statistics

```
cmdlog summary [--date YYYY-MM-DD] -> read events -> stats stdout
```

## Steps

1. Set `req.Subcommand = "summary"`.
2. Seed fixture events for the target date.

## Context

- Output format: `N commands, M directories, top: <path> (count)`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Subcommand = "summary"
	return nil
}
```