FROM alpine:3.15 AS dms-server

LABEL maintainer="junsheng.wu <junsheng.wu@cetccloud.com>"

ADD bin/ /usr/local/bin/

ENTRYPOINT [ "/usr/local/bin/dms-server" ]

EXPOSE 9200
