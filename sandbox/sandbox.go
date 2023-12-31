package sandbox

import (
	"errors"
	"github.com/adrg/xdg"
	"regexp"
)

var SandboxInstallDir = xdg.DataHome + "/easysandbox/sandboxes/"


type SandboxInfo struct {
	Name         string
	RootTemplate string
	HomeTemplate string
	MaxMemory    uint32
	MaxCPUs      uint32
}

func NewSandboxInfo(sandboxName, rootTemplate, homeTemplate string, maxMem uint32, maxCPUs uint32) (*SandboxInfo, error) {
	sandboxNameValid, _ := regexp.MatchString("^[a-zA-Z0-9]{1,12}$", sandboxName)
	if !sandboxNameValid {
		return nil, errors.New("err invalid sandbox name")
	}
	return &SandboxInfo{
		Name:         sandboxName,
		RootTemplate: rootTemplate,
		HomeTemplate: homeTemplate,
		MaxMemory: maxMem,
		MaxCPUs: maxCPUs,
	}, nil
}

type SandboxCreator interface {
	CreateSandbox(SandboxInfo) error
}

type Sandbox interface {
	StartSandbox(SandboxInfo, discardChanges bool) error
	StopSandbox(SandboxInfo, hard bool) error
	RebootSandbox(SandboxInfo, hard bool) error
	DeleteSandbox(SandboxInfo) error
	GetSandboxStatus(SandboxInfo) (string, error)
}

type SandboxExec interface {
	//SSHShell(Sandbox, SSHArgs) error
	//SSHExec(Sandbox) error
	XPRAExec(sandbox SandboxInfo, command string, args ***string) error
}
