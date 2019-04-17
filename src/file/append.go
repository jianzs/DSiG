package file

import "common"

func (kp *Keeper) AppendFile(args *common.FileArgs, reply *common.FileReply) error {
	filename := args.Filename
	data := args.Content
	err := common.AppendFile(filename, data)
	if err != nil {
		reply.Err = err
		return err
	}
	reply.Content = data
	return nil
}
