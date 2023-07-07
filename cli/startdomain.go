package cli

import (
	"fmt"
	"os"

	"git.voidnet.tech/kev/easysandbox/domains"
)

func StartDomain() error {
	if len(os.Args) != 3 {
		FatalStderr("usage: easysandbox start <domain-name>", 2)
	}

	domainName := os.Args[2]


	err := domains.StartDomain(domainName)

	if err != nil {
		FatalStderr("Failed to start domain: "+domainName + "\n" + err.Error(), 1)
	}
	fmt.Println("Started domain: " + domainName)
	return nil

}
