.PHONY: format
format:
	@gofumpt -l -w .

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: build
build:
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -C cmd/ -tags netgo -ldflags="-w -s" -o ../bin/DesignManager
	@GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -C cmd/ -tags netgo -ldflags="-w -s" -o ../bin/DesignManager-aarch64

.PHONY: run
run: build
	@./bin/DesignManager

.PHONY: test
test:
	@go test --cover -v ./...

.PHONY: clean
clean:
	@rm -rf ./bin


