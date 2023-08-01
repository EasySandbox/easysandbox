package gui

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"git.voidnet.tech/kev/easysandbox/domains"

	"github.com/estebangarcia21/subprocess"
	"github.com/hashicorp/go-retryablehttp"
)

var XPRA_BIN_NAME = "xpra"

var retryClient = retryablehttp.NewClient()

var DomainIsRunningError *domains.DomainIsRunningError

func RunGUIApplication(ipmapperSocket string,
	domain string, program string, args ...string) error {

	startDomainErr := domains.StartDomain(domain, domains.GetVirtInstallArgs(domain))
	if !errors.As(startDomainErr, &DomainIsRunningError) {
		return startDomainErr
	}

	var domainIPRes *http.Response
	var getDomainIPError error
	gotIPMapperResult := false

	retryClient.RetryMax = 10
	notFoundTries := retryClient.RetryMax
	httpFailTries := retryClient.RetryMax

	for !gotIPMapperResult {
		domainIPRes, getDomainIPError = retryClient.Get(
			fmt.Sprint("http://", ipmapperSocket, "/get?key=", domain))

		if getDomainIPError != nil {
			httpFailTries--
			if httpFailTries == 0 {
				return fmt.Errorf("ipMapper returned error after %d: %w", httpFailTries, getDomainIPError)
			}
		}
		if domainIPRes.StatusCode == 404 {
			notFoundTries--
			if notFoundTries == 0 {
				return fmt.Errorf("IPMapper returned 404 10 times in a row")
			}
			time.Sleep(2 * time.Second)
		} else if domainIPRes.StatusCode > 299 {
			return fmt.Errorf("ipmapper returned status code %d", domainIPRes.StatusCode)
		} else {
			gotIPMapperResult = true
		}
	}
	domainIPBytes, bodyReadErr := io.ReadAll(domainIPRes.Body)
	if bodyReadErr != nil {
		return fmt.Errorf("failed to read IPMapper response body: %w", bodyReadErr)
	}

	xpraCmdArgs := subprocess.Args(
		"start",
		"--dpi=100",
		"--ssh=ssh -o StrictHostKeyChecking=accept-new",
		"ssh:user@"+string(domainIPBytes),
		"--exit-with-children",
		"--start-child='"+program+"'")

	cmdObj := subprocess.New(XPRA_BIN_NAME, xpraCmdArgs)

	return cmdObj.Exec()

}
