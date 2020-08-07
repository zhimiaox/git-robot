FROM golang:1.14-alpine AS build

LABEL maintainer=""

WORKDIR /build

ADD . /build

ENV GO111MODULE=on \
    GOPROXY="https://goproxy.cn" \
    GOARCH=amd64 \
    CGO_ENABLED=1 \
    GOOS=linux 

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories \
&& apk add --no-cache git gcc musl-dev protoc libprotoc protobuf libprotobuf protobuf-dev \
&& go get -u github.com/golang/protobuf/proto \
&& go get -u github.com/golang/protobuf/protoc-gen-go \
&& export PATH="/root/go/bin:$PATH"

RUN go build -tags musl -o gitrobot_server_linux_amd64 /build
RUN rm -rf data/work-git && mkdir -p data/work-git && cd data/work-git \
&& git init && git checkout -b robot

FROM alpine:3.12 AS prod

COPY --from=build /build/docker-entrypoint.sh /usr/local/bin/
COPY --from=build /build/gitrobot_server_linux_amd64 /usr/local/bin/
COPY --from=build /build/config.toml /zhimiao/
COPY --from=build /build/html /zhimiao/html
COPY --from=build /build/data/work-git /zhimiao/data/work-git

RUN set -eux \
    && chmod a+x /usr/local/bin/docker-entrypoint.sh \
    && addgroup -S -g 1000 zhimiao \
    && adduser -S -G zhimiao -u 1000 zhimiao \
    && chown zhimiao:zhimiao /zhimiao \
    && sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories \
    && apk upgrade -U -a --no-cache \
    && apk add --no-cache \
    'su-exec>=0.2'

EXPOSE 1319

WORKDIR /zhimiao

ENTRYPOINT ["docker-entrypoint.sh"]

CMD ["gitrobot_server_linux_amd64"]