package domains

import (
	_ "embed"
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
		return fmt.Errorf("faiiled to get host IP: %w", err)
	}

	ipmapperSystemdServiceFile, err := os.CreateTemp("", "ipmapper.service")

	ipmapperSystemdServiceFile.Write(ipmapperSystemdServiceData)

	defer func() {
		os.Remove(ipmapperSystemdServiceFile.Name())
	}()

	if err != nil {
		return fmt.Errorf("failed to setup systemd service file for VM: %w", err)
	}

	ipmapperSystemdTimerFile, err := os.CreateTemp("", "ipmapper.timer")

	ipmapperSystemdTimerFile.Write(ipmapperSystemdTimerData)

	defer func() {
		os.Remove(ipmapperSystemdTimerFile.Name())
	}()

	if err != nil {
		return fmt.Errorf("failed to setup systemd timer file for VM: %w", err)
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
		return fmt.Errorf("failed to configure domain root: %w", err.Error())
	}

	fmt.Println("Password for root: " + pass)

	return nil
}

func StartDomain(name string, virtInstallArgs subprocess.Option) error {

	conn, err := libvirt.NewConnect("qemu:///session")
	if err != nil {
		return fmt.Errorf("error connecting to libvirt: %w", err)
	}

	domain, lookupError := conn.LookupDomainByName(name)
	if lookupError == nil {
		isActive, isActiveErr := domain.IsActive()
		if isActiveErr != nil {
			return fmt.Errorf("error checking if domain is active: %w", isActiveErr)
		}
		if isActive {
			return &DomainIsRunningError{
				Domain: name,
				Msg:    "cannot start domain that is already running",
			}
		}
	}
	rootCloneFile := DomainsDir + name + "/root.qcow2"

	// We create a root qcow2 with a backing file of the root template
	// If it already exists, it is overwritten
	rootToUse, getBackingFileErr := GetBackingFilePath(rootCloneFile)
	if getBackingFileErr != nil {
		return fmt.Errorf("error getting backing file: %w", getBackingFileErr)
	}

	if createDomainRootErr := createBackingFile(rootToUse, name, true); createDomainRootErr != nil {
		return fmt.Errorf("error creating domain root: %w", createDomainRootErr)
	}

	if configureDomainRootErr := configureEphemeralRoot(rootCloneFile, name); configureDomainRootErr != nil {
		return fmt.Errorf("error configuring domain root: %w", configureDomainRootErr)
	}

	// kind of a hack, because we could use the libvirt api, but that involves XML
	if virtInstallCmdErr := subprocess.New("virt-install", virtInstallArgs).Exec(); virtInstallCmdErr != nil {
		return fmt.Errorf("error running virt-install: %w", virtInstallCmdErr)
	}

	if domainCreationError := domain.Create(); domainCreationError != nil {
		return fmt.Errorf("error creating libvirt domain: %w", domainCreationError)
	}

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

	var shutdownAttemptTime = time.Now().Unix()
	for shutdownState != libvirt.DOMAIN_SHUTOFF {
		shutdownState, _, shutdownStateErr = domain.GetState()

		if shutdownStateErr != nil {
			return fmt.Errorf("error getting domain state: %w", shutdownStateErr)
		}
		time.Sleep(50 * time.Millisecond)

		if time.Now().Unix()-shutdownAttemptTime > 5 {
			if shutdownErr := domain.Shutdown(); shutdownErr != nil {
				return fmt.Errorf("Error shutting down libvirt domain %s: %w", name, shutdownErr)
			}
			shutdownAttemptTime = time.Now().Unix()
		}

	}

	if domainUndefineErr := domain.Undefine(); domainUndefineErr != nil {
		return fmt.Errorf("Error deleting libvirt domain (undefining) %s: %w", name, domainUndefineErr)
	}
	return nil
}
