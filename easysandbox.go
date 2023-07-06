package main

import (
	"fmt"
	//"flag"
	"os"
	"os/exec"

	//"libvirt.org/go/libvirt"

	"git.voidnet.tech/kev/easysandbox/cli"
	"git.voidnet.tech/kev/easysandbox/domains"
	"git.voidnet.tech/kev/easysandbox/filesystem"
	"git.voidnet.tech/kev/easysandbox/gui"
	"git.voidnet.tech/kev/easysandbox/templates"
	//"github.com/adrg/xdg"
)

func main() {

	// conn, err := libvirt.NewConnect("qemu:///system")
	// if err != nil {
	// 	cli.FatalStderr("Failed to connect to qemu:///system", 1)
	// }
	// defer conn.Close()

	if _, err := exec.LookPath(gui.XPRA_BIN_NAME); err != nil {
		cli.PrintStderr(fmt.Sprintf("%s not found in path\nIt is required for rendering VM apps", gui.XPRA_BIN_NAME))
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		// TODO show help
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
	case "import-root-template":
		fmt.Println("import-root-template")
	case "import-home-template":
		fmt.Println("import-home-template")

	case "create-domain":
		cli.DoCreateDomain()
	case "start-domain":
		fmt.Println("start-domain")
	default:
		fmt.Println("cmd:", cmd)
	}
	//flag.Parse()
	//fmt.Println("listTemplatesPtr:", *listTemplatesPtr)
}
