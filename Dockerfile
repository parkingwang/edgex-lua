FROM alpine:3.7

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.cloud.tencent.com/g' /etc/apk/repositories \
    && apk add curl ffmpeg \
    && rm -rf /var/cache/apk/*

# 执行文件名称，须与 name.txt 中一致
COPY edgenode-lua /bin/

ONBUILD COPY application.toml /etc/edgex/

CMD ["edgenode-lua"]