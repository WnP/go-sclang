# go-sclang

```
▶ go-sclang --help

Usage of go-sclang:
Wrap sclang CLI in order to provide an HTTP API, expose:

POST application/json {"Code": "'sclang code'.postln", "Stdout": true}

if Stdout is true then sclang return value is returned.

Flags:

  -b int
        -buffer-size shorthand (default 1024)
  -buffer-size int
        Stdout and Stderr buffer size (default 1024)
  -h string
        -host shorthand (default "localhost")
  -host string
        Server host (default "localhost")
  -p int
        -port shorthand (default 5533)
  -port int
        Server port (default 5533)
  -t string
        -timeout shorthand (default "4s")
  -timeout string
        sclang server timeout (default "4s")
```

Here is the query model:

```go
// Query represent go-sclang HTTP query
type Query struct {
	// Code string to pass to sclang
	Code string
	// if true sclang return value is returned
	Stdout bool
	// if true kill sclang and sc-synth server
	Kill bool
	// if true kill sc-synth and reload sclang (sending SIGUSR1)
	Reload bool
}
```
