package domains

import (
	"fmt"
	"time"

	"github.com/estebangarcia21/subprocess"
	"libvirt.org/go/libvirt"

	"errors"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetBackingFile(path string) (string, error) {
	cmd := exec.Command("qemu-img", "info", path)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "backing file: ") {
			backingFilePath := strings.TrimPrefix(line, "backing file: ")
			return filepath.Base(backingFilePath), nil
		}
	}

	return "", errors.New("no backing file found")
}

func StartDomain(name string) error {

	rootCloneFile := DomainsDir + name + "/" + "root.qcow2"
	homeFile := DomainsDir + name + "/" + "home.qcow2"

	fmt.Println("getting backing file for " + rootCloneFile)
	// We create a root qcow2 with a backing file of the root template
	// If it already exists, it is overwritten
	rootToUse, getBackingFileErr := GetBackingFile(rootCloneFile)
	if getBackingFileErr != nil {
		return getBackingFileErr
	}
	fmt.Println("Using root template: " + rootToUse)
	createDomainRootErr := createBackingFile(rootToUse, name, true)

	if createDomainRootErr != nil {
		return createDomainRootErr
	}

	// kind of a hack, because we could use the libvirt api, but that involves XML
	virtInstallCmd := subprocess.New("virt-install", subprocess.Args("--os-variant", "fedora38", "--virt-type=kvm",
		"--name="+name, "--ram", "6000", "--vcpus=6", "--virt-type=kvm", "--hvm", "--network=nat,type=user",
		"--disk", rootCloneFile+",target.bus=sata", "--disk", homeFile+",target.bus=sata", "--import", "--install",
		"no_install=yes", "--noreboot"))

	virtInstallCmdErr := virtInstallCmd.Exec()

	if virtInstallCmdErr != nil {
		return virtInstallCmdErr
	}

	conn, err := libvirt.NewConnect("qemu:///session")
	if err != nil {
		return err
	}
	defer conn.Close()

	domain, err := conn.LookupDomainByName(name)
	if err != nil {
		return err
	}
	defer domain.Free()
	domain.Create()

	return nil
}

func StopDomain(name string) error {
	conn, err := libvirt.NewConnect("qemu:///session")
	if err != nil {
		return err
	}
	defer conn.Close()

	domain, err := conn.LookupDomainByName(name)
	if err != nil {
		return err
	}
	defer domain.Free()

	domain.Shutdown()


	// loop until domain is shut off
	for {
		state, returnInt, stateErr := domain.GetState()
		if returnInt == -1  {
			return errors.New("error getting domain state")
		}
		if stateErr != nil {
			return stateErr
		}
		time.Sleep(1 * time.Second)
		if state == libvirt.DOMAIN_SHUTOFF {
			break
		}
	}
	return domain.Undefine()

}
