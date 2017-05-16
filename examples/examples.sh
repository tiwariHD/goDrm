#!/bin/bash

#for static builds
go build -ldflags '-extldflags "-static -lm"' ../goDrm.go
#rm -r ../drmFiles

echo "#cat drm.conf"
cat drm.conf

echo "#env CDI_COMMAND=VERSION ./goDrm < drm.conf" 
env CDI_COMMAND=VERSION ./goDrm < drm.conf 

echo "#env CDI_COMMAND=INFO ./goDrm < drm.conf" 
env CDI_COMMAND=INFO ./goDrm < drm.conf 

echo "#env CDI_VERSION=0.0.1 CDI_COMMAND=ADD CDI_REQUEST=gpu:1 CDI_REQUEST_ID=1234 ./goDrm < drm.conf"
env CDI_VERSION=0.0.1 CDI_COMMAND=ADD CDI_REQUEST=gpu:1 CDI_REQUEST_ID=1234 ./goDrm < drm.conf

echo "#env CDI_VERSION=0.0.1 CDI_COMMAND=ADD CDI_REQUEST=gpu:2 CDI_REQUEST_ID=1234 ./goDrm < drm.conf"
env CDI_VERSION=0.0.1 CDI_COMMAND=ADD CDI_REQUEST=gpu:2 CDI_REQUEST_ID=1234 ./goDrm < drm.conf

echo "#env CDI_VERSION=0.0.1 CDI_COMMAND=DEL CDI_REQUEST_ID=1234 ./goDrm < drm.conf"
env CDI_VERSION=0.0.1 CDI_COMMAND=DEL CDI_REQUEST_ID=1234 ./goDrm < drm.conf

echo "#env CDI_VERSION=0.0.1 CDI_COMMAND=ADD CDI_REQUEST=gpu:1,gpu-memory=2048Mi CDI_REQUEST_ID=1234 ./goDrm < drm.conf"
env CDI_VERSION=0.0.1 CDI_COMMAND=ADD CDI_REQUEST=gpu:1,gpu-memory=2048Mi CDI_REQUEST_ID=1234 ./goDrm < drm.conf

echo "#env CDI_VERSION=0.0.1 CDI_COMMAND=DEL CDI_REQUEST_ID=3456 ./goDrm < drm.conf"
env CDI_VERSION=0.0.1 CDI_COMMAND=DEL CDI_REQUEST_ID=3456 ./goDrm < drm.conf

echo "#env CDI_VERSION=0.999 CDI_COMMAND=ADD CDI_REQUEST=gpu:1 CDI_REQUEST_ID=1234 ./goDrm < drm.conf"
env CDI_VERSION=0.999 CDI_COMMAND=ADD CDI_REQUEST=gpu:1 CDI_REQUEST_ID=1234 ./goDrm < drm.conf

echo "#env CDI_VERSION=0.0.1 CDI_COMMAND=ADD CDI_REQUEST=gpu:1 ./goDrm < drm.conf"
env CDI_VERSION=0.0.1 CDI_COMMAND=ADD CDI_REQUEST=gpu:1 ./goDrm < drm.conf

echo "#env CDI_VERSION=0.0.1 CDI_COMMAND=MYCMD ./goDrm < drm.conf"
env CDI_VERSION=0.0.1 CDI_COMMAND=MYCMD ./goDrm < drm.conf

echo "#env CDI_VERSION=0.0.1 CDI_REQUEST=gpu:1 ./goDrm < drm.conf"
env CDI_VERSION=0.0.1 CDI_REQUEST=gpu:1 ./goDrm < drm.conf

