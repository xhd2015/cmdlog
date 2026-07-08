# Scenario

**Feature**: cmdlog timeline filters events by time range

```
cmdlog timeline [--from RFC3339] [--to RFC3339] -> filtered chronological output
```

## Steps

1. Set `req.Subcommand = "timeline"`.
2. Seed events spanning a known time range; set `--from` / `--to` per leaf.

```go
func Setup(t *testing.T, req *Request) error {
	req.Subcommand = "timeline"
	return nil
}
```