package cli

import (
	"fmt"
	"os"

	"git.voidnet.tech/kev/easysandbox/sandbox"
	"git.voidnet.tech/kev/easysandbox/virtproviders"
)

func DoCreateSandbox() {
	if len(os.Args) < 5 {
		FatalStderr("usage: easysandbox create-sandbox <home-template-name> <root-template-name> <sandbox-name>", 4)
	}

	homeTemplate := os.Args[2]
	rootTemplate := os.Args[3]
	name := os.Args[4]

	sandboxInfo, invalidSandboxInfoErr := sandbox.NewSandboxInfo(name, rootTemplate, homeTemplate, 4096, 6)
	if invalidSandboxInfoErr != nil {
		FatalStderr("Invalid sandbox info: "+invalidSandboxInfoErr.Error(), 4)
	}

	provider, getProviderError := virtproviders.GetProviderFromCLIFlag()

	if getProviderError != nil {
		FatalStderr("Failed to load virt provider: "+getProviderError.Error(), 5)
	}

	sandboxAPISymbolString := "CreateSandbox"
	sandboxAPISymbol, err := provider.Lookup(sandboxAPISymbolString)
	if err != nil {
		FatalStderr(fmt.Sprintf("Error looking up sandbox %s API: %s", sandboxAPISymbolString, err), 1)
	}

	sandboxCreationError := sandboxAPISymbol.(func(sandbox.SandboxInfo) error)(*sandboxInfo)

	if sandboxCreationError != nil {
		FatalStderr("Failed to create sandbox: "+err.Error(), 5)
	} else {
		fmt.Println("Sandbox created successfully")
	}
}
