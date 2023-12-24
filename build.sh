#!/bin/sh
rm ~/.local/share/easysandbox/sandboxes/sandboxchristmas/root-prepared;
cd ../easysandbox-libvirt
go build -buildmode=plugin -o ../easysandbox/plugins/libvirt.so main.go
cd ../easysandbox
make