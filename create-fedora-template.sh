#!/bin/bash

if [ ! -f xpra.repo ]; then
    wget https://xpra.org/repos/Fedora/xpra.repo
fi

virt-builder fedora-39 --size 7G --format qcow2 --output template.qcow2 \
--copy-in xpra.repo:/etc/yum.repos.d/ \
--update \
--install xpra \
--firstboot-command "useradd -d /home/user user" \
--append-line /etc/fstab:"/dev/sdb1 /home/user/ ext4 defaults 0 1"