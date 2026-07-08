# Scenario

**Feature**: cmdlog today shows chronological timeline for local today

```
cmdlog today -> read events.jsonl -> filter local today -> group by hour+cwd -> stdout
```

## Preconditions

- Events seeded for local today's date (or empty/missing log per leaf).

## Steps

1. Set `req.Subcommand = "today"`.
2. Seed fixture events when leaf requires non-empty log.

## Context

- Empty or missing log must print a friendly message and exit 0.
- User-facing stdout must end with trailing newline.

```go
func Setup(t *testing.T, req *Request) error {
	req.Subcommand = "today"
	return nil
}
```