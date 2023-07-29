package cli

import (
	"flag"
	"os"

	"git.voidnet.tech/kev/easysandbox/gui"
)

func DomainExec() error {

	var ipmapperSocket string

	if len(os.Args) < 4 {
		FatalStderr("usage: easysandbox domain-exec <domain-name> <program> <args...>", 2)
		os.Exit(1)
	}

	domain := os.Args[2]
	prog := os.Args[3]
	args := os.Args[4:]

	flag.StringVar(&ipmapperSocket, "ipmapper-address", "127.0.0.1:8080", "IPMapper socket address")
	flag.Parse()

	runErr := gui.RunGUIApplication(ipmapperSocket, domain, prog, args...)
	if runErr != nil {
		FatalStderr("Failed to run GUI application: "+runErr.Error(), 1)
	}
	return nil

}
