# protoc-gen-gin

根据`proto`文件生成`gin`框架

注意：
* 使用了`embed`特性，需要`go1.16`以上版本，`go mod`注意版本信息
    
安装
```shell
go get -u github.com/x-lambda/protoc-gen-gin
```

使用
```shell
protoc -I ./example/rpc/ \
--go_out ./example/rpc --go_opt=paths=source_relative \
--gin_out ./example/rpc --gin_opt=paths=source_relative ./example/rpc/demo/v0/demo.proto
```

example: https://github.com/x-lambda/protoc-gen-gin/tree/master/example
```shell
example
    |`cmd
        |`help         使用说明
        |`job          job系统
        |`server       api注册模块和server启动
    |`dao              db数据访问层
    |`rpc              api接口定义层
    |`server           server实现层
    |`service          service层
    |`util             工具包
    |-main.go
```