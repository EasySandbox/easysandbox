package domains

import (
	"libvirt.org/go/libvirt"
)

func GetDomainIP(name string) (string, error) {
	conn, err := libvirt.NewConnect("qemu:///session")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	domain, err := conn.LookupDomainByName(name)
	if err != nil {
		return "", err
	}
	defer domain.Free()

	interfaces, err := domain.ListAllInterfaceAddresses(libvirt.DOMAIN_INTERFACE_ADDRESSES_SRC_LEASE)
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		for _, addr := range iface.Addrs {
			if addr.Type == libvirt.IP_ADDR_TYPE_IPV4 {
				return addr.Addr, nil
			}
		}

	}

	return "", nil
}
