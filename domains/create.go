package domains

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/estebangarcia21/subprocess"

	"git.voidnet.tech/kev/easysandbox/filesystem"
	"git.voidnet.tech/kev/easysandbox/templates"
)

func createBackingFile(template string, name string, isRoot bool) error {
	// ensure template ends with .qcow2
	if !strings.HasSuffix(template, ".qcow2") {
		template += ".qcow2"
	}
	var templatePath string
	var targetFile string
	if isRoot {
		templatePath = templates.RootTemplateDir + template
		targetFile = DomainsDir + name + "/root.qcow2"
	} else {
		templatePath = templates.HomeTemplateDir + template
		targetFile = DomainsDir + name + "/home.qcow2"
	}

	fmt.Println("Creating domain with backing file: " + templatePath + " at " + targetFile)

	return subprocess.New(
		"qemu-img",
		subprocess.Arg("create"),
		subprocess.Arg("-f"),
		subprocess.Arg("qcow2"),
		subprocess.Arg("-F"),
		subprocess.Arg("qcow2"),
		subprocess.Arg("-b"),
		subprocess.Arg(templatePath),
		subprocess.Arg(targetFile)).Exec()
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
	createDomainRootError := createBackingFile(rootTemplate, name, true)
	if createDomainRootError != nil {
		return createDomainRootError
	}
	//homeCopyError := copy.Copy(templates.HomeTemplateDir+homeTemplate, DomainsDir+name+"/home.qcow2")
	homeCopyError := createBackingFile(homeTemplate, name, false)
	if homeCopyError != nil {
		return homeCopyError
	}
	return nil

}
