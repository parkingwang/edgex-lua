# EdgeX LUA - lua脚本节点

EdgeX 节点实现，可以调用Lua脚本来执行远程指令。在本示例中，通过Lua脚本调用ffmpeg截图，并通过curl上传到内部系统。

## Docker Image

基于alpine底层系统的镜像，支持以下三个CPU架构：

1. arm32v7
1. arm64v8
1. amd64

> registry.cn-shenzhen.aliyuncs.com/edge-x/edgenode-lua:0.4.1 