# EasySandbox

EasySandbox is a Linux tool for managing sandboxes for linux programs.

Most functionality is implemented in libraries that implement sandboxing,
with such that wrap existing software useful for linux sandboxing such
as qemu (via libvirt), bubblewrap, or others.

This program merely directs the plugins according to user intent.

# Building

1. Clone this repo
2. Clone [easysandbox-libvirt](https://github.com/easysandbox/easysandbox-libvirt) to same top level directory
3. cd into easysandbox
4. run ./build.sh

# Creating Templates

Will add instructions for this soon, and later tools to assist with it.