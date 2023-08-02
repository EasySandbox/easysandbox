package domains

import (
	"errors"
	"fmt"
	"os"

	"libvirt.org/go/libvirt"
)

func DeleteDomain(name string) error {

	conn, err := libvirt.NewConnect("qemu:///session")
	if err != nil {
		return fmt.Errorf("failed to connect to libvirt: %v", err)
	}

	domain, domainLookupErr := conn.LookupDomainByName(name)
	if domainLookupErr == nil {
		err = domain.Destroy()
		if err != nil && !errors.As(err, &libvirt.Error{Code: libvirt.ERR_OPERATION_INVALID}) {
			return fmt.Errorf("failed to destroy domain: %w", err)
		}
		if err := domain.Undefine(); err != nil {
			return fmt.Errorf("failed to undefine domain: %w", err)
		}
	}
	if err := os.RemoveAll(DomainsDir + "/" + name); err != nil {
		return fmt.Errorf("failed to delete domain directory/files: %w", err)
	}
	return nil
}
