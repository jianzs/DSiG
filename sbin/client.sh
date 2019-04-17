export GOPATH=$PWD/../
go run src/main/keeper_bin.go client &

 go run src/main/client_bin.go <local-ip> <master-ip> <absolute/path/to/executor> <absolute/path/to/input/files>