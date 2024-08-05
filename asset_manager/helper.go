package asset_manager

import (
	"os"
)

func LoadFileOnDisk(filePath string) (bytes []byte, err error) {
	stat, err := os.Stat(filePath)
	if err != nil {
		return
	}
	size := stat.Size()
	bytes = make([]byte, size)
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return
	}
	_, err = file.ReadAt(bytes, 0)
	return
}

// escape " into ""
func sqlEscape(s string) string {
	res := ""
	for _, r := range s {
		if r == '"' {
			res += string(`"`)
		}
		res += string(r)
	}
	return res
}
