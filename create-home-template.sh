#!/bin/bash

virt-make-fs --format=qcow2 --type=ext4 --size=1G --partition=gpt \
templatedata home.qcow2