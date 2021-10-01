# GOOS：darwin、freebsd、linux、windows
# GOARCH：386、amd64、arm、s390x

all: darwin darwin_arm64 linux mips64le mipsle windows

darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags "-w -s" -o bin/BUCTNet-Login_macOS_amd64 ./cmd

darwin_arm64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags "-w -s" -o bin/BUCTNet-Login_macOS_arm64 ./cmd

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-w -s" -o bin/BUCTNet-Login_linux_amd64 ./cmd

mips64le:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "-w -s" -o bin/BUCTNet-Login_linux_mips64le ./cmd

mipsle:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "-w -s" -o bin/BUCTNet-Login_linux_mipsle ./cmd

windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "-w -s" -o bin/BUCTNet-Login_windows_amd64.exe ./cmd

clean:
	rm -rf ./bin

# .PHONY:publish
# .PHONY:darwin
# .PHONY:linux
# .PHONY:windows
# .PHONY:clean