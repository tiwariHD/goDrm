# grpcDrm

Using grpc and protobuf3 messages for communication between server and client code.

Run Server then Client:

#grpcDrmServer -f grpcDrmServer/drm.conf

#grpcDrmClient

## Output of grpcDrmClient.
```
#VERSION: Empty{}
{
    "cdiVersion": "0.0.2",
    "supportedVersions": [
        "0.0.1",
        "0.0.2"
    ]
}

#INFO: Empty{}
{
    "cdiVersion": "0.0.2",
    "gpu": 2,
    "devices": [
        "0:1:0.0",
        "0:0:2.0"
    ]
}

#ADD: AddRequest{Version: "0.0.1", Request: "gpu:1", RequestId: "1234"}
{
    "cdiVersion": "0.0.1",
    "devices": [
        "/dev/dri/card1",
        "/dev/dri/renderD129"
    ]
}

#ADD: AddRequest{Version: "0.0.1", Request: "gpu:2", RequestId: "1234"}
{
    "addError": {
        "cdiVersion": "0.0.1",
        "code": 102,
        "msg": "Resource unavailable",
        "details": "No of GPU available: 0"
    }
}

#ADD: AddRequest{Version: "0.0.1", Request: "gpu:1, gpu-memory=2048Mi", RequestId: "1234"}
{
    "addError": {
        "cdiVersion": "0.0.1",
        "code": 4,
        "msg": "Resource spec is not supported",
        "details": "Unsupported resource request: gpu:1, gpu-memory=2048Mi"
    }
}

#ADD: AddRequest{Version: "0.999", Request: "gpu:1", RequestId: "1234"}
{
    "addError": {
        "cdiVersion": "0.0.2",
        "code": 1,
        "msg": "CDI Version of env-var not supported",
        "details": "Unsupported version: 0.999"
    }
}

#ADD: AddRequest{Version: "0.0.1", Request: "gpu:1"}
{
    "addError": {
        "cdiVersion": "0.0.1",
        "code": 103,
        "msg": "Container ID not specified"
    }
}

#DEL: DelRequest{Version: "0.0.1", RequestId: "1234"}
{
    "cdiVersion": "0.0.1"
}

#DEL: DelRequest{Version: "0.0.1", RequestId: "3456"}
{
    "delError": {
        "cdiVersion": "0.0.1",
        "code": 5,
        "msg": "Unknown container ID"
    }
}
```
