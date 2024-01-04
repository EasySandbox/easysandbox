package cli

import (
	"fmt"
	"os"

	"git.voidnet.tech/kev/easysandbox/virtproviders"
)

func SandboxGUIExecute() error {
	if len(os.Args) < 3 {
		FatalStderr("usage: easysandbox gui-exec <sandbox-name>", 2)
	}

	sandboxName := os.Args[2]
	commandArgs := os.Args[3:]

	provider, getProviderError := virtproviders.GetProviderFromCLIFlag()

	if getProviderError != nil {
		FatalStderr("Failed to load virt provider: "+getProviderError.Error(), 5)
	}
	sandboxAPISymbolString := "GUIExecute"
	sandboxAPISymbol, err := provider.Lookup(sandboxAPISymbolString)
	if err != nil {
		FatalStderr(fmt.Sprintf("Error looking up sandbox %s API: %s", sandboxAPISymbolString, err), 1)
	}
	sandboxGUIExecuteErr := sandboxAPISymbol.(func(string, ...string) error)(sandboxName, commandArgs...)

	if sandboxGUIExecuteErr != nil {
		FatalStderr("Failed to execute sandbox gui: "+sandboxGUIExecuteErr.Error(), 5)
	}
	fmt.Println("Sandbox GUI attached successfully")
	return nil

}

func GUIAttach() error {
	if len(os.Args) < 3 {
		FatalStderr("usage: easysandbox gui-attach <sandbox-name>", 2)
	}

	sandboxName := os.Args[2]

	provider, getProviderError := virtproviders.GetProviderFromCLIFlag()

	if getProviderError != nil {
		FatalStderr("Failed to load virt provider: "+getProviderError.Error(), 5)
	}
	sandboxAPISymbolString := "GUIAttach"
	sandboxAPISymbol, err := provider.Lookup(sandboxAPISymbolString)
	if err != nil {
		FatalStderr(fmt.Sprintf("Error looking up sandbox %s API: %s", sandboxAPISymbolString, err), 1)
	}
	sandboxGUIAttachErr := sandboxAPISymbol.(func(string) error)(sandboxName)

	if sandboxGUIAttachErr != nil {
		FatalStderr("Failed to attach sandbox gui: "+sandboxGUIAttachErr.Error(), 5)
	}
	fmt.Println("Sandbox GUI attached successfully")
	return nil

}
