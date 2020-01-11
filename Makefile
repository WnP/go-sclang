build:
	mkdir -p bin
	go build -ldflags="-s -w" -o ./bin/go-sclang ./cmd/go-sclang
	go build -ldflags="-s -w" -o ./bin/go-sclang-client ./cmd/go-sclang-client

compress:
	upx --brute ./bin/go-sclang
	upx --brute ./bin/go-sclang-client

clean:
	rm bin/* || return 0
