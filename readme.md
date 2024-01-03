# 脚手架

项目结合 gozero 与 grpc-gateway 搭建, 
在自动生成 grpc 的 http api的基础上提供了利用 gozero 编写自定义 http api 的能力.

关键目录结构
```
├── Dockerfile                  //构建线上镜像的dockerfile
├── dsl                         //放dsl文件
├── etc                         //配置文件
├── go.mod
├── go.sum
├── helm                        //helm chart配置
├── internal
│   ├── api                     // http api 代码目录, 骨架由goctl生成
│   ├── config
│   ├── logic                   // rpc逻辑
│   ├── server
│   └── svc                     //服务上下文, rpc 与 api 共用
│       └── servicecontext.go
├── Makefile                    //构建命令
├── pb
└── simple.go
```

## how it works
程序会开启两个端口, 一个为rpc服务端口(8080), 一个为api服务端口(8888). 
grpc-gateway 生成的 http 服务在启动时被配置为直连本机 rpc 服务.
自定义 http api 服务由开发者自行控制如何链接(或者链接哪些) rpc 服务.

api 与 rpc 服务公用同一个 svc.ServiceContext 实例.

代码生成都通过makefile实现.

## 开发流程
### rpc服务
在proto文件中定义好服务和接口, 执行命令:
```shell
make zrpc
```
在 ./internal/logic 目录中的对应文件中编写逻辑
> 如需暴露 rpc 接口为 http api , 
> 需按照 grpc-gateway 文档中的写法修改proto文件的对应接口, 
> 然后执行:
> ```shell
>  make gw
> ```

### 自定义 http api
在一些特殊情况下, 需要创建自定义 http api.
(eg. 文件上传,而 grpc-gateway 并不支持)

在api文件中定义好服务和接口, 执行命令:
```shell
make zapi
```
在 ./internal/api/logic 目录中的对应文件中编写逻辑

## 参考文档
* [grpc-gateway](https://grpc-ecosystem.github.io/grpc-gateway/)
* [gozero](https://go-zero.dev/docs)
