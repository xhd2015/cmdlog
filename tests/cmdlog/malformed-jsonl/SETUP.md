# Scenario

**Feature**: query commands tolerate malformed JSONL lines

```
events.jsonl with bad line + good line -> cmdlog today -> warning on stderr, good line processed
```

## Preconditions

- `events.jsonl` contains at least one invalid JSON line and one valid line.

## Steps

1. Set `req.Subcommand = "malformed"` (runs `cmdlog today` against corrupt file).
2. Load fixture from `testdata/malformed-events.jsonl`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Subcommand = "malformed"
	return nil
}
```