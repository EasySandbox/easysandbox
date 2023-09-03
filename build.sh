#!/bin/sh


cp -Rf ../easysandbox-libvirt/ temp/
go build -buildmode=plugin -o ../plugins/libvirt.so temp/easysandbox-libvirt/main.go
make