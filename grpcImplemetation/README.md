# grpcDrm

Using grpc and protobuf3 messages for communication between server and client code.

Run Server then Client:

#grpcDrmServer < drm.conf
#grpcDrmClient

## Output of client code. CDI_COMMAND, etc are mentioned just for example.
```
#env CDI_COMMAND=VERSION ./goDrmCdi < drm.conf
{
        "cdiVersion": "0.0.2",
            "supportedVersions": [
                    "0.0.1",
                            "0.0.2"
                                ]
}
#env CDI_COMMAND=INFO ./goDrmCdi < drm.conf
{
        "cdiVersion": "0.0.2",
            "gpu": 2,
                "devices": [
                        "0:0:2.0",
                                "0:1:0.0"
                                    ]
}
#env CDI_VERSION=0.0.1 CDI_COMMAND=ADD CDI_REQUEST=gpu:1 CDI_REQUEST_ID=1234 ./goDrmCdi < drm.conf
{
        "cdiVersion": "0.0.1",
            "devices": [
                    "/dev/dri/card0",
                            "/dev/dri/renderD128"
                                ]
}
#env CDI_VERSION=0.0.1 CDI_COMMAND=ADD CDI_REQUEST=gpu:2 CDI_REQUEST_ID=1234 ./goDrmCdi < drm.conf
{
        "cdiVersion": "0.0.1",
            "code": 102,
                "msg": "Resource unavailable",
                    "details": "No of GPU available: 0"
}
#env CDI_VERSION=0.0.1 CDI_COMMAND=DEL CDI_REQUEST_ID=1234 ./goDrmCdi < drm.conf
{
        "cdiVersion": "0.0.1"
}
#env CDI_VERSION=0.0.1 CDI_COMMAND=ADD CDI_REQUEST=gpu:1,gpu-memory=2048Mi CDI_REQUEST_ID=1234 ./goDrmCdi < drm.conf
{
        "cdiVersion": "0.0.1",
            "code": 4,
                "msg": "Resource spec is not supported",
                    "details": "Unsupported resource request: gpu:1, gpu-memory=2048Mi"
}
#env CDI_VERSION=0.0.1 CDI_COMMAND=DEL CDI_REQUEST_ID=3456 ./goDrmCdi < drm.conf
{
        "cdiVersion": "0.0.1",
            "code": 5,
                "msg": "Unknown container ID"
}
#env CDI_VERSION=0.999 CDI_COMMAND=ADD CDI_REQUEST=gpu:1 CDI_REQUEST_ID=1234 ./goDrmCdi < drm.conf
{
        "cdiVersion": "0.0.2",
            "code": 1,
                "msg": "CDI Version of env-var not supported",
                    "details": "Unsupported version: 0.999"
}
#env CDI_VERSION=0.0.1 CDI_COMMAND=ADD CDI_REQUEST=gpu:1 ./goDrmCdi < drm.conf
{
        "cdiVersion": "0.0.1",
            "code": 103,
                "msg": "Container ID not specified"
}
#env CDI_VERSION=0.0.1 CDI_COMMAND=MYCMD ./goDrmCdi < drm.conf
{
        "cdiVersion": "0.0.1",
            "code": 3,
                "msg": "Command not supported",
                    "details": "Unsupported command MYCMD"
}
#env CDI_VERSION=0.0.1 CDI_REQUEST=gpu:1 ./goDrmCdi < drm.conf
{
        "cdiVersion": "0.0.1",
            "code": 100,
                "msg": "Command not specified"
}
```
