build:
	@go build -o bin/db2db ./cmd/db2db
run: build
	@./bin/db2db -s=sqlserver://sa:1234Belma.@localhost:1433?database=fro -t=sqlserver://sa:1234Belma.@localhost:1432?database=to
test:
	@go test -v ./..