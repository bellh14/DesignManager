build:
	@go build -C cmd/ -o ../bin/DesignManager

run: build
	@./bin/DesignManager

test:
	@go test -v ./...