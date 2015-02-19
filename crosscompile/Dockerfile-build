FROM daaku/goruntime
COPY geoip-linux-amd64 /geoip-service

EXPOSE 5000

CMD ["/geoip-service", "-db=/data/geodb.mmdb"]
VOLUME /data/geodb.mmdb
