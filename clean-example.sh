#!/bin/bash

docker rmi rpi-ws281x-builder-armv7
docker rmi rpi-ws281x-builder-arm64

rm examples/*/*-armv7
rm examples/*/*-arm64
