build:
	@go build -o bin/DesignManager

run: build
	@./bin/DesignManager

test:
	@go test -v ./...