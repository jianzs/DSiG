cd ../
export GOPATH=$PWD/
go run src/main/keeper_bin.go worker &

go run src/main/worker_bin.go <local-ip> <master-ip>