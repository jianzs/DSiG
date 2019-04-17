package file

import "common"

func (kp *Keeper) WriteFile(args *common.FileArgs, reply *common.FileReply) error {
	filename := args.Filename
	data := args.Content
	err := common.WriteFile(filename, data)
	if err != nil {
		reply.Err = err
		return err
	}
	reply.Content = data
	return nil
}
