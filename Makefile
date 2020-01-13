build:
	mkdir -p bin
	go build -ldflags="-s -w" -o ./bin/go-sclang ./cmd/go-sclang
	go build -ldflags="-s -w" -o ./bin/go-sclang-client ./cmd/go-sclang-client

compress:
	upx --brute ./bin/go-sclang
	upx --brute ./bin/go-sclang-client

clean:
	rm bin/* || return 0
	rm -rf dist

dist:
	for platform in linux darwin; do \
		for arch in amd64 arm; do \
			if [ $$platform != darwin ] || [ $$arch != arm ]; then \
				for bin in go-sclang go-sclang-client; do \
					GOOS=$$platform GOARCH=$$arch go build -ldflags="-s -w" -o ./dist/$$bin-$$platform-$$arch ./cmd/$$bin; \
				done; \
			fi; \
		done; \
	done

	for f in `find dist -executable -type f` ; do \
		upx --brute $$f; \
		sha256sum $$f > $$f-sha256sum.txt; \
	done
