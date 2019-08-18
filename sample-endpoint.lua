---
--- 通过Lua执行一个截图指令,并返回截图结果
---

function endpoint_serve(vnid, seqid, body)
    print("接收参数, VnId: ", vnid)
    print("接收参数, SeqId: ", seqid)
    print("接收参数, Body: ", body)
    image_path = "/tmp/IRAIN_EDGE_CAPTURE.png"
    cap_cmd = { "avconv",
                 "-i", "rtsp://USER:PASSWORD@camera0.edge.irain.io/11",
                 "-t", "0.001",
                 "-f", "image2",
                 "-vframes", "1",
                 image_path
    }
    cap_exec = table.concat(cap_cmd, " ")
    print("执行截图指令: ", cap_exec)
    cap_state = os.execute(cap_exec)
    if 0 == cap_state then
        upl_cmd = {
            "curl", "http://api.edgex.io:5580/upload/" .. seqid .. "/image",
            "-H", '"X-Token: FOO_BAR"',
            "-F", '"file=@' .. image_path .. '"',
            "-v"
        }
        upl_exec = table.concat(upl_cmd, " ")
        print("执行上传指令: ", upl_exec)
        upl_state = os.execute(upl_exec)
        if 0 == upl_state then
            return "SUCCESS", nil
        else
            return nil, "UPL_EXEC:(" .. upl_exec .. "):" .. upl_state
        end
    else
        return nil, "CAP_EXEC:(" .. cap_exec .. "):" .. cap_state
    end
end