package sandbox

import (
	"github.com/adrg/xdg"
)

var SandboxInstallDir = xdg.DataHome + "/easysandbox/sandboxes/"


type SandboxInfo struct {
	Name         string
	RootTemplate string
	HomeTemplate string
	MaxMemory    uint32
	MaxCPUs      uint32
}

func NewSandboxInfo(name, rootTemplate, homeTemplate string, maxMem uint32, maxCPUs uint32) *SandboxInfo {

	return &SandboxInfo{
		Name:         name,
		RootTemplate: rootTemplate,
		HomeTemplate: homeTemplate,
		MaxMemory: maxMem,
		MaxCPUs: maxCPUs,
	}
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
