package cli

import (
	"fmt"
	"os"
	//"git.voidnet.tech/kev/easysandbox/cli"
	"git.voidnet.tech/kev/easysandbox/domains"
)

func PrintIPAddress() {

	if len(os.Args) != 3 {
		FatalStderr("usage: easysandbox get-ip <domain>", 1)
	}

	ip, err := domains.GetDomainIP(os.Args[2])
	if err != nil {
		FatalStderr("Failed to get IP address: "+err.Error(), 1)
	}
	fmt.Println(ip)

}
