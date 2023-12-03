program = ./go-tetris
arg =

.PHONY: all
all: check build

.PHONY: check
check:
	go vet main/*.go
	go vet model/*.go

.PHONY: build
build:
	go build -o go-tetris main/main.go

.PHONY: clean
clean:
	rm -f $(program)

.PHONY: run
run:
	go run main/main.go
