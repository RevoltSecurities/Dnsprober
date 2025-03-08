package utils

import (
	"errors"
	"os"
)

func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return !info.IsDir()
}

func CanRead(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsPermission(err) {
			return false, errors.New("read permission denied")
		}
		return false, err
	}
	file.Close()
	return true, nil
}

func CanWrite(path string) (bool, error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		if os.IsPermission(err) {
			return false, errors.New("write permission denied")
		}
		return false, err
	}
	file.Close()
	return true, nil
}

func HasPermission(path string) (bool, error) {
	canRead, errRead := CanRead(path)
	if errRead != nil {
		return false, errRead
	}

	canWrite, errWrite := CanWrite(path)
	if errWrite != nil {
		return false, errWrite
	}

	return canRead && canWrite, nil
}

func CreateTmpFile() (string, error) {
	tmpFile, err := os.CreateTemp(os.TempDir(), "")
	if err != nil {
		return "", err
	}
	name := tmpFile.Name()
	if err := tmpFile.Close(); err != nil {
		return "", err
	}
	return name, nil
}

func CreateTmpDir(prefix string) (string, error) {
	dir, err := os.MkdirTemp(os.TempDir(), prefix)
	if err != nil {
		return "", err
	}
	return dir, nil
}
