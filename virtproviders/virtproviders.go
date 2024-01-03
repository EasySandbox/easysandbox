package virtproviders

import (
	"flag"
	"fmt"
	"plugin"
)

var pluginDir = "plugins"

type VirtProviderLoadFailureError struct {
	Name string
}

func (e *VirtProviderLoadFailureError) Error() string {
	return fmt.Sprintf("virt provider %s not found", e.Name)
}

func GetProviderFromCLIFlag() (*plugin.Plugin, error) {

	var providerType string
	flag.StringVar(&providerType, "type", "libvirt", "virtualization provider type")
	flag.Parse()

	fmt.Println("Attempting to load")
	fmt.Println(fmt.Sprintf("%s/%s.so", pluginDir, providerType))
	provider, providerLoadError := plugin.Open(fmt.Sprintf("%s/%s.so", pluginDir, providerType))
	if providerLoadError != nil {
		return nil, &VirtProviderLoadFailureError{Name: providerType}
	}
	fmt.Println("loaded")
	return provider, nil

}
