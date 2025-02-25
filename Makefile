build:
	go build -o setup-tools main.go

run:
	./setup-tools

init:
	./setup-tools init

install-zsh:
	./setup-tools install-zsh

clean:
	rm -f setup-tools

all: build run
