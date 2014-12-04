FROM golang:1.3-onbuild

COPY *.mmdb .

#RUN curl -O http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.mmdb.gz
#RUN gunzip GeoLite2-City.mmdb.gz

EXPOSE 5000
