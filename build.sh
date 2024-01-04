#!/bin/sh
cd ../easysandbox-libvirt
go build -buildmode=plugin -o ../easysandbox/release/plugins/libvirt.so main.go
cd ../easysandbox
make