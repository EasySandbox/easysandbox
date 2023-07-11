package cli

import (
	"fmt"
	"os"

	"git.voidnet.tech/kev/easysandbox/domains"
)

func StartDomain() error {
	var userProvidedArgs []string
	if len(os.Args) < 3 {
		FatalStderr("usage: easysandbox start <domain-name>", 2)
	} else if len(os.Args) > 3 {
		userProvidedArgs = os.Args[3:]
	}

	domainName := os.Args[2]

	virtInstallArgs := domains.GetVirtInstallArgs(domainName, userProvidedArgs...)


	err := domains.StartDomain(domainName, virtInstallArgs)

	if err != nil {
		FatalStderr("Failed to start domain: "+domainName + "\n" + err.Error(), 1)
	}
	fmt.Println("Started domain: " + domainName)
	return nil
}

func StopDomain() error {
	if len(os.Args) != 3 {
		FatalStderr("usage: easysandbox stop <domain-name>", 2)
	}

	domainName := os.Args[2]

	err := domains.StopDomain(domainName)

	if err != nil {
		FatalStderr("Failed to stop domain: "+domainName + "\n" + err.Error(), 1)
	}
	fmt.Println("Stopped domain: " + domainName)
	return nil
}
