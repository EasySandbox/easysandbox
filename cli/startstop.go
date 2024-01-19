package cli

import (
	"fmt"
	"os"

	"git.voidnet.tech/kev/easysandbox/virtproviders"
)

func StartSandbox() error {
	if len(os.Args) < 3 {
		FatalStderr("usage: easysandbox start <sandbox-name>", 2)
	}

	sandboxName := os.Args[2]

	provider, getProviderError := virtproviders.GetProviderFromCLIFlag()

	if getProviderError != nil {
		FatalStderr("Failed to load virt provider: "+getProviderError.Error(), 5)
	}
	sandboxAPISymbolString := "StartSandbox"
	sandboxAPISymbol, err := provider.Lookup(sandboxAPISymbolString)
	if err != nil {
		FatalStderr(fmt.Sprintf("Error looking up sandbox %s API: %s", sandboxAPISymbolString, err), 1)
	}
	sandboxStartError := sandboxAPISymbol.(func(string) error)(sandboxName)

	if sandboxStartError != nil {
		FatalStderr("Failed to start sandbox: "+sandboxStartError.Error(), 5)
	}
	fmt.Println("Sandbox started successfully")
	return nil

}

func StopSandbox() error {
	if len(os.Args) < 3 {
		FatalStderr("usage: easysandbox stop <sandbox-name>", 2)
	}

	sandboxName := os.Args[2]

	provider, getProviderError := virtproviders.GetProviderFromCLIFlag()

	if getProviderError != nil {
		FatalStderr("Failed to load virt provider: "+getProviderError.Error(), 5)
	}
	sandboxAPISymbolString := "StopSandbox"
	sandboxAPISymbol, err := provider.Lookup(sandboxAPISymbolString)
	if err != nil {
		FatalStderr(fmt.Sprintf("Error looking up sandbox %s API: %s", sandboxAPISymbolString, err), 1)
	}
	sandboxStopError := sandboxAPISymbol.(func(string) error)(sandboxName)

	if sandboxStopError != nil {
		FatalStderr("Failed to stop sandbox: "+sandboxStopError.Error(), 5)
	}
	fmt.Println("Sandbox stopped successfully")
	return nil

}
