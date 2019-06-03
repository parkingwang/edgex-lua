ARG IMAGE=alpine:3.7
FROM $IMAGE

LABEL MAINTAINER="yoojiachen@gmail.com"

# 执行文件名称，须与 name.txt 中一致
COPY endpoint-lua /bin/

ONBUILD COPY application.toml /etc/edgex/

CMD ["endpoint-lua"]