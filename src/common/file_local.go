package common

import (
	"constant"
	"io/ioutil"
	"os"
	"strings"
)

func ReadFile(filename string) (string, error) {
	if filename[0] != '/' {
		filename = constant.FILE_PREFIX + filename
	}
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func WriteFile(filename string, content string) error {
	if filename[0] != '/' {
		filename = constant.FILE_PREFIX + filename
	}
	err := ensureParentDir(filename)
	if err != nil {
		return err
	}

	outFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	defer outFile.Close()
	if err != nil {
		return err
	}

	_, err = outFile.Write([]byte(content))
	outFile.Write([]byte{'\n'})
	if err != nil {
		return err
	}
	return nil
}

func AppendFile(filename string, content string) error {
	if filename[0] != '/' {
		filename = constant.FILE_PREFIX + filename
	}
	err := ensureParentDir(filename)
	if err != nil {
		return err
	}

	outFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	defer outFile.Close()
	if err != nil {
		return err
	}

	_, err = outFile.Write([]byte(content))
	outFile.Write([]byte{'\n'})
	if err != nil {
		return err
	}
	return nil
}

// removeFile is a simple wrapper around os.Remove that logs errors.
func RemoveFile(filename string) error {
	if filename[0] != '/' {
		filename = constant.FILE_PREFIX + filename
	}
	err := os.Remove(filename)
	if err == os.ErrNotExist {
		return nil
	} else if err != nil {
		return err
	}
	return nil
}

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 写入文件前，先保证父级目录存在
func ensureParentDir(name string) error {
	dir := getParentDir(name)

	if !Exists(dir) {
		err := ensureParentDir(dir)
		if err != nil {
			return err
		}
		err = os.Mkdir(dir, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func getParentDir(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}
