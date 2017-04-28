```
#Device info stored as:
[
    {
        "Nodes": {
            "Primary": "/dev/dri/card1",
            "Control": "/dev/dri/controlD65",
            "Render": "/dev/dri/renderD129"
        },
        "Info": {
            "BusInfo": {
                "Domain": 0,
                "Bus": 1,
                "Dev": 0,
                "Func": 0
            },
            "DevInfo": {
                "VendorId": 4098,
                "DeviceId": 38340,
                "SubVendorid": 6058,
                "SubDeviceId": 8463,
                "RevisionId": 0
            }
        }
    },
    {
        "Nodes": {
            "Primary": "/dev/dri/card0",
            "Control": "/dev/dri/controlD64",
            "Render": "/dev/dri/renderD128"
        },
        "Info": {
            "BusInfo": {
                "Domain": 0,
                "Bus": 0,
                "Dev": 2,
                "Func": 0
            },
            "DevInfo": {
                "VendorId": 32902,
                "DeviceId": 10818,
                "SubVendorid": 6058,
                "SubDeviceId": 8467,
                "RevisionId": 7
            }
        }
    }
]

#cat drm.conf
{
    "cdiVersion": "0.0.1",
    "name": "drm-gpus",
    "plugin": "drm",
    "args": {
        "device_node_type": "all"
    }
}

#env CDI_VERSION=0.0.1 CDI_COMMAND=ADD CDI_REQUEST=gpu:1 CDI_CONTAINERID=1234 ../../../../bin/goDrm < drm.conf
{
    "cdiVersion": "0.0.1",
    "PrimaryNode": [
        "/dev/dri/card1"
    ]
}

#env CDI_VERSION=0.0.1 CDI_COMMAND=ADD CDI_REQUEST=gpu:2 CDI_CONTAINERID=1234 ../../../../bin/goDrm < drm.conf
{
    "cdiVersion": "0.0.1",
    "code": 6,
    "msg": "Resource unavailable",
    "details": "No of GPU available: 1"
}

#env CDI_VERSION=0.0.1 CDI_COMMAND=DEL CDI_CONTAINERID=1234 ../../../../bin/goDrm < drm.conf
{
    "cdiVersion": "0.0.1"
}

#env CDI_VERSION=0.0.1 CDI_COMMAND=ADD CDI_REQUEST=gpu:1,gpu-memory=2048Mi CDI_CONTAINERID=1234
../../../../bin/goDrm < drm.conf
{
    "cdiVersion": "0.0.1",
    "code": 7,
    "msg": "Resource sub-type unsupported",
    "details": "Unsupported resource request: gpu:1,gpu-memory=2048Mi"
}

#env CDI_VERSION=0.0.1 CDI_COMMAND=DEL CDI_CONTAINERID=3456 ../../../../bin/goDrm < drm.conf
{
    "cdiVersion": "0.0.1",
    "code": 8,
    "msg": "Unknown container ID"
}

#env CDI_VERSION=0.999 CDI_COMMAND=ADD CDI_REQUEST=gpu:1 CDI_CONTAINERID=1234 ../../../../bin/goDrm < drm.conf
{
    "cdiVersion": "0.0.1",
    "code": 4,
    "msg": "CDI Version not supported",
    "details": "Unsupported version: 0.999"
}

#env CDI_VERSION=0.0.1 CDI_COMMAND=ADD CDI_REQUEST=gpu:1 ../../../../bin/goDrm < drm.conf
{
    "cdiVersion": "0.0.1",
    "code": 9,
    "msg": "Container ID not specified"
}

#env CDI_VERSION=0.0.1 CDI_COMMAND=MYCMD ../../../../bin/goDrm < drm.conf
{
    "cdiVersion": "0.0.1",
    "code": 3,
    "msg": "Command not supported",
    "details": "Unsupported command MYCMD"
}

#env CDI_VERSION=0.0.1 CDI_REQUEST=gpu:1 ../../../../bin/goDrm < drm.conf
{
    "cdiVersion": "0.0.1",
    "code": 2,
    "msg": "Command not specified"
}
```
