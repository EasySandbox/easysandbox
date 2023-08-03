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
		return fmt.Errorf("Failed to generate password: %w", err)
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
			"--root-password", "password:"+pass,
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

	conn, err := libvirt.NewConnect("qemu:///session")
	if err != nil {
		return err
	}
	defer conn.Close()

	domain, lookupError := conn.LookupDomainByName(name)
	if lookupError == nil {
		isActive, isActiveErr := domain.IsActive()
		if isActiveErr != nil {
			return isActiveErr
		}
		if isActive {
			return &DomainIsRunningError{
				Msg: "Cannot start domain that is already running",
			}
		}
		domain.Free()
	}
	rootCloneFile := DomainsDir + name + "/" + "root.qcow2"

	// We create a root qcow2 with a backing file of the root template
	// If it already exists, it is overwritten
	rootToUse, getBackingFileErr := GetBackingFilePath(rootCloneFile)
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

	domain, err = conn.LookupDomainByName(name)
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
		return fmt.Errorf("error connecting to libvirt: %w", err)
	}

	var shutdownState libvirt.DomainState
	var shutdownStateErr error

	domain, err := conn.LookupDomainByName(name)
	if err != nil {
		return fmt.Errorf("error looking up domain: %w", err)
	}

	var shutdownErr error

	var shutdownAttemptTime = time.Now().Unix()
	for shutdownState != libvirt.DOMAIN_SHUTOFF {
		shutdownState, _, shutdownStateErr = domain.GetState()

		if shutdownStateErr != nil {
			return fmt.Errorf("error getting domain state: %w", shutdownStateErr)
		}
		time.Sleep(50 * time.Millisecond)

		if time.Now().Unix()-shutdownAttemptTime > 5 {
			shutdownErr = domain.Shutdown()
			if shutdownErr != nil {
				return fmt.Errorf("Error shutting down libvirt domain %s: %w", name, shutdownErr)
			}
			shutdownAttemptTime = time.Now().Unix()
		}

	}
	domainUndefineErr := domain.Undefine()
	if domainUndefineErr != nil {
		return fmt.Errorf("Error deleting libvirt domain (undefining) %s: %w", name, domainUndefineErr)
	}
	return nil
}
