package cli

import (
	"errors"
	"fmt"
	"os"

	"git.voidnet.tech/kev/easysandbox/sandbox"
	"git.voidnet.tech/kev/easysandbox/virtproviders"
)

func DoCreateDomain() {
	if len(os.Args) < 5 {
		FatalStderr("usage: easysandbox create-sandbox <home-template-name> <root-template-name> <sandbox-name>", 4)
	}

	homeTemplate := os.Args[2]
	rootTemplate := os.Args[3]
	name := os.Args[4]

	sandboxInfo := *sandbox.NewSandboxInfo(name, rootTemplate, homeTemplate, 4096, 6)

	provider, getProviderError := virtproviders.GetProviderFromCLIFlag()

	if errors.Is(getProviderError, &virtproviders.VirtProviderLoadFailureError{}) {
		FatalStderr("Failed to load virt provider: " + errors.Unwrap(getProviderError).Error(), 5)
	}

	sandboxAPISymbolString := "CreateSandbox"
	sandboxAPISymbol, err := provider.Lookup(sandboxAPISymbolString)
	if err != nil {
		FatalStderr(fmt.Sprintf("Error looking up sandbox %s API: %s", sandboxAPISymbolString, err), 1)
	}

	sandboxCreationError := sandboxAPISymbol.(func(sandbox.SandboxInfo) error)(sandboxInfo)

	if sandboxCreationError != nil {
		FatalStderr("Failed to create sandbox: "+err.Error(), 5)
	} else {
		fmt.Println("Sandbox created successfully")
	}
}
