# EasySandbox

EasySandbox is a Linux tool for managing sandboxes for linux programs.

Most functionality is implemented in libraries that implement sandboxing,
with such that wrap existing software useful for linux sandboxing such
as qemu (via libvirt), bubblewrap, or others. **Currently, only libvirt is available**

"We have Qubes at home": this project is meant to be distro agnostic and reduce the need to live in the Qubes ecosystem. This project is not as secure as Qubes, but the author dogfoods it and believes it migate most malware.

Current features:

* libvirt VM based sandboxes
* seamless GUI windowing for the VM (via xpra)
* user home directory templates
* root file system templates
* qcow2 backed images, aka linked clones for root images
* ephemeral roots - changes to VM root are lost after reboot

Planned features:

* podman based sandbox
* colored GUI windows like in Qubes
* file mounts
* network filtering
* GUI sandbox manager
* Qubes-type clipboard
* Permission system for microphone, webcam, attaching USB

# Usage and install

1. Install dependencies: libvirt, qemu, xpra, virtinstall, [guestfs-tools](https://libguestfs.org/guestfs-faq.1.html#binaries)
2. Download the release
3. Create your first home and root templates. (see below)
4. Create your first sandbox: ./easysandbox create-sandbox home_template_name root_template_name sandbox_name
5. Run sandbox: ./easysandbox start-sandbox sandbox_name
6. Attach GUI: ./easysandbox attach-gui sandbox_name
7. Run GUI program: ./easysandbox gui-exec sandbox_name xterm

## Creating Templates


You need a root and home template to get started

A root template (a linux install in a qcow2 file) needs:

* xpra installed (no need to create service). recommended to use official xpra package (not distro repo)
* a user named user (password doesn't matter)
* setting to mount home disk at boot: /etc/fstab: "/dev/sdb1 /home/user/ ext4 defaults 0 1"
* place it in ~/.local/share/easysandbox/root-templates/


A home template needs:

* ext4 partition in qcow2 file. put your favorite dotfiles here
* place it in ~/.local/share/easysandbox/home-templates/



# Building

1. Clone this repo
2. Clone [easysandbox-libvirt](https://github.com/easysandbox/easysandbox-libvirt) to same top level directory
3. cd into easysandbox
4. run ./build.sh
