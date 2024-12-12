# 使用示例

文件api/service.api和api/service@goctl-proto.api为以下示例所用到的api文件，工作目录为项目根目录。

## 常规使用

```
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

## 使用@goctl-proto标识指定生成某个或某几个rpc

service@goctl-proto.api中的“Ping”和“GetMock”这两个接口的描述中包含"@goctl-proto"，因此在不使用其它参数的情况下只会生成这两个接口的rpc信息。

```
goctl-proto proto --input ./example/api/service@goctl-proto.api --output ./example/proto
或者
goctl api plugin -plugin goctl-proto="proto" -api ./example/api/service@goctl-proto.api -dir ./example/proto
```

[点击此处](https://github.com/liferod/goctl-proto/blob/main/example/proto/4_@goctl-proto.proto)查看生成的proto文件。
