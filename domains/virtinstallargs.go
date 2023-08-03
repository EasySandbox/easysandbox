package domains

import (
	"fmt"
	"runtime"

	"github.com/estebangarcia21/subprocess"
)

func GetVirtInstallArgs(domainName string, args ...string) subprocess.Option {
	return subprocess.Args(GetVirtInstallArgsString(domainName, args...)...)
}

func GetVirtInstallArgsString(domainName string, args ...string) []string {
	rootCloneFile := DomainsDir + domainName + "/" + "root.qcow2"
	homeFile := DomainsDir + domainName + "/" + "home.qcow2"

	mandatoryArgs := []string{
		"--name", domainName,
		"--disk", rootCloneFile + ",target.bus=sata",
		"--disk", homeFile + ",target.bus=sata",
		"--import",
		"--hvm",
		"--network", "bridge=virbr0",
		"--virt-type", "kvm",
		"--install", "no_install=yes",
		"--noreboot",
	}
	defaultArgs := map[string]string{
		"--memory":     "4096",
		"--vcpus":      fmt.Sprintf("%d", runtime.NumCPU()),
		"--os-variant": "linux2022",
	}

	var overriddenArgs []string
	for i := 0; i < len(args); i += 2 {
		arg := args[i]
		value := args[i+1]
		delete(defaultArgs, arg)
		overriddenArgs = append(overriddenArgs, arg, value)
	}

	allArgs := append(mandatoryArgs, overriddenArgs...)
	for arg, value := range defaultArgs {
		allArgs = append(allArgs, arg, value)
	}

	return allArgs
}

