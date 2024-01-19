package cli

import (
	"fmt"
	"os"

	"git.voidnet.tech/kev/easysandbox/virtproviders"
)

func ListProviders() {
	providers, getProvidersErr := virtproviders.GetAvailableVirtProviders()
	if getProvidersErr != nil {
		PrintStderr("Error: " + getProvidersErr.Error())
		os.Exit(1)
	}
	for _, provider := range providers {
		fmt.Println(provider)
	}

}
