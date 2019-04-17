package common

// 通用的文件读写操作

import (
	"constant"
)

func ReadFileRemote(srv, filename string) (string, error) {
	args := FileArgs{Filename: filename}
	var reply FileReply

	err := Call(srv, constant.READ_FILE, args, &reply)
	if err != nil {
		return "", err
	}

	if reply.Err != nil {
		return "", reply.Err
	}

	return reply.Content, nil
}

func WriteFileRemote(srv, filename, content string) error {
	args := FileArgs{filename, content}
	var reply FileReply

	err := Call(srv, constant.WRITE_FILE, args, &reply)
	if err != nil {
		return err
	}

	if reply.Err != nil {
		return reply.Err
	}

	return nil
}

func AppendFileRemote(srv, filename, content string) error {
	args := FileArgs{filename, content}
	var reply FileReply

	err := Call(srv, constant.APPEND_FILE, args, &reply)
	if err != nil {
		return err
	}

	if reply.Err != nil {
		return reply.Err
	}
	return nil
}

func RemoveFileRemote(srv, filename string) error {
	args := FileArgs{Filename:filename}
	var reply FileReply

	err := Call(srv, constant.REMOVE_FILE, args, &reply)
	if err != nil {
		return err
	}

	if reply.Err != nil {
		return reply.Err
	}
	return nil
}
