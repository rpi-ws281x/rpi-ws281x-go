#!/bin/bash

docker buildx build --no-cache --platform linux/arm/v7 -t rpi-ws281x-builder-armv7 --load --file docker/app-builder/Dockerfile .
docker buildx build --no-cache --platform linux/arm64  -t rpi-ws281x-builder-arm64 --load --file docker/app-builder/Dockerfile .

for i in examples/*; do
    app=$(basename $i)
    docker run --platform linux/arm/v7 --rm -v "$(pwd)/$i":"/usr/src/$app" -w /usr/src/$app rpi-ws281x-builder-armv7 go build -o "$app-armv7" -v
    docker run --platform linux/arm64 --rm -v "$(pwd)/$i":"/usr/src/$app"  -w /usr/src/$app rpi-ws281x-builder-arm64 go build -o "$app-arm64" -v
done

file examples/*/*-armv7
file examples/*/*-arm64
