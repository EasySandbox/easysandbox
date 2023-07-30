# Makefile

# Binary names
EASYSANDBOX_BIN := release/easysandbox
IPMAPPER_BIN := release/ipmapper

# Go build command
GO_BUILD := go build

# Go source files
MAIN_SRC := main.go
IPMAPPER_SRC := ipmapper/main.go

.PHONY: all tidy easysandbox ipmapper clean

all: tidy easysandbox ipmapper

tidy:
	go mod tidy

easysandbox:
	$(GO_BUILD) -o $(EASYSANDBOX_BIN) $(MAIN_SRC)

ipmapper:
	$(GO_BUILD) -o $(IPMAPPER_BIN) $(IPMAPPER_SRC)

clean:
	rm -rf release
