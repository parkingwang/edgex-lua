package main

import (
	"github.com/nextabc-lab/edgex-go"
	"github.com/yoojia/go-value"
	"github.com/yuin/gopher-lua"
)

//
// Author: 陈哈哈 yoojiachen@gmail.com
// 使用Lua脚本引擎客户端作为Endpoint，接收gRPC控制指令，并返回执行结果

func main() {
	edgex.Run(func(ctx edgex.Context) error {
		config := ctx.LoadConfig()
		name := value.Of(config["Name"]).String()
		rpcAddress := value.Of(config["RpcAddress"]).String()
		scriptFile := value.Of(config["Script"]).String()

		endpoint := ctx.NewEndpoint(edgex.EndpointOptions{
			Name:    name,
			RpcAddr: rpcAddress,
		})

		script := lua.NewState(lua.Options{
			CallStackSize: 8,
			RegistrySize:  8,
		})

		if err := script.DoFile(scriptFile); nil != err {
			ctx.Log().Panic("Lua load script failed: ", err)
		}else{
			ctx.Log().Debug("Load script: ", scriptFile)
		}

		endpoint.Serve(func(in edgex.Message) (out edgex.Message) {
			// 先函数，后参数，正序入栈:
			script.Push(script.GetGlobal("endpointMain"))
			// Arg 1
			script.Push(lua.LString(string(in.Body())))
			// Call
			if err := script.PCall(1, 2, nil); nil != err {
				return edgex.NewMessageString(name, "EX=ERR:"+err.Error())
			} else {
				retData := script.ToString(1)
				retErr := script.ToString(2)
				script.Pop(2)
				if "" != retErr {
					return edgex.NewMessageString(name, "EX=ERR:"+retErr)
				} else {
					return edgex.NewMessageString(name, retData)
				}
			}
		})

		endpoint.Startup()
		defer endpoint.Shutdown()

		return ctx.TermAwait()
	})
}
