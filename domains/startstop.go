package domains

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/estebangarcia21/subprocess"
	"libvirt.org/go/libvirt"

	"git.voidnet.tech/kev/easysandbox/getourip"
)

var ipmapperClientPath = "./release/ipmapper"

//go:embed templatedata/ipmapper.service
var ipmapperSystemdServiceData []byte

//go:embed templatedata/ipmapper.timer
var ipmapperSystemdTimerData []byte

func configureEphemeralRoot(rootPath string, domain string) error {
	var err error
	pass, err := generateRandomPassword(12)

	if err != nil {
		return errors.New("Failed to generate random password: " + err.Error())
	}

	hostIP, err := getourip.GetOurIP()
	if err != nil {
		return err
	}

	ipmapperSystemdServiceFile, err := os.CreateTemp("", "ipmapper.service")

	ipmapperSystemdServiceFile.Write(ipmapperSystemdServiceData)

	defer func() {
		os.Remove(ipmapperSystemdServiceFile.Name())
	}()

	if err != nil {
		return errors.New("Failed to setup systemd service file for VM: " + err.Error())
	}

	ipmapperSystemdTimerFile, err := os.CreateTemp("", "ipmapper.timer")

	ipmapperSystemdTimerFile.Write(ipmapperSystemdTimerData)

	defer func() {
		os.Remove(ipmapperSystemdTimerFile.Name())
	}()

	if err != nil {
		return errors.New("Failed to setup systemd timer file for VM: " + err.Error())
	}

	err = subprocess.New(
		"virt-customize",
		subprocess.Args(
			"-a", rootPath,
			"--no-selinux-relabel",
			"--firstboot-command", "setenforce 0",
			"--root-password", "password:a",
			"--upload",
			fmt.Sprintf("%s:%s", ipmapperClientPath, "/bin/ipmapper"),
			"--upload",
			fmt.Sprintf("%s:%s", ipmapperSystemdServiceFile.Name(), "/etc/systemd/system/ipmapper.service"),
			"--upload",
			fmt.Sprintf("%s:%s", ipmapperSystemdTimerFile.Name(), "/etc/systemd/system/ipmapper.timer"),
			"--append-line", "/etc/fstab:/dev/sdb1 /user/ ext4 defaults 0 1",
			"--append-line", fmt.Sprintf("/etc/hosts:%s hostsystem", hostIP),
			"--firstboot-command", "systemctl daemon-reload",
			"--firstboot-command", "systemctl enable ipmapper.timer",
			"--firstboot-command", "systemctl start ipmapper.timer",
			"--delete", "/etc/ssh/*_key",
			"--delete", "/etc/ssh/*.pub",
			"--hostname", domain,
		),
	).Exec()

	if err != nil {
		return errors.New("Failed to configure domain root: " + err.Error())
	}

	fmt.Println("Password for root: " + pass)

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

	configureDomainRootErr := configureEphemeralRoot(rootCloneFile, name)
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
	shutdownSignalTimeout := time.After(10 * time.Second)
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
		select {
		case <-shutdownSignalTimeout:
			domain.Shutdown()
			shutdownSignalTimeout = time.After(10 * time.Second)
		default:
			continue
		}
	}
	return domain.Undefine()

}
