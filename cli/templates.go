package cli

import (
	"errors"
	"fmt"
	"os"
	"plugin"

	"git.voidnet.tech/kev/easysandbox/virtproviders"
)

func doShowTemplates(providerPlugin *plugin.Plugin, symbolName string) error {

	sandboxAPISymbol, err := providerPlugin.Lookup(symbolName)
	if err != nil {
		return fmt.Errorf("error looking up sandbox %s API: %s", symbolName, err)
	}
	templatesList, sandboxGetTemplatesErr := sandboxAPISymbol.(func() ([]string, error))()

	if sandboxGetTemplatesErr != nil {
		return errors.New("Failed to list libvirt templates for sandbox: " + sandboxGetTemplatesErr.Error())
	}

	for _, template := range templatesList {
		fmt.Println(template)
	}

	return nil
}

func PrintTemplates() {

	if len(os.Args) < 2 {
		FatalStderr("usage: easysandbox list-templates", 2)
	}

	provider, getProviderError := virtproviders.GetProviderFromCLIFlag()
	if getProviderError != nil {
		FatalStderr("Failed to load sandbox provider: " + getProviderError.Error(), 5)
	}

	fmt.Println("Home templates:")
	showTemplatesErr := doShowTemplates(provider, "GetHomeTemplatesList")
	if showTemplatesErr != nil {
		FatalStderr("Failed to list home templates: "+showTemplatesErr.Error(), 5)
	}
	fmt.Println("Root templates:")
	showTemplatesErr = doShowTemplates(provider, "GetRootTemplatesList")
	if showTemplatesErr != nil {
		FatalStderr("Failed to list root templates: "+showTemplatesErr.Error(), 5)
	}

}
