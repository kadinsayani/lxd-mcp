.PHONY: all clean build

all: clean build

clean:

build:
	go build -o lxd-mcp-server .
