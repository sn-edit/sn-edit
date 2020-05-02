binary_name = sn-edit

.PHONY: all
all: linux_amd64 darwin_amd64 windows_amd64 checksums

.PHONY: linux_amd64
linux_amd64:
	CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -v -a -gcflags=-trimpath=$$HOME -asmflags=-trimpath=$$HOME -o build/$(binary_name)-linux-amd64 -ldflags "-linkmode external -extldflags -static"

.PHONY: darwin_amd64
darwin_amd64:
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -v -a -gcflags "all=-trimpath=$$PWD;$$HOME" -asmflags "all=-trimpath=$$PWD;$$HOME" -o build/$(binary_name)-darwin-amd64

.PHONY: windows_amd64
windows_amd64:
	CGO_ENABLED=1 CC=/usr/local/bin/x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -v -a -gcflags=-trimpath=$$HOME -asmflags=-trimpath=$$HOME -o build/$(binary_name)-windows-amd64.exe

.PHONY: checksums
checksums:
	shasum -a 256 build/* > build/checksum.txt