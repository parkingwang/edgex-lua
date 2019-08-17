ARG IMAGE=alpine:3.7
FROM $IMAGE

# 执行文件名称，须与 name.txt 中一致
COPY edgenode-lua /bin/

ONBUILD COPY application.toml /etc/edgex/

CMD ["edgenode-lua"]