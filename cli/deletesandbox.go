package cli

import (
	"fmt"
	"os"
	"errors"
	"git.voidnet.tech/kev/easysandbox/virtproviders"

)

func DeleteSandbox()  {
	if len(os.Args) < 3 {
		FatalStderr("usage: easysandbox delete <sandbox-name>", 2)
	}

	sandboxName := os.Args[2]

	provider, getProviderError := virtproviders.GetProviderFromCLIFlag()

	if errors.Is(getProviderError, &virtproviders.VirtProviderLoadFailureError{}) {
		FatalStderr("Failed to load virt provider: " + errors.Unwrap(getProviderError).Error(), 5)
	}
	sandboxAPISymbolString := "DeleteSandbox"
	sandboxAPISymbol, err := provider.Lookup(sandboxAPISymbolString)
	if err != nil {
		FatalStderr(fmt.Sprintf("Error looking up sandbox %s API: %s", sandboxAPISymbolString, err), 1)
	}
	sandboxDeleteError := sandboxAPISymbol.(func(string) error)(sandboxName)

	if sandboxDeleteError != nil {
		FatalStderr("Failed to stop sandbox: "+sandboxDeleteError.Error(), 5)
	}
	fmt.Println("Sandbox deleted successfully")

}