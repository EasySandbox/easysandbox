package cli

import (
	"fmt"
	"os"
	//"git.voidnet.tech/kev/easysandbox/cli"
	"git.voidnet.tech/kev/easysandbox/templates"
	"git.voidnet.tech/kev/easysandbox/domains"
)

func PrintDomainsList(){
	// print all directories in domains.DomainsDir
	files, err := os.ReadDir(domains.DomainsDir)

	if err != nil {
	  fmt.Println("Error reading domains directory:", err)
	  return
	}

	for _, f := range files {
	  if f.IsDir() {
		fmt.Println(f.Name())
	  }
	}
}


func PrintTemplatesList() {
	templateList, err := templates.GetRootTemplatePaths()
	if err != nil && err.(*os.PathError).Err.Error() == "no such file or directory" {
		FatalStderr("Template directory does not exist: "+templates.RootTemplateDir, 1)
	} else if err != nil {
		FatalStderr("Failed to get templates: "+err.Error(), 1)
	}

	if len(templateList) == 0 {
		fmt.Println("No root templates found.")
	} else {
		fmt.Println("Root templates:")
		for _, template := range templateList {
			fmt.Println(template)
		}
	}

	templateList, err = templates.GetHomeTemplatePaths()
	if err != nil && err.(*os.PathError).Err.Error() == "no such file or directory" {
		FatalStderr("Template directory does not exist: "+templates.HomeTemplateDir, 1)
	} else if err != nil {
		FatalStderr("Failed to get templates: "+err.Error(), 1)
	}

	if len(templateList) == 0 {
		fmt.Println("No home templates found.")
	} else {
		fmt.Println("\nHome templates:")
		for _, template := range templateList {
			fmt.Println(template)
		}
	}
}