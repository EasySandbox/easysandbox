package templates

import (
	"os"
)

func getDiskFilesInDir(dir string) (names []string, err error) {
	file, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	names, err = file.Readdirnames(0)
	var filteredNames []string
	if err != nil {
		return nil, err
	}
	// only show qcow2 files
	for i := 0; i < len(names); i++ {
		if names[i][len(names[i])-6:] == ".qcow2" {
			filteredNames = append(filteredNames, names[i])
		}
	}


	return filteredNames, nil
}



func GetRootTemplatePaths() (paths []string, err error) {
	// get template files in TemplateDir

	files, err := getDiskFilesInDir(RootTemplateDir)

	return files, err

}

func GetHomeTemplatePaths() (paths []string, err error) {
	// get template files in TemplateDir

	files, err := getDiskFilesInDir(HomeTemplateDir)

	return files, err

}
