package cli

import (
	"fmt"
	"os"
	"git.voidnet.tech/kev/easysandbox/domains"
)

func DoCreateDomain() {
	if len(os.Args) != 5 {
		FatalStderr("usage: easysandbox create-domain <home-template> <root-template> <name>", 4)
	}
	homeTemplate := os.Args[2]
	rootTemplate := os.Args[3]
	name := os.Args[4]

	err := domains.CreateDomain(homeTemplate, rootTemplate, name)
	if err != nil {
		FatalStderr("Failed to create domain: " +err.Error(), 5)
	} else {
		fmt.Println("Domain created successfully")
	}
}
