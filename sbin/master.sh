cd ../
export GOPATH=$PWD/
go run src/main/keeper_bin.go master &

go run src/main/master_bin.go