# Scenario

**Feature**: status reports partial when bash.sh exists without profile marker

```
# bash.sh present, no marker in .bash_profile
cmdlog integration bash --status -> partial, exit 1
```

## Steps

1. Pre-seed `bash.sh` only (no profile marker).

```go
func Setup(t *testing.T, req *Request) error {
	req.PreExistingBashSh = minimalBashSh()
	return nil
}
```