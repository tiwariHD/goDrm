package main

import (
    "os"
    "fmt"
    "io"
    "strings"
    "strconv"
    "path/filepath"
    "encoding/json"
    "github.com/tiwariHD/goDrmSys"
)

const (
    INFO        = "INFO"
    ADD         = "ADD"
    DEL         = "DEL"
    EMPTY       = ""
    DRMPATH     = "../drmFiles"
)

const (
    ERR_CONF        = 1
    ERR_NO_CMD      = 2
    ERR_UN_CMD      = 3
    ERR_VERSION     = 4
    ERR_OTHER       = 5
    ERR_AVAIL_RES   = 6
    ERR_TYPE_RES    = 7
    ERR_UN_CONTID   = 8
    ERR_NO_CONTID   = 9
)

var ERR_MSG = []string{"", "Configuration file read error",
                    "Command not specified",
                    "Command not supported",
                    "CDI Version not supported",
                    "Other error",
                    "Resource unavailable",
                    "Resource sub-type unsupported",
                    "Unknown container ID",
                    "Container ID not specified"}

type InfoReply struct {
    CdiVersion      string      `json:"cdiVersion"`
    Device          []string    `json:"device"`
}

type AddReply struct {
    CdiVersion      string      `json:"cdiVersion"`
    PrimaryNode     []string    `json:"device`
}

type ErrorReply struct {
    CdiVersion      string      `json:"cdiVersion"`
    Code            int         `json:"code"`
    Msg             string      `json:"msg"`
    Details         string      `json:"details,omitempty"`
}

type DelReply struct {
    CdiVersion      string      `json:"cdiVersion"`
}

type ArgStr struct {
    Device_node_type string
}

type NodePaths struct {
    DirPath         []string
    DevNames        []goDrmSys.DeviceNodes
}

type DrmConf struct {
    CdiVersion      string
    Name            string
    Plugin          string
    Args            ArgStr
}

var conf DrmConf

func jsonOut(i interface{}) {

    out, _ := json.MarshalIndent(i, "", "    ")
    fmt.Println(string(out))
}

func isDirEmpty(dpath string) bool {

    f, err := os.Open(dpath)
    if err != nil {
        jsonOut(ErrorReply{conf.CdiVersion, ERR_OTHER, ERR_MSG[ERR_OTHER],
        fmt.Sprintf("Error: %s", err)})
        os.Exit(1)
    }
    defer f.Close()

    // reads atleast 1 name from directory
    _, err = f.Readdirnames(1)
    if err == io.EOF {
        return true
    } else {
        return false
    }

}

func getDeviceId(d goDrmSys.DeviceInfo) string {

    // returns pci ids
    return d.BusInfo.GetBusInfo()
}

func info() (r InfoReply, e ErrorReply, nodes NodePaths) {

    dev := goDrmSys.GetDevices()

    r.CdiVersion = conf.CdiVersion
    r.Device = make([]string, len(dev))

    nodes.DirPath = make([]string, len(dev))
    nodes.DevNames = make([]goDrmSys.DeviceNodes, len(dev))

    // parses device data
    for i := 0; i <  len(dev); i++ {
        r.Device[i] = getDeviceId(dev[i].Info)
        nodes.DirPath[i] = filepath.Join(DRMPATH, r.Device[i])
        nodes.DevNames[i] = dev[i].Nodes
    }

    // creates folder for gpus identified by pci ids
    if isDirEmpty(DRMPATH) == true {
        for i := 0; i < len(dev); i++ {
        // make dirs for devs
            if err := os.Mkdir(nodes.DirPath[i], os.ModePerm); err != nil {
                e = ErrorReply{conf.CdiVersion, ERR_OTHER,
                ERR_MSG[ERR_OTHER], fmt.Sprintf("Error: %s", err)}
                return
            }
        }
    }

    return
}

func add(num int, conId string) (r AddReply, e ErrorReply) {

    _, err, nodes := info()
    if err.Code != 0 {
        e = err
        return
    }

    // count no of free gpus
    count := 0
    idx := make([]int, len(nodes.DirPath))
    for i, dpath := range nodes.DirPath {
        if isDirEmpty(dpath) == true {
            idx[count] = i
            count++
        }
    }

    // less no of gpu available
    if num > count {
        e = ErrorReply{conf.CdiVersion, ERR_AVAIL_RES, ERR_MSG[ERR_AVAIL_RES],
        fmt.Sprintf("No of GPU available: %d", count)}
        return
    }

    r.CdiVersion = conf.CdiVersion
    r.PrimaryNode = make([]string, num)


    // create container id folder inside gpu folder
    for i := 0; i < num; i++ {
        dpath := filepath.Join(nodes.DirPath[idx[i]], string(conId))
        if err := os.Mkdir(dpath, os.ModePerm); err != nil {
            e = ErrorReply{conf.CdiVersion, ERR_OTHER,
            ERR_MSG[ERR_OTHER], fmt.Sprintf("Error: %s", err)}
            return
        }
        // returns primary node
        r.PrimaryNode[i] = nodes.DevNames[idx[i]].Primary
    }

    return
}

func del(conId string) (r DelReply, e ErrorReply) {

    _, err, nodes := info()
    if err.Code != 0 {
        e = err
        return
    }

    // search for container id inside gpu folder then delete
    found := false
    for _, dpath := range nodes.DirPath {
        fpath := filepath.Join(dpath, conId)

        if _, err := os.Stat(fpath); err == nil {

            if er:= os.Remove(fpath); er != nil {

                e = ErrorReply{conf.CdiVersion, ERR_OTHER,
                ERR_MSG[ERR_OTHER], fmt.Sprintf("Error: %s", er)}
                return
            }

            found = true
        }
    }

    if found == false {
        e = ErrorReply{conf.CdiVersion, ERR_UN_CONTID, ERR_MSG[ERR_UN_CONTID], ""}
    } else {
        r = DelReply{conf.CdiVersion}
    }

    return
}

func main() {

    dec := json.NewDecoder(os.Stdin)
    if err := dec.Decode(&conf); err != nil {
        jsonOut(ErrorReply{conf.CdiVersion, ERR_CONF, ERR_MSG[ERR_CONF], ""})
        os.Exit(1)
    }

    cdiVer := os.Getenv("CDI_VERSION")
    if conf.CdiVersion != cdiVer {
        jsonOut(ErrorReply{conf.CdiVersion, ERR_VERSION, ERR_MSG[ERR_VERSION],
        fmt.Sprintf("Unsupported version: %s", cdiVer)})
        os.Exit(1)
    }

    if !goDrmSys.DrmAvailable() {
        jsonOut(ErrorReply{conf.CdiVersion, ERR_OTHER,
        ERR_MSG[ERR_OTHER], fmt.Sprintf("DRM unavailable!!")})
        os.Exit(1)
    }

    command := os.Getenv("CDI_COMMAND")
    switch command {

        case INFO:
            r, e, _ := info()
            if e.Code == 0 {
                // No error
                jsonOut(r)
            } else {
                jsonOut(e)
                os.Exit(1)
            }

        case ADD:
            id := os.Getenv("CDI_CONTAINERID")
            if id == EMPTY {
                jsonOut(ErrorReply{conf.CdiVersion, ERR_NO_CONTID,
                ERR_MSG[ERR_NO_CONTID], ""})
                os.Exit(1)
            }

            req := os.Getenv("CDI_REQUEST")
            num, err := strconv.Atoi(strings.TrimPrefix(req, "gpu:"))
            if err != nil {
                jsonOut(ErrorReply{conf.CdiVersion, ERR_TYPE_RES,
                ERR_MSG[ERR_TYPE_RES],
                fmt.Sprintf("Unsupported resource request: %s", req)})
                os.Exit(1)
            }

            r, e := add(num, id)
            if e.Code == 0 {
                // No error
                jsonOut(r)
            } else {
                jsonOut(e)
                os.Exit(1)
            }

        case DEL:
            id := os.Getenv("CDI_CONTAINERID")
            if id == EMPTY {
                jsonOut(ErrorReply{conf.CdiVersion, ERR_NO_CONTID,
                ERR_MSG[ERR_NO_CONTID], ""})
                os.Exit(1)
            }

            r, e := del(id)
            if e.Code == 0 {
                // No error
                jsonOut(r)
            } else {
                jsonOut(e)
                os.Exit(1)
            }

        case EMPTY:
            jsonOut(ErrorReply{conf.CdiVersion, ERR_NO_CMD,
            ERR_MSG[ERR_NO_CMD], ""})
            os.Exit(1)

        default:
            jsonOut(ErrorReply{conf.CdiVersion, ERR_UN_CMD, ERR_MSG[ERR_UN_CMD],
            fmt.Sprintf("Unsupported command %s", command)})
            os.Exit(1)
    }
}

