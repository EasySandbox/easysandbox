package gui

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"

	//"github.com/estebangarcia21/subprocess"

	"git.voidnet.tech/kev/easysandbox/domains"

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

	retryClient.RetryMax = 10
	notFoundTries := 10
	for {
		domainIPRes, getDomainIPError = retryClient.Get(
			fmt.Sprint("http://", ipmapperSocket, "/get?key=", domain))

		if getDomainIPError != nil {
			return getDomainIPError
		}
		if domainIPRes.StatusCode == 404 {
			notFoundTries--
			if notFoundTries == 0 {
				return fmt.Errorf("IPMapper returned 404 10 times in a row")
			}
			time.Sleep(2 * time.Second)
		} else if domainIPRes.StatusCode > 299 {
			return fmt.Errorf("IPMapper returned status code %d", domainIPRes.StatusCode)
		} else {
			break
		}
	}
	domainIPBytes, bodyErr := io.ReadAll(domainIPRes.Body)
	if bodyErr != nil {
		return bodyErr
	}


	// todo convert this to subprocess
	cmd := exec.Command(XPRA_BIN_NAME, "start", "--dpi=100", "--ssh=ssh -o StrictHostKeyChecking=accept-new", "ssh:user@"+string(domainIPBytes), "--exit-with-children", "--start-child='"+program+"'")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()
	return cmd.Wait()

}
