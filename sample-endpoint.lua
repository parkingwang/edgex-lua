---
--- 通过Lua执行一个截图指令,并返回截图结果
---

function startup()
    print(">> LUA Startup")
    os.execute("ls /tmp")
end

function shutdown()
    print(">> LUA Shutdown")
end

function endpoint_serve(vnid, seqid, body)
    image_path = "/tmp/camera-capture"..seqid..".jpg"
    cap_cmd = { "ffmpeg",
                 "-i", "rtsp://admin:pass@192.168.1.4",
                 "-t", "0.001",
                 "-vframes", "1",
                 "-loglevel", "error",
                 image_path
    }
    cap_exec = table.concat(cap_cmd, " ")
    print("执行截图指令: ", cap_exec)
    cap_state = os.execute(cap_exec)
    if 0 == cap_state then
        upl_cmd = {
            "curl", "http://api.edgex.domain/upload/" .. seqid .. "/image",
            "-H", '"Form-Type: file"',
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