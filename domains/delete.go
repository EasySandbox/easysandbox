package domains

import(
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
		if err != nil {
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