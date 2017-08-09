FROM       alpine:3.4
MAINTAINER Adobe Ethos Dev Team <Ethos_Dev@adobe.com>


# add dependencies
RUN apk add --no-cache \
      bash \
      curl \
      openssh-client 

ADD zk-cli-linux-amd64 /usr/bin/zk-cli

ENTRYPOINT ["/usr/bin/zk-cli"]
CMD        []
