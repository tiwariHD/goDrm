package main

import (
    "os"
    "fmt"
    "encoding/json"
    "github.com/tiwariHD/goDrmSys"
)

const (
    Info        = "INFO"
    Add         = "ADD"
    Del         = "DEL"
)

type InfoReply struct {
    DeviceIds       []string
}

type AddReply struct {
    PrimaryNode     string
}

type ErrorReply struct {
    Code            int
    Message         string
}

func jsonOut(i interface{}) {
    out, _ := json.MarshalIndent(i, "", "    ")
    fmt.Println(string(out))
}

func getDeviceId(d goDrmSys.DeviceInfo) string {

    //add
    return d.BusInfo.GetBusInfo()
}

func info() (r InfoReply, e ErrorReply){

    dev := goDrmSys.GetDevices()
    r.DeviceIds = make([]string, len(dev))
    jsonOut(dev)
    for i := 0; i <  len(dev); i++ {
        r.DeviceIds[i] = getDeviceId(dev[i].Info)
    }
    //jsonOut(r)
    return
}

func add(devId string) (r AddReply, e ErrorReply){
    //add
    dev := goDrmSys.GetDevices()
    var d *goDrmSys.Device
    for i := 0; i < len(dev); i++ {
        if devId == getDeviceId(dev[i].Info) {
            d = &dev[i]
        }
    }

    if d != nil {
        r.PrimaryNode = d.Nodes.Primary
    } else {
        e.Code = 1
        e.Message = fmt.Sprintf("%s not found", devId)
    }

    return
}

func main() {

    command := os.Getenv("CDI_COMMAND")

    if !goDrmSys.DrmAvailable() {
        fmt.Println("DRM unavailable!!")
        os.Exit(1)
    }

    switch command {

        case Info:
            r, e := info()
            if e.Code == 0 {
                // No error
                jsonOut(r)
            } else {
                jsonOut(e)
                os.Exit(1)
            }

        case Add:
            id := os.Getenv("CDI_DEVICE_ID")
            r, e := add(id)
            if e.Code == 0 {
                // No error
                jsonOut(r)
            } else {
                jsonOut(e)
                os.Exit(1)
            }

        case Del:
            //pass

        default:
            jsonOut(ErrorReply{1, command})
            os.Exit(1)
    }
}

