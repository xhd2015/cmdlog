# Scenario

**Feature**: cmdlog dirs groups commands by working directory for a date

```
cmdlog dirs [--date YYYY-MM-DD] -> group by cwd -> stdout
```

## Steps

1. Set `req.Subcommand = "dirs"`.
2. Seed events across multiple cwds for the target date.

```go
func Setup(t *testing.T, req *Request) error {
	req.Subcommand = "dirs"
	return nil
}
```