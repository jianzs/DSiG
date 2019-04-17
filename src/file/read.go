package file

import "common"

func (kp *Keeper) ReadFile(args *common.FileArgs, reply *common.FileReply) error {
	filename := args.Filename
	data, err := common.ReadFile(filename)
	if err != nil {
		reply.Err = err
		return err
	}
	reply.Content = data
	return nil
}
