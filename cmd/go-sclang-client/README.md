# go-sclang-client

```
▶ go-sclang-client --help

Usage of go-sclang-client:
Reads code from stdin and send it to go-sclang server, eg:

        echo '"Hello from go-sclang-client".postln' | go-sclang-client

Flags:

  -h string
        -host shorthand  (default "http://localhost")
  -host string
        Server host (default "http://localhost")
  -k    -kill shorthand
  -kill
        If true kill sclang server
  -o    -stdout shorthand
  -p int
        -port shorthand  (default 5533)
  -port int
        Server port (default 5533)
  -r    -kill shorthand
  -reload
        If true reload sclang server by sending SIGUSR1
  -stdout
        If true return sclang value
  -t string
        -timeout shorthand (default "10s")
  -timeout string
        sclang server timeout (default "10s")
```
