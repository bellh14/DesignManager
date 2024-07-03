.PHONY: build
build:
	@go build -C cmd/ -o ../bin/DesignManager

.PHONY: run
run: build
	@./bin/DesignManager

.PHONY: test
test:
	@go test --cover -v ./...

.PHONY: clean
clean:
	@rd /s /q "bin"
	@rm -f ./bin