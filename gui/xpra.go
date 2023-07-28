package gui

import (
	"fmt"
	"io"
	"github.com/estebangarcia21/subprocess"
	"github.com/hashicorp/go-retryablehttp"
)

var XPRA_BIN_NAME = "xpra"

var retryClient = retryablehttp.NewClient()

func RunGUIApplication(ipmapperSocket string,
	domain string, program string, args ...string) error {

	retryClient.RetryMax = 5

	domainIPRes, getDomainIPError := retryClient.Get(
		fmt.Sprint("http://", ipmapperSocket, "/?key=", domain))

	if getDomainIPError != nil {
		return getDomainIPError
	}
	if domainIPRes.StatusCode > 299 {
		return fmt.Errorf("IPMapper returned status code %d", domainIPRes.StatusCode)
	}
	domainIPBytes, bodyErr := io.ReadAll(domainIPRes.Body)
	if bodyErr != nil {
		return bodyErr
	}

	subprocessArgsForGUI := subprocess.Args(
		"start",
		"--dpi=100",
		"--ssh=ssh ssh:user@" + string(domainIPBytes),
		"--exit-with-children",
		"--start-child=" + program)

	return subprocess.New(
		XPRA_BIN_NAME, subprocessArgsForGUI, subprocess.Args(args...)).Exec()
}
