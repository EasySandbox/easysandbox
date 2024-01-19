package virtproviders

import (
	"fmt"
	"os"
	"plugin"
	"strings"

	"github.com/adrg/xdg"
)

var pluginDirs = []string{"/usr/lib64/easysandbox/", xdg.DataHome + "/easysandbox/plugins/", "release/easysandbox-plugins/", "easysandbox-plugins/"}

type VirtProviderLoadFailureError struct {
	Name string
}

func (e *VirtProviderLoadFailureError) Error() string {
	return fmt.Sprintf("virt provider %s not found", e.Name)
}

func GetAvailableVirtProviders() ([]string, error) {
	var availableProviders []string

	for _, dir := range pluginDirs {
		if _, err := os.Stat(dir); err == nil {
			dirEntries, err := os.ReadDir(dir)
			if err != nil {
				return nil, err
			}

			for _, dirEntry := range dirEntries {
				if !dirEntry.IsDir() && strings.HasSuffix(dirEntry.Name(), ".so") {
					providerName := strings.TrimSuffix(dirEntry.Name(), ".so")
					availableProviders = append(availableProviders, providerName)
				}
			}
		}
	}

	if len(availableProviders) == 0 {
		return nil, fmt.Errorf("no virtualization providers found")
	}

	return availableProviders, nil
}

func GetProviderFromCLIFlag() (*plugin.Plugin, error) {
	var providerType string

	for i := 0; i < len(os.Args); i++ {
		if strings.HasPrefix(os.Args[i], "-type=") {
			providerType = os.Args[i]
			providerType = strings.ReplaceAll(providerType, "-type=", "")
			break
		}
	}

	if providerType == "" {
		providerType = "libvirt"
	}

	var pluginDir string
	for _, dir := range pluginDirs {
		if _, err := os.Stat(dir); err == nil {
			pluginDir = dir
			break
		}
	}
	if pluginDir == "" {
		return nil, fmt.Errorf("no plugin directory found")
	}
	if !strings.HasSuffix(pluginDir, "/") {
		pluginDir += "/"
	}

	fmt.Printf("%s%s.so\n", pluginDir, providerType)
	pluginFilePath := fmt.Sprintf("%s%s.so", pluginDir, providerType)


	provider, providerLoadError := plugin.Open(pluginFilePath)
	if providerLoadError != nil {
		return nil, &VirtProviderLoadFailureError{Name: providerType}
	}
	return provider, nil


}
