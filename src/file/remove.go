package file

import "common"

func (kp *Keeper) RemoveFile(args *common.FileArgs, reply *common.FileReply) error {
	filename := args.Filename
	err := common.RemoveFile(filename)
	if err != nil {
		reply.Err = err
		return err
	}
	return nil
}
