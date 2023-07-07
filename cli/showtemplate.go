package cli

import (
	"fmt"
	"os"

	"git.voidnet.tech/kev/easysandbox/domains"
)

func ShowTemplate() {
	if len(os.Args) != 3 {
		FatalStderr("usage: easysandbox show-template <domain-name>", 2)
	}

	domainName := os.Args[2]

	template, err := domains.GetBackingFile(domains.DomainsDir + "/" + domainName + "/root.qcow2")
	if err != nil {
		FatalStderr("Failed to get backing file: "+err.Error(), 1)
	}
	fmt.Println(template)

}