FROM golang:1.4-cross
MAINTAINER Klaus Post <klauspost@gmail.com>

ADD crosscompile.sh /usr/local/bin/crosscompile.sh
RUN chmod +x /usr/local/bin/crosscompile.sh

CMD ["/usr/local/bin/crosscompile.sh"]
