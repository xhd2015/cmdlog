# Scenario

**Feature**: malformed JSONL line skipped with stderr warning

```
# fixture: one corrupt line, one valid today's event
cmdlog today -> stderr warning + stdout includes valid command only
```

## Steps

1. Seed `testdata/malformed-events.jsonl` (written relative to leaf dir at assert time;
   Run loads via `EventsFixtureFile` from leaf SETUP — need to fix path).

## Context

- Valid line command: `valid-command`.
- Invalid line is not valid JSON.

```go
func Setup(t *testing.T, req *Request) error {
	goodLine := fmt.Sprintf(`{"ts":%q,"cwd":"/tmp","cmd":"valid-command"}`, utcTodayAt(9, 52))
	req.EventsRawContent = "{not valid json at all}\n" + goodLine
	return nil
}
```