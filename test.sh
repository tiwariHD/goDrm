#!/bin/bash

#for static builds
go install -ldflags '-extldflags "-static -lm"'

echo "#cat drm.conf"
cat drm.conf

echo "#env CDI_VERSION=0.0.1 CDI_COMMAND=ADD CDI_REQUEST=gpu:1 CDI_CONTAINERID=1234 ../../../../bin/goDrm < drm.conf"
env CDI_VERSION=0.0.1 CDI_COMMAND=ADD CDI_REQUEST=gpu:1 CDI_CONTAINERID=1234 ../../../../bin/goDrm < drm.conf

echo "#env CDI_VERSION=0.0.1 CDI_COMMAND=ADD CDI_REQUEST=gpu:2 CDI_CONTAINERID=1234 ../../../../bin/goDrm < drm.conf"
env CDI_VERSION=0.0.1 CDI_COMMAND=ADD CDI_REQUEST=gpu:2 CDI_CONTAINERID=1234 ../../../../bin/goDrm < drm.conf

echo "#env CDI_VERSION=0.0.1 CDI_COMMAND=DEL CDI_CONTAINERID=1234 ../../../../bin/goDrm < drm.conf"
env CDI_VERSION=0.0.1 CDI_COMMAND=DEL CDI_CONTAINERID=1234 ../../../../bin/goDrm < drm.conf

echo "#env CDI_VERSION=0.0.1 CDI_COMMAND=ADD CDI_REQUEST=gpu:1,gpu-memory=2048Mi CDI_CONTAINERID=1234 ../../../../bin/goDrm < drm.conf"
env CDI_VERSION=0.0.1 CDI_COMMAND=ADD CDI_REQUEST=gpu:1,gpu-memory=2048Mi CDI_CONTAINERID=1234 ../../../../bin/goDrm < drm.conf

echo "#env CDI_VERSION=0.0.1 CDI_COMMAND=DEL CDI_CONTAINERID=3456 ../../../../bin/goDrm < drm.conf"
env CDI_VERSION=0.0.1 CDI_COMMAND=DEL CDI_CONTAINERID=3456 ../../../../bin/goDrm < drm.conf

echo "#env CDI_VERSION=0.999 CDI_COMMAND=ADD CDI_REQUEST=gpu:1 CDI_CONTAINERID=1234 ../../../../bin/goDrm < drm.conf"
env CDI_VERSION=0.999 CDI_COMMAND=ADD CDI_REQUEST=gpu:1 CDI_CONTAINERID=1234 ../../../../bin/goDrm < drm.conf

echo "#env CDI_VERSION=0.0.1 CDI_COMMAND=ADD CDI_REQUEST=gpu:1 ../../../../bin/goDrm < drm.conf"
env CDI_VERSION=0.0.1 CDI_COMMAND=ADD CDI_REQUEST=gpu:1 ../../../../bin/goDrm < drm.conf

echo "#env CDI_VERSION=0.0.1 CDI_COMMAND=MYCMD ../../../../bin/goDrm < drm.conf"
env CDI_VERSION=0.0.1 CDI_COMMAND=MYCMD ../../../../bin/goDrm < drm.conf

echo "#env CDI_VERSION=0.0.1 CDI_REQUEST=gpu:1 ../../../../bin/goDrm < drm.conf"
env CDI_VERSION=0.0.1 CDI_REQUEST=gpu:1 ../../../../bin/goDrm < drm.conf

