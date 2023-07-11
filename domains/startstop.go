package domains

import (
	"fmt"
	"time"

	"github.com/estebangarcia21/subprocess"
	"libvirt.org/go/libvirt"

	"errors"
)

func configureDomainRoot(rootPath string, domain string) error {
	var err error
	pass, err := generateRandomPassword(12)
	if err != nil {
		return errors.New("Failed to generate random password: " + err.Error())
	}
	fmt.Println("Password for root: " + pass)
	err = subprocess.New(
		"virt-customize",
		subprocess.Args(
			"-a", rootPath,
			"--no-selinux-relabel",
			"--root-password", "password:"+pass,
			"--delete", "/etc/ssh/*_key",
			"--delete", "/etc/ssh/*.pub",
			"--hostname", domain,
		),
	).Exec()

	if err != nil {
		return errors.New("Failed to configure domain root: " + err.Error())
	}

	return nil
}

func StartDomain(name string, virtInstallArgs subprocess.Option) error {

	rootCloneFile := DomainsDir + name + "/" + "root.qcow2"

	// We create a root qcow2 with a backing file of the root template
	// If it already exists, it is overwritten
	rootToUse, getBackingFileErr := GetBackingFile(rootCloneFile)
	if getBackingFileErr != nil {
		return getBackingFileErr
	}

	createDomainRootErr := createBackingFile(rootToUse, name, true)
	if createDomainRootErr != nil {
		return createDomainRootErr
	}


	// kind of a hack, because we could use the libvirt api, but that involves XML
	virtInstallCmd := subprocess.New("virt-install", virtInstallArgs)

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

	configureDomainRootErr := configureDomainRoot(rootCloneFile, name)
	if configureDomainRootErr != nil {
		return configureDomainRootErr
	}

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
		if returnInt == -1 {
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
