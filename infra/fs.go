package infra

import (
	"fmt"
	"io/fs"
	"os"
)

func HandlePermissionOrFileSystemError(err error) {
	if os.IsPermission(err) {
		// if no permission: quit and ask for permission
		fmt.Println("Require root")
		os.Exit(1)
	} else {
		// if file system error: panic quit
		fmt.Println("File system error")
		os.Exit(1)
	}
}

func MkdirIfNotExist(dir string, perm fs.FileMode) {
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			mkdirErr := os.MkdirAll(dir, perm)
			if mkdirErr != nil {
				HandlePermissionOrFileSystemError(mkdirErr)
			}
		} else {
			HandlePermissionOrFileSystemError(err)
		}
	}
}
