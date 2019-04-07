package common

// 通用的文件读写操作

import (
	"constant"
	"io/ioutil"
	"os"
	"strings"
)

func ReadFileFrMaster(mr string, filename string) (string, error) {
	args := ReadFileArgs{Filename: filename}
	var reply ReadFileReply

	err := Call(SrvAddr(mr, constant.MASTER_RPC), constant.READ_FILE, args, &reply)
	if err != nil {
		return "", err
	}

	if reply.Err != nil {
		return "", reply.Err
	}

	return reply.Content, nil
}

func WriteFileToMaster(mr string, filename string, content string) error {
	args := WriteFileArgs{filename, content}
	var reply WriteFileReply

	err := Call(SrvAddr(mr, constant.MASTER_RPC), constant.WRITE_FILE, args, &reply)
	if err != nil {
		return err
	}

	if reply.Err != nil {
		return reply.Err
	}

	return nil
}

func ReadFile(filename string) (string, error) {
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