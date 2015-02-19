#!/bin/sh

docker rmi gocross 
docker build --tag="gocross" .
docker run  --rm -it -v "$(pwd)":/usr/src/myapp -w /usr/src/myapp gocross

mkdir build
cp geoip-linux-amd64-static build/geoip-linux-amd64-static
cp Dockerfile-build build/Dockerfile
cd build
docker rmi geoip-service
docker build --tag="geoip-service" .
docker save geoip-service > ../geoip-service-image.tar
cd ..
bzip2 geoip-service-image.tar
rm -rf build
