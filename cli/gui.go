package cli

import (
	"os"

	"git.voidnet.tech/kev/easysandbox/xpra"
)

func GUIAttach() {
	sandboxName := os.Args[2]
	if startXpraErr := xpra.StartXpraClient(sandboxName); startXpraErr != nil {
		FatalStderr("Failed to attach xpra to sandbox: "+startXpraErr.Error(), 5)
	}
}

func SandboxExecute() {
	sandboxName := os.Args[2]
	sandboxCmd := os.Args[3]
	if execSandboxErr := xpra.RunXpraCommand(sandboxName, sandboxCmd); execSandboxErr != nil {
		FatalStderr("Failed to execute command in sandbox: "+execSandboxErr.Error(), 6)
	}
}
