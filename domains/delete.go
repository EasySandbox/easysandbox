package domains

import (
	"errors"
	"os"

	"libvirt.org/go/libvirt"
)

func DeleteDomain(name string) error {

	conn, err := libvirt.NewConnect("qemu:///session")
	if err != nil {
		return err
	}
	defer conn.Close()

	domain, domainLookupErr := conn.LookupDomainByName(name)
	if domainLookupErr == nil {
		defer domain.Free() // if this line is placed out of the if scope the func breaks
		err = domain.Destroy()
		if err != nil && !errors.As(err, &libvirt.Error{Code: libvirt.ERR_OPERATION_INVALID}) {
			return err
		}
	}

	err = os.RemoveAll(DomainsDir + "/" + name)
	if err != nil {
		return err
	}
	if domainLookupErr == nil{
		return domain.Undefine()
	}
	return nil
}