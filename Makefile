program = ./go-tetris
arg =

.PHONY: all
all: lint build

.PHONY: lint
lint:
	go vet main.go
	go vet model/*.go

.PHONY: build
build:
	go build -o $(program)

.PHONY: clean
clean:
	rm -f $(program)

.PHONY: run
run:
	go run main.go

.PHONY: test
test:
	go test -v ./model