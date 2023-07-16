# Makefile

# Binary names
EASYSANDBOX_BIN := release/easysandbox
IPMAPPER_BIN := release/ipmapper
IPREPORTER_BIN := release/ipreporter

# Go build command
GO_BUILD := go build

# Go source files
MAIN_SRC := main.go
IPMAPPER_SRC := ipmapper/main.go
IPREPORTER_SRC := ipreporter/main.go

.PHONY: all easysandbox ipmapper ipreporter clean

all: easysandbox ipmapper ipreporter

easysandbox:
	$(GO_BUILD) -o $(EASYSANDBOX_BIN) $(MAIN_SRC)

ipmapper:
	$(GO_BUILD) -o $(IPMAPPER_BIN) $(IPMAPPER_SRC)

ipreporter:
	$(GO_BUILD) -o $(IPREPORTER_BIN) $(IPREPORTER_SRC)

clean:
	rm -rf release
