# Go SCLang

Provide an HTTP server and an HTTP client for [SuperCollider sclang][1].

See the server and client readmes for more informations.

## Install

Download the latest release from the release and add them to your path, eg.:

```console
$ # Fist create a temporary directory
$ TMP_DIR=`mktemp -d`
$ cd $TMP_DIR
$ # Download go-sclang
$ wget https://github.com/WnP/go-sclang/releases/download/1.0.0/go-sclang-linux-amd64
$ # Check if checksum is correct
$ wget https://github.com/WnP/go-sclang/releases/download/1.0.0/go-sclang-linux-amd64-sha256sum.txt
$ cat go-sclang-linux-amd64-sha256sum.txt | sha256sum -c -
go-sclang-linux-amd64: OK
$ # Now make it executable and mv it to your path
$ chmod +x go-sclang-linux-amd64 && mv go-sclang-linux-amd64 ~/.local/bin/go-sclang

$ # Do the same for go-sclang-client
$ wget https://github.com/WnP/go-sclang/releases/download/1.0.0/go-sclang-client-linux-amd64
$ wget https://github.com/WnP/go-sclang/releases/download/1.0.0/go-sclang-client-linux-amd64-sha256sum.txt
$ cat go-sclang-client-linux-amd64-sha256sum.txt | sha256sum -c -
go-sclang-client-linux-amd64: OK
$ chmod +x go-sclang-client-linux-amd64 && mv go-sclang-client-linux-amd64 ~/.local/bin/go-sclang-client

$ Clean your temporary directory
$ cd ~ && rm -rf $TMP_DIR
```

## Build from source

Require [go][2] `1.13.x`.

Note that the compress phase is optional and require [upx][3],
but it reduce the final binaries size to ~30% of the original size.

```console
$ git clone https://github.com/WnP/go-sclang.git
$ cd go-sclang
$ make build
$ make compress  # Optional
$ mv ./bin/* ~/.local/bin
```

[1]: https://supercollider.github.io/
[2]: https://golang.org/
[3]: https://upx.github.io/
