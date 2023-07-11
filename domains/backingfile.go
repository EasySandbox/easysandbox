package domains

import (
	"errors"
	"os/exec"
	"path/filepath"
	"strings"
)


func GetBackingFile(path string) (string, error) {
	cmd := exec.Command("qemu-img", "info", path)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "backing file: ") {
			backingFilePath := strings.TrimPrefix(line, "backing file: ")
			return filepath.Base(backingFilePath), nil
		}
	}

	return "", errors.New("no backing file found")
}
