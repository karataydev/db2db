build:
	@go build -o bin/db2db ./cmd/db2db
run: build
	@./bin/db2db -s sqlserver://sa:1234Belma.@localhost:1433?database=fro -t sqlserver://sa:1234Belma.@localhost:1432?database=to -ct=f
test:
	@go test -v ./..

compile:
	GOOS=freebsd GOARCH=386 go build -o bin/db2db-freebsd-386 main.go
	GOOS=linux GOARCH=386 go build -o bin/db2db-linux-386 main.go
	GOOS=windows GOARCH=386 go build -o bin/db2db-windows-386 main.go

build-all: windows linux darwin
	@echo version: $(VERSION)

EXECUTABLE=db2db
VERSION=$(shell git describe --tags)
WINDOWS=$(EXECUTABLE)_windows_amd64_$(VERSION).exe
LINUX=$(EXECUTABLE)_linux_amd64_$(VERSION)
DARWIN=$(EXECUTABLE)_darwin_amd64_$(VERSION)

windows: $(WINDOWS)

linux: $(LINUX)

darwin: $(DARWIN)

$(WINDOWS):
	env GOOS=windows GOARCH=amd64 go build -v -o bin/$(WINDOWS) -ldflags="-s -w -X main.version=$(VERSION)" ./cmd/db2db/main.go

$(LINUX):
	env GOOS=linux GOARCH=amd64 go build -v -o bin/$(LINUX) -ldflags="-s -w -X main.version=$(VERSION)" ./cmd/db2db/main.go

$(DARWIN):
	env GOOS=darwin GOARCH=amd64 go build -v -o bin/$(DARWIN) -ldflags="-s -w -X main.version=$(VERSION)" ./cmd/db2db/main.go

clean:
	rm -f $(WINDOWS) $(LINUX) $(DARWIN)