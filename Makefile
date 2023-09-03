# Makefile

# Binary names
EASYSANDBOX_BIN := release/easysandbox

# Go build command
GO_BUILD := go build

# Go source files
MAIN_SRC := main.go

.PHONY: all tidy easysandbox  clean

all: tidy easysandbox ipmapper

tidy:
	go mod tidy

easysandbox:
	$(GO_BUILD) -o $(EASYSANDBOX_BIN) $(MAIN_SRC)

clean:
	rm -rf release
