package main

import (
	"fmt"
	"os"

	"git.voidnet.tech/kev/easysandbox/cli"
	"git.voidnet.tech/kev/easysandbox/filesystem"
	"git.voidnet.tech/kev/easysandbox/sandbox"
	"git.voidnet.tech/kev/easysandbox/templates"
)

var Version = "0.0.0"

func init() {
	//defer trace.Stop()
}

func main() {

	if len(os.Args) < 2 {
		cli.PrintStderr("usage: easysandbox <command>")
		os.Exit(2)
	}

	func() {
		directoriesToCreate := []string{
			templates.HomeTemplateDir,
			templates.RootTemplateDir,
			sandbox.SandboxInstallDir,
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
	case "list-sandboxs":
		fallthrough
	case "list-sandboxes":
		cli.PrintDomainsList()
	case "create-sandbox":
		cli.DoCreateDomain()
	case "delete-sandbox":
		cli.DeleteSandbox()
	case "start-sandbox":
		cli.StartDomain()
	case "stop-sandbox":
		cli.StopSandbox()
	//case "get-sandbox-ip":
	//cli.PrintIPAddress()
	case "version":
		fmt.Println(Version)
	default:
		fmt.Println("unknown cmd:", cmd)
	}

}
