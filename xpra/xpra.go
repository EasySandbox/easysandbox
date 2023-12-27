package xpra

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"git.voidnet.tech/kev/easysandbox/sandbox"
	"github.com/estebangarcia21/subprocess"
)

var XPRA_BIN_NAME = "xpra"

func getSandboxXPRAConnectionString(sandboxName string) (string, error) {

	portFileData, portFileErr := os.ReadFile(filepath.Join(sandbox.SandboxInstallDir, sandboxName, "xpra-port"))
	if portFileErr != nil {
		return "", portFileErr
	}

	return fmt.Sprintf("tcp://127.0.0.1:%s", string(portFileData)), nil

}

func isXpraAttached(sandboxName string) (bool, error) {
	connString, getConnStringErr := getSandboxXPRAConnectionString(sandboxName)
	if getConnStringErr != nil {
		return false, getConnStringErr
	}
	attachSemaphorePath := filepath.Join(sandbox.SandboxInstallDir, sandboxName, "xpra-attach")
	semaphoreContents, semaphoreReadErr := os.ReadFile(attachSemaphorePath)
	if semaphoreReadErr != nil {
		if os.IsNotExist(semaphoreReadErr) {
			return false, nil
		}
		return false, semaphoreReadErr
	}
	if string(semaphoreContents) == connString {
		return true, nil
	}
	os.Remove(attachSemaphorePath)
	return false, nil
}

func StartXpraClient(sandboxName string) error {

	xpraAttached, xpraAttachedErr := isXpraAttached(sandboxName)
	if xpraAttachedErr != nil {
		return xpraAttachedErr
	}
	if xpraAttached {
		return errors.New("xpra is already attached")
	}

	connString, getConnStringErr := getSandboxXPRAConnectionString(sandboxName)
	if getConnStringErr != nil {
		return getConnStringErr
	}
	attachSemaphorePath := filepath.Join(sandbox.SandboxInstallDir, sandboxName, "xpra-attach")
	defer os.Remove(attachSemaphorePath)
	xpraAttachSemaphoreErr := os.WriteFile(attachSemaphorePath, []byte(connString), 0644)
	if xpraAttachSemaphoreErr != nil {
		return xpraAttachSemaphoreErr
	}

	return subprocess.New("xpra", subprocess.Args("attach", connString, "--splash=no", "--dpi=100")).Exec()
}

func RunXpraCommand(sandboxName string, args ...string) error {

	xpraAttached, xpraAttachedErr := isXpraAttached(sandboxName)
	if xpraAttachedErr != nil {
		return xpraAttachedErr
	}
	if !xpraAttached {
		return errors.New("xpra is not attached")
	}

	connString, getConnStringErr := getSandboxXPRAConnectionString(sandboxName)
	if getConnStringErr != nil {
		return getConnStringErr
	}

	return subprocess.New("xpra", subprocess.Args("control", connString, "start", strings.Join(args, " "))).Exec()

}
