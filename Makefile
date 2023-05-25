## test: runs all test
test:
	@go test -v ./...

## cover: opens coverage in browser
cover:
	@go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

## coverage: displays test coverage
coverage:
	@go test -cover ./..

## build_cli: builds the command line tool celeritas and copies it to myapp
build_cli:
	@go build -o ../myapp/celeritas ./cmd/cli

## build: builds the command line tool to dist directory
build: build_windows build_linux

build_windows:
	@GOOS=windows GOARCH=amd64 go build -o ./dist/cli.exe ./cmd/cli

build_linux:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./dist/cli ./cmd/cli