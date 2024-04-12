package converter

import (
	"os"
	"strings"
)

// GetFileList ...
func GetFileList(id string, path string) []string {
	var list []string
	files, _ := os.ReadDir(path)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.HasPrefix(file.Name(), id) {
			list = append(list, file.Name())
		}
	}
	return list
}
