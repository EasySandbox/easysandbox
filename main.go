package main

import (
	"fmt"
	"os"
	"os/exec"


	"git.voidnet.tech/kev/easysandbox/cli"
	"git.voidnet.tech/kev/easysandbox/domains"
	"git.voidnet.tech/kev/easysandbox/filesystem"
	"git.voidnet.tech/kev/easysandbox/gui"
	"git.voidnet.tech/kev/easysandbox/templates"
)

func main() {

	if _, err := exec.LookPath(gui.XPRA_BIN_NAME); err != nil {
		cli.PrintStderr(fmt.Sprintf("%s not found in path\nIt is required for rendering VM apps", gui.XPRA_BIN_NAME))
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		cli.PrintStderr("usage: easysandbox <command>")
		os.Exit(2)
	}

	func() {
		directoriesToCreate := []string{
			templates.HomeTemplateDir,
			templates.RootTemplateDir,
			domains.DomainsDir,
		}
		for _, dir := range directoriesToCreate {
			dir := dir // necessary because go is weird
			defer func() {
				if err := filesystem.CreateDirectory(dir); err != nil {
					cli.FatalStderr("Failed to create directory: "+dir, 3)
				}
			}()
		}
	}()
	cmd := os.Args[1]

	switch cmd {
	case "list-templates":
		cli.PrintTemplatesList()
	case "create-domain":
		cli.DoCreateDomain()
	case "start-domain":
		cli.StartDomain()
	case "stop-domain":
		cli.StopDomain()
	case "get-domain-ip":
		cli.PrintIPAddress()
	case "show-template":
		cli.ShowTemplate()
	default:
		fmt.Println("unknown cmd:", cmd)
	}

}
