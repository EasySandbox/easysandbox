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
		cli.PrintSandboxList()
	case "create-sandbox":
		cli.DoCreateSandbox()
	case "delete-sandbox":
		cli.DeleteSandbox()
	case "start-sandbox":
		cli.StartSandbox()
	case "stop-sandbox":
		cli.StopSandbox()
	case "attach-gui":
		cli.GUIAttach()
	case "sandbox-exec":
		fallthrough
	case "sandbox-execute":
		cli.SandboxExecute()
	case "version":
		fmt.Println(Version)
	default:
		fmt.Println("unknown cmd:", cmd)
	}

}
