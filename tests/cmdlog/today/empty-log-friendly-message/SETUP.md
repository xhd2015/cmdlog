# Scenario

**Feature**: today handles empty or missing log gracefully

```
# no events.jsonl exists
cmdlog today -> friendly message -> exit 0
```

## Steps

1. Do not seed any events (fresh isolated HOME).
2. Run `cmdlog today`.

## Context

- Must not error or panic on missing file.

```go
func Setup(t *testing.T, req *Request) error {
	req.SeedEvents = nil
	req.EventsFixtureFile = ""
	req.EventsRawContent = ""
	return nil
}
```