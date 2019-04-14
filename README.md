# Distributed System in Go

使用Go语言实现一个简单的MapReduce分布式系统

## 功能
1. Master, Worker 单独启动。
2. 用户可多次提交`Job`。
3. 用户可自定义`Map`和`Reduce`函数。通过改写`main/executor_bin_example.bin`中的`MapFunc`和`ReduceFunc`。

## 缺点
1. 只能在Master节点创建Client提交。
2. 文件都存在Master上。
3. Worker执行任务都是对Master上的文件进行读取和写入。
4. 输入文件只能使用绝对路径

## 计划
1. Worker与Client直接连接。
2. 设计失败反馈。
3. 设计心跳，增强容错。

## 配置
修改存放中间文件及结果的目录前缀：修改`src/common_file.go`

## 启动
### Master
```bash
go run src/main/master_bin.go
```

### Worker
```bash
go run src/main/worker_bin.go <local-ip> <master-ip>
```

### Client
```bash
 go run src/main/client_bin.go <path/to/executor> <absolute/path/to/input/files>
```
executor的文件名省略最后`.go`后缀，否则go会认为是多go文件执行

