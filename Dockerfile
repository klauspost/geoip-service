FROM golang:1.3-onbuild
MAINTAINER Klaus Post <klauspost@gmail.com>

EXPOSE 5000

CMD ["app", "-db=/data/geodb.mmdb"]
VOLUME /data/geodb.mmdb