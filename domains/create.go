package domains

import (
	"errors"
	"os"
	"strings"

	"github.com/estebangarcia21/subprocess"
	"github.com/otiai10/copy"

	"git.voidnet.tech/kev/easysandbox/filesystem"
	"git.voidnet.tech/kev/easysandbox/templates"
)

func createDomainRoot(template string, name string) error {
	// ensure template ends with .qcow2
	if !strings.HasSuffix(template, ".qcow2") {
		template += ".qcow2"
	}
	rootCloneFile := DomainsDir + name + "/" + template

	return subprocess.New(
		"qemu-img",
		subprocess.Arg("create"),
		subprocess.Arg("-f"),
		subprocess.Arg("qcow2"),
		subprocess.Arg("-F"),
		subprocess.Arg("qcow2"),
		subprocess.Arg("-b"),
		subprocess.Arg(templates.RootTemplateDir+template),
		subprocess.Arg(rootCloneFile)).Exec()
}

func CreateDomain(homeTemplate string, rootTemplate string, name string) error {
	existsFunc := func(path string) (bool, error) {
		_, err := os.Stat(path)
		if err == nil {
			return true, nil
		}
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	// ensure domain that would be created doesn't already exist
	exists, existsErr := existsFunc(DomainsDir + name + "/home.qcow2")
	if existsErr != nil {
		return existsErr
	}
	if exists {
		return errors.New("domain already exists")
	}

	// ensure homeTemplate ends with .qcow2
	if !strings.HasSuffix(homeTemplate, ".qcow2") {
		homeTemplate += ".qcow2"
	}
	directoryCreateError := filesystem.CreateDirectory(DomainsDir + name)
	if directoryCreateError != nil {
		return directoryCreateError
	}
	createDomainRootError := createDomainRoot(rootTemplate, name)
	if createDomainRootError != nil {
		return createDomainRoot(rootTemplate, name)
	}
	homeCopyError := copy.Copy(templates.HomeTemplateDir+homeTemplate, DomainsDir+name+"/home.qcow2")
	if homeCopyError != nil {
		return homeCopyError
	}
	return nil

}
