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
		boardId := value.Of(config["BoardId"]).String()
		majorId := value.Of(config["MajorId"]).String()
		devType := value.Of(config["DeviceType"]).String()
		if "" == boardId || "" == majorId || "" == devType {
			log.Panic("未设置参数：BoardId/MajorId/DeviceType")
		}
		if "" == scriptFile {
			log.Panic("未设置LuaScript文件")
		}

		ctx.InitialWithConfig(config)
		endpoint := ctx.NewEndpoint(edgex.EndpointOptions{
			NodePropertiesFunc: FuncEndpointProperties(boardId, majorId, devType),
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

		endpoint.Serve(func(in edgex.Message) (out []byte) {
			unionId := in.UnionId()
			eventId := in.EventId()
			log.Debugf("接收到RPC控制指令: UnionId= %s, EventId= %d", unionId, eventId)
			// 先函数，后参数，正序入栈:
			script.Push(script.GetGlobal("endpoint_serve"))
			// 三个参数： unionId, EventId, Body
			script.Push(lua.LString(unionId))
			script.Push(lua.LNumber(eventId))
			script.Push(lua.LString(string(in.Body())))
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
func FuncEndpointProperties(boardId, majorId, devType string) func() edgex.MainNodeProperties {
	return func() edgex.MainNodeProperties {
		return edgex.MainNodeProperties{
			NodeType:   edgex.NodeTypeEndpoint,
			Vendor:     "EdgeX",
			ConnDriver: "Script/LUA",
			VirtualNodes: []*edgex.VirtualNodeProperties{
				{
					BoardId:     boardId,
					MajorId:     majorId,
					MinorId:     "LUA",
					DeviceType:  devType,
					Description: fmt.Sprintf("%s/%s/LUA", boardId, majorId),
					Virtual:     false,
					StateCommands: map[string]string{
						"TRIGGER": "AT+NOP",
					},
				},
			},
		}
	}
}
