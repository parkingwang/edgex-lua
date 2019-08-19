package main

import (
	"fmt"
	"github.com/nextabc-lab/edgex-go"
	"github.com/yoojia/go-value"
	"github.com/yuin/gopher-lua"
)

//
// Author: 陈哈哈 yoojiachen@gmail.com
// 使用Lua脚本引擎客户端作为Endpoint，接收RPC控制指令，并返回执行结果

func main() {
	edgex.Run(func(ctx edgex.Context) error {
		config := ctx.LoadConfig()
		log := ctx.Log()
		scriptFile := value.Of(config["Script"]).String()
		majorId := value.Of(config["MajorId"]).String()
		minorId := value.Of(config["MinorId"]).String()
		if "" == majorId || "" == minorId {
			log.Panic("未设置参数：MajorId/MinorId")
		}
		if "" == scriptFile {
			log.Panic("未设置LuaScript文件")
		}

		ctx.InitialWithConfig(config)
		endpoint := ctx.NewEndpoint(edgex.EndpointOptions{
			NodePropertiesFunc: FuncEndpointProperties(majorId, minorId),
		})

		script := lua.NewState(lua.Options{
			CallStackSize: 8,
			RegistrySize:  8,
		})

		log.Debug("加载脚本文件: ", scriptFile)
		if err := script.DoFile(scriptFile); nil != err {
			log.Panic("加载脚本出错: ", err)
		}

		// startup
		script.Push(script.GetGlobal("startup"))
		if err := script.PCall(0, 0, nil); nil != err {
			log.Panic("脚本Startup执行错误：", err)
		}
		// shutdown
		defer func() {
			script.Push(script.GetGlobal("shutdown"))
			if err := script.PCall(0, 0, nil); nil != err {
				log.Error("脚本Shutdown执行错误：", err)
			}
		}()

		endpoint.Serve(func(in edgex.Message) []byte {
			// 先函数，后参数，正序入栈:
			script.Push(script.GetGlobal("endpoint_serve"))
			vnid := in.VirtualNodeId()
			eventId := in.EventId()
			body := string(in.Body())
			// 三个参数： vnId, EventId, Body
			script.Push(lua.LString(vnid))
			script.Push(lua.LNumber(eventId))
			script.Push(lua.LString(body))
			// Call
			if err := script.PCall(3, 2, nil); nil != err {
				return []byte("EX=ERR:" + err.Error())
			} else {
				retData := script.ToString(1)
				retErr := script.ToString(2)
				script.Pop(2)
				if "" != retErr {
					return []byte("EX=ERR:" + retErr)
				} else {
					return []byte("EX=OK:" + retData)
				}
			}
		})

		endpoint.Startup()
		defer endpoint.Shutdown()

		return ctx.TermAwait()
	})
}

// 创建EndpointNode函数
func FuncEndpointProperties(majorId string, minorId string) func() edgex.MainNodeProperties {
	return func() edgex.MainNodeProperties {
		return edgex.MainNodeProperties{
			NodeType:   edgex.NodeTypeEndpoint,
			Vendor:     "EdgeX",
			ConnDriver: "Script/LUA",
			VirtualNodes: []*edgex.VirtualNodeProperties{
				{
					VirtualId:   fmt.Sprintf("LUA-%s-%s", majorId, minorId),
					MajorId:     majorId,
					MinorId:     minorId,
					Description: fmt.Sprintf("%s:%s-脚本驱动", majorId, minorId),
					Virtual:     false,
					StateCommands: map[string]string{
						"TRIGGER": "AT+NOP",
					},
				},
			},
		}
	}
}
