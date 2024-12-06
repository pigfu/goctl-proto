# 使用示例

文件api/service.api为以下示例所用到的api文件，工作目录为项目根目录。

## 常规使用

```go
goctl-proto proto --input ./example/api/service.api --output ./example/proto
或者
goctl api plugin -plugin goctl-proto="proto" -api ./example/api/service.api -dir ./example/proto
```

[点击此处](https://github.com/liferod/goctl-proto/blob/main/example/proto/1_normal.proto)查看生成的proto文件。

## 仅生成某个或某几个rpc

```
goctl-proto proto --input ./example/api/service.api --output ./example/proto --inc CreateMock --inc GetMock
或者
goctl api plugin -plugin goctl-proto="proto --inc CreateMock --inc GetMock" -api ./example/api/service.api -dir ./example/proto
```

[点击此处](https://github.com/liferod/goctl-proto/blob/main/example/proto/2_include.proto)查看生成的proto文件。

## 排除某个或某几个rpc

```
goctl-proto proto --input ./example/api/service.api --output ./example/proto --exc Ping --exc GetMock
或者
goctl api plugin -plugin goctl-proto="proto --exc Ping --exc GetMock" -api ./example/api/service.api -dir ./example/proto
```

[点击此处](https://github.com/liferod/goctl-proto/blob/main/example/proto/3_exclude.proto)查看生成的proto文件。
