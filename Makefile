program = ./go-tetris
arg =

.PHONY: all
all: lint build

.PHONY: lint
lint:
	go vet main/*.go
	go vet model/*.go

.PHONY: build
build:
	go build -o $(program) main/main.go

.PHONY: clean
clean:
	rm -f $(program)

.PHONY: run
run:
	go run main/main.go
