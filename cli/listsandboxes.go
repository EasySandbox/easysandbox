package cli

import (
	"fmt"
	"os"

	"git.voidnet.tech/kev/easysandbox/sandbox"
)

func PrintSandboxList() {
	// print all directories in sandbox.SandboxInstallDir
	files, err := os.ReadDir(sandbox.SandboxInstallDir)

	if err != nil {
		fmt.Println("Error reading domains directory:", err)
		return
	}

	for _, f := range files {
		if f.IsDir() {
			fmt.Println(f.Name())
		}
	}
}


