package main

import (
	"fmt"
	"os"

	"git.voidnet.tech/kev/easysandbox/cli"
	"git.voidnet.tech/kev/easysandbox/sandbox"
)

var Version = "0.0.0"

func main() {

	if len(os.Args) < 2 {
		cli.PrintStderr("usage: easysandbox <command>")
		os.Exit(2)
	}

	createDirectory := func(path string) (err error) {
		return os.MkdirAll(path, 0700)
	}

	func() {
		directoriesToCreate := []string{
			sandbox.SandboxInstallDir,
		}
		for _, dir := range directoriesToCreate {
			dir := dir // necessary because go is weird
			defer func() {
				if err := createDirectory(dir); err != nil {
					cli.FatalStderr("Failed to create directory: "+dir, 3)
				}
			}()
		}
	}()
	cmd := os.Args[1]

	switch cmd {
	case "list-providers":
		cli.ListProviders()
	case "list-templates":
		cli.PrintTemplates()
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
	case "gui-attach":
		cli.GUIAttach()
	case "gui-exec":
		fallthrough
	case "gui-execute":
		cli.SandboxGUIExecute()
	case "version":
		fmt.Println(Version)
	default:
		fmt.Println("unknown cmd:", cmd)
	}

}
