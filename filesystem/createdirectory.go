package filesystem

import (
	"os"
)

func CreateDirectory(path string) (err error) {
	return os.MkdirAll(path, 0700)
}
