#!/bin/bash

#for static builds
go build -ldflags '-extldflags "-static -lm"' ../goDrmCdi.go
#rm -r ../drmFiles

echo "#cat drm.conf"
cat drm.conf

echo "#env CDI_COMMAND=VERSION ./goDrmCdi < drm.conf" 
env CDI_COMMAND=VERSION ./goDrmCdi < drm.conf 

echo "#env CDI_COMMAND=INFO ./goDrmCdi < drm.conf" 
env CDI_COMMAND=INFO ./goDrmCdi < drm.conf 

echo "#env CDI_VERSION=0.0.1 CDI_COMMAND=ADD CDI_REQUEST=gpu:1 CDI_REQUEST_ID=1234 ./goDrmCdi < drm.conf"
env CDI_VERSION=0.0.1 CDI_COMMAND=ADD CDI_REQUEST=gpu:1 CDI_REQUEST_ID=1234 ./goDrmCdi < drm.conf

echo "#env CDI_VERSION=0.0.1 CDI_COMMAND=ADD CDI_REQUEST=gpu:2 CDI_REQUEST_ID=1234 ./goDrmCdi < drm.conf"
env CDI_VERSION=0.0.1 CDI_COMMAND=ADD CDI_REQUEST=gpu:2 CDI_REQUEST_ID=1234 ./goDrmCdi < drm.conf

echo "#env CDI_VERSION=0.0.1 CDI_COMMAND=DEL CDI_REQUEST_ID=1234 ./goDrmCdi < drm.conf"
env CDI_VERSION=0.0.1 CDI_COMMAND=DEL CDI_REQUEST_ID=1234 ./goDrmCdi < drm.conf

echo "#env CDI_VERSION=0.0.1 CDI_COMMAND=ADD CDI_REQUEST=gpu:1,gpu-memory=2048Mi CDI_REQUEST_ID=1234 ./goDrmCdi < drm.conf"
env CDI_VERSION=0.0.1 CDI_COMMAND=ADD CDI_REQUEST=gpu:1,gpu-memory=2048Mi CDI_REQUEST_ID=1234 ./goDrmCdi < drm.conf

echo "#env CDI_VERSION=0.0.1 CDI_COMMAND=DEL CDI_REQUEST_ID=3456 ./goDrmCdi < drm.conf"
env CDI_VERSION=0.0.1 CDI_COMMAND=DEL CDI_REQUEST_ID=3456 ./goDrmCdi < drm.conf

echo "#env CDI_VERSION=0.999 CDI_COMMAND=ADD CDI_REQUEST=gpu:1 CDI_REQUEST_ID=1234 ./goDrmCdi < drm.conf"
env CDI_VERSION=0.999 CDI_COMMAND=ADD CDI_REQUEST=gpu:1 CDI_REQUEST_ID=1234 ./goDrmCdi < drm.conf

echo "#env CDI_VERSION=0.0.1 CDI_COMMAND=ADD CDI_REQUEST=gpu:1 ./goDrmCdi < drm.conf"
env CDI_VERSION=0.0.1 CDI_COMMAND=ADD CDI_REQUEST=gpu:1 ./goDrmCdi < drm.conf

echo "#env CDI_VERSION=0.0.1 CDI_COMMAND=MYCMD ./goDrmCdi < drm.conf"
env CDI_VERSION=0.0.1 CDI_COMMAND=MYCMD ./goDrmCdi < drm.conf

echo "#env CDI_VERSION=0.0.1 CDI_REQUEST=gpu:1 ./goDrmCdi < drm.conf"
env CDI_VERSION=0.0.1 CDI_REQUEST=gpu:1 ./goDrmCdi < drm.conf

