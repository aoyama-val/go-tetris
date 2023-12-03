package = github.com/aoyama-val/go-tetris
program = ./go-tetris
arg =

.PHONY: all
all: check build

.PHONY: check
check:
	go vet $(package)

.PHONY: build
build:
	go build

.PHONY: clean
clean:
	rm -f $(program)

.PHONY: run
run:
	go run main.go model.go
