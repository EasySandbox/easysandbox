package cli

import (
	"os"
	"fmt"
	"git.voidnet.tech/kev/easysandbox/domains"
)

func DeleteDomainCmd() {
	if len(os.Args) != 3 {
		FatalStderr("usage: easysandbox delete-domain <domain-name>", 2)
	}

	domainName := os.Args[2]

	err := domains.DeleteDomain(domainName)
	if err != nil {
		FatalStderr("Failed to delete domain: "+err.Error(), 1)
	} else {
		fmt.Println("Domain deleted successfully")
	}
}