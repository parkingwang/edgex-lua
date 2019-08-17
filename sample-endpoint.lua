---
--- 通过Lua执行一个截图指令,并返回截图结果文件
---

function endpoint_serve(args, frameStr)
    commands = { "avconv",
                 "-i", "rtsp://USER:PASSWORD@camera0.edge.irain.io/11",
                 "-t", "0.001",
                 "-f", "image2",
                 "-vframes", "1",
                 "/tmp/IRAIN_EDGE_CAPTURE.png"
    }
    cmd = table.concat(commands, " ")
    print("正在执行命令: ", cmd)
    ret = os.execute(cmd)
    if 0 == ret then
        return "/tmp/IRAIN_EDGE_CAPTURE.png", nil
    else
        return nil, "执行命令错误[:" .. ret .. "], " .. cmd
    end
end