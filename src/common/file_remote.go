package common

// 通用的文件读写操作

import (
	"constant"
)

func ReadFileRemote(addr string, filename string) (string, error) {
	args := FileArgs{Filename: filename}
	var reply FileReply

	err := Call(SrvAddr(addr, constant.FILE_RPC), constant.READ_FILE, args, &reply)
	if err != nil {
		return "", err
	}

	if reply.Err != nil {
		return "", reply.Err
	}

	return reply.Content, nil
}

func WriteFileRemote(addr string, filename string, content string) error {
	args := FileArgs{filename, content}
	var reply FileReply

	err := Call(SrvAddr(addr, constant.FILE_RPC), constant.WRITE_FILE, args, &reply)
	if err != nil {
		return err
	}

	if reply.Err != nil {
		return reply.Err
	}

	return nil
}

func AppendFileRemote(addr, filename, content string) error {
	args := FileArgs{filename, content}
	var reply FileReply

	err := Call(SrvAddr(addr, constant.FILE_RPC), constant.APPEND_FILE, args, &reply)
	if err != nil {
		return err
	}

	if reply.Err != nil {
		return reply.Err
	}
	return nil
}
