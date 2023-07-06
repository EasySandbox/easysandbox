package templates

import (
	"os"
)

func getFilesInDir(dir string) (names []string, err error) {
	file, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	names, err = file.Readdirnames(0)
	if err != nil {
		return nil, err
	}
	return names, nil
}



func GetRootTemplatePaths() (paths []string, err error) {
	// get template files in TemplateDir

	files, err := getFilesInDir(RootTemplateDir)

	return files, err

}

func GetHomeTemplatePaths() (paths []string, err error) {
	// get template files in TemplateDir

	files, err := getFilesInDir(HomeTemplateDir)

	return files, err

}
