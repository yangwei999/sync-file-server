FROM openeuler/openeuler:23.03 as BUILDER
RUN dnf update -y && \
    dnf install -y golang && \
    go env -w GOPROXY=https://goproxy.cn,direct

MAINTAINER zengchen1024<chenzeng765@gmail.com>

# build binary
COPY . /go/src/github.com/opensourceways/repo-file-server
RUN cd /go/src/github.com/opensourceways/repo-file-server && CGO_ENABLED=1 go build -v -o ./repo-file-server ./server

# copy binary config and utils
FROM openeuler/openeuler:22.03
RUN dnf -y update && \
    dnf in -y shadow && \
    groupadd -g 1000 repo-file-server && \
    useradd -u 1000 -g repo-file-server -s /bin/bash -m repo-file-server &&\
    mkdir -p /opt/app/

USER repo-file-server

# overwrite config yaml
COPY  --chown=repo-file-server --from=BUILDER /go/src/github.com/opensourceways/repo-file-server/repo-file-server /opt/app

WORKDIR /opt/app/
ENTRYPOINT ["/opt/app/repo-file-server"]
