package util

import (
	"bytes"
	"fmt"
	"path/filepath"
)

func ConvertToAbsolutePath(relativePath string) string {
	p, err := filepath.Abs(relativePath)
	if err != nil {
		fmt.Println(err)
	}
	return p
}

func StringConcat(strs ...string) string {
	buf := new(bytes.Buffer)
	for i := 0; i < len(strs); i++ {
		buf.WriteString(strs[i])
	}
	return buf.String()
}
