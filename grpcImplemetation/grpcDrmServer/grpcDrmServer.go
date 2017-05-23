package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/tiwariHD/goDrmSys"

	pb "github.com/tiwariHD/commandProto"
	"golang.org/x/net/context"
)

//constants for commands and path
const (
	VERSION       = "VERSION"
	INFO          = "INFO"
	ADD           = "ADD"
	DEL           = "DEL"
	EMPTY         = ""
	DRMPATH       = "../drmFiles"
	VERSIONPLUGIN = "0.0.2"
	port          = ":50051"
)

//constants for error codes in spec (0-99)
const (
	ERR_NONE = iota
	ERR_VERSION_ENV
	ERR_VERSION_CONF
	ERR_COMMAND_UNSUPPORT
	ERR_RESOURCE_UNSUPPORT
	ERR_REQUESTID_UNKNOWN
)

//constants for other errors (>99)
const (
	ERR_CMD_EMPTY = iota + 100
	ERR_OTHER
	ERR_RESOURCE_UNAVAIL
	ERR_REQUESTID_EMPTY
	ERR_CONF_READ
)

//global variables for plugin
var (
	//errorMsg contains description for error codes
	errorMsg = map[int]string{
		ERR_NONE:               "No Error",
		ERR_VERSION_ENV:        "CDI Version of env-var not supported",
		ERR_VERSION_CONF:       "CDI Version of config not supported",
		ERR_COMMAND_UNSUPPORT:  "Command not supported",
		ERR_RESOURCE_UNSUPPORT: "Resource spec is not supported",
		ERR_REQUESTID_UNKNOWN:  "Unknown container ID",
		ERR_CMD_EMPTY:          "Command not specified",
		ERR_OTHER:              "Other error",
		ERR_RESOURCE_UNAVAIL:   "Resource unavailable",
		ERR_REQUESTID_EMPTY:    "Container ID not specified",
		ERR_CONF_READ:          "Configuration file read error",
	}

	//supportedVersions of the plugin
	supportedVersions = []string{"0.0.1", "0.0.2"}

	//envCdiVersion stores cdiVersion from environment variable
	envCdiVersion string

	//conf stores conf file values in struct
	conf DrmConf
)

//DrmConf structure for storing conf file details
type DrmConf struct {
	CdiVersion string `json:"cdiVersion"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Plugin     string `json:"plugin"`
	Args       ArgStr `json:"args"`
}

//VersionReply for VERSION command output
type VersionReply struct {
	CdiVersion        string   `json:"cdiVersion"`
	SupportedVersions []string `json:"supportedVersions"`
}

//InfoReply for REPLY command output
type InfoReply struct {
	CdiVersion string   `json:"cdiVersion"`
	Gpu        int      `json:"gpu"`
	Devices    []string `json:"devices"`
}

//AddReply for ADD command output
type AddReply struct {
	CdiVersion string   `json:"cdiVersion"`
	Devices    []string `json:"devices"`
}

//DelReply for DEL command output
type DelReply struct {
	CdiVersion string `json:"cdiVersion"`
}

//ErrorReply for error message output
type ErrorReply struct {
	CdiVersion string `json:"cdiVersion"`
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
	Details    string `json:"details,omitempty"`
}

//ArgStr contains fields for arguments in conf file
type ArgStr struct {
	WantDeviceNodes   []string `json:"want_device_nodes"`
	VendorIDWhitelist []string `json:"vendorid_whitelist"`
}

//GpuInfo contain device names and path of directory
type GpuInfo struct {
	Num      int
	DeviceID []string
	VendorID []string
	DirPath  []string
	DevNames []goDrmSys.DeviceNodes
}

//getJsonStruct prints json ouput on stdin
func getJsonStruct(i interface{}) []byte {
	out, _ := json.MarshalIndent(i, "", "    ")
	//fmt.Printf("%s\n", out)
	return out
}

//getErrorStruct is wrapper function for fetching error details
func (e *ErrorReply) getErrorStruct(err int, detail string) {
	//*e = ErrorReply{envCdiVersion, err, errorMsg[err], msg}
	e.CdiVersion = envCdiVersion
	e.Code = err
	e.Msg = errorMsg[err]
	e.Details = detail
}

//isDirExists checks whether directory already exists
func isDirExists(dpath string) (bool, ErrorReply) {
	var e ErrorReply
	if _, err := os.Stat(DRMPATH); err == nil {
		return true, e
	} else if os.IsNotExist(err) {
		//False case, pass
	} else {
		e.getErrorStruct(ERR_OTHER, fmt.Sprintf("Error: %s", err))
	}

	return false, e
}

//isDirEmpty checks whether directory is empty
func isDirEmpty(dpath string) (bool, ErrorReply) {
	var r bool
	var e ErrorReply
	f, err := os.Open(dpath)
	if err != nil {
		e.getErrorStruct(ERR_OTHER, fmt.Sprintf("Error: %s", err))
	} else {
		defer f.Close()

		// reads atleast 1 name from directory
		_, err = f.Readdirnames(1)
		if err == io.EOF {
			r = true
		}
	}
	return r, e
}

//makeDir creates a new directory specified by path, path can be relative also
func makeDir(dpath string) ErrorReply {
	var e ErrorReply
	if err := os.Mkdir(dpath, os.ModePerm); err != nil {
		e.getErrorStruct(ERR_OTHER, fmt.Sprintf("Error: %s", err))
	}
	return e
}

//getPciBusID fetches pcibus info from goDrmSys package
func getPciBusID(d *goDrmSys.DeviceInfo) string {
	// returns bus ids
	return d.BusInfo.GetBusInfo()
}

//getPciVendorID fetches pcidevice info from goDrmSys package
func getPciVendorID(d *goDrmSys.DeviceInfo) string {
	// returns vendor id
	return (strings.Split(d.DevInfo.GetDevInfo(), ":")[0])
}

//getDevices fetches all the info for gpus from goDrmSys package
func getDevices() (GpuInfo, ErrorReply) {
	var e ErrorReply
	var nodes GpuInfo
	//check drm available on host
	if !goDrmSys.DrmAvailable() {
		e.getErrorStruct(ERR_OTHER, fmt.Sprintf("DRM unavailable!!"))
		return nodes, e
	}

	dev := goDrmSys.GetDevices()
	nodes.Num = len(dev)

	// parses device data
	for i := 0; i < len(dev); i++ {
		devID := getPciBusID(&dev[i].Info)
		nodes.DeviceID = append(nodes.DeviceID, devID)
		nodes.VendorID = append(nodes.VendorID, getPciVendorID(&dev[i].Info))
		nodes.DirPath = append(nodes.DirPath, filepath.Join(DRMPATH, devID))
		nodes.DevNames = append(nodes.DevNames, dev[i].Nodes)
	}

	//check directory exists otherwise create
	if exists, err := isDirExists(DRMPATH); err.Code == ERR_NONE {
		if exists == false {
			makeDir(DRMPATH)
		}
	} else {
		e = err
		return nodes, e
	}
	// creates folder for gpus identified by pci ids
	if empty, err := isDirEmpty(DRMPATH); err.Code == ERR_NONE {
		if empty == true {
			for i := 0; i < len(dev); i++ {
				// make dirs for devs
				makeDir(nodes.DirPath[i])
			}
		}
	} else {
		e = err
		return nodes, e
	}

	return nodes, e
}

//gpuInWhitelist checks for gpu vendor id from configuration file
func gpuInWhitelist(vID string) bool {
	found := false
	//iterate over list from conf
	for _, vConfID := range conf.Args.VendorIDWhitelist {
		if vID == vConfID {
			found = true
			break
		}
	}
	return found
}

//checkVersion compares version with supportedVersions array
func checkVersion(ver string) bool {
	found := false
	//iterate over supportedVersions[]
	for _, confVer := range supportedVersions {
		if ver == confVer {
			found = true
			break
		}
	}
	return found
}

func checkCdiVersions() ErrorReply {
	var e ErrorReply
	//compare version of configuration file with supportedVersions
	if checkVersion(conf.CdiVersion) == false {
		e.getErrorStruct(ERR_VERSION_CONF,
			fmt.Sprintf("Unsupported version: %s", conf.CdiVersion))
		e.CdiVersion = VERSIONPLUGIN
		return e
	}

	//compare version of environment variable with supportedVersions
	//envCdiVersion = GetVersion()
	if checkVersion(envCdiVersion) == false {
		e.getErrorStruct(ERR_VERSION_ENV,
			fmt.Sprintf("Unsupported version: %s", envCdiVersion))
		e.CdiVersion = VERSIONPLUGIN
	}
	return e
}

//VERSION command reply
func version() VersionReply {
	var r VersionReply
	r.CdiVersion = VERSIONPLUGIN
	r.SupportedVersions = supportedVersions
	return r
}

//INFO command reply
func info(nodes *GpuInfo) InfoReply {
	var r InfoReply
	r.CdiVersion = VERSIONPLUGIN
	r.Gpu = nodes.Num

	for i := 0; i < nodes.Num; i++ {
		r.Devices = append(r.Devices, nodes.DeviceID[i])
	}
	return r
}

//ADD command reply
func add(nodes *GpuInfo, num int, conID string) (AddReply, ErrorReply) {
	var r AddReply
	var e ErrorReply
	r.CdiVersion = envCdiVersion

	//check if whitelist is populated
	whitelist := false
	if len(conf.Args.VendorIDWhitelist) > 0 {
		whitelist = true
	}

	// count no of free gpus
	count := 0
	availableGpus := make([]int, nodes.Num)
	for i, dpath := range nodes.DirPath {
		if empty, err := isDirEmpty(dpath); err.Code == ERR_NONE {
			if empty == true {
				if whitelist == true && gpuInWhitelist(nodes.VendorID[i]) == false {
					continue
				} else {
					availableGpus[count] = i
					count++
				}
			}
		} else {
			e = err
			return r, e
		}
	}

	// less no of gpu available
	if num > count {
		e.getErrorStruct(ERR_RESOURCE_UNAVAIL, fmt.Sprintf("No of GPU available: %d", count))
		return r, e
	}

	// create container id folder inside gpu folder
	for i := 0; i < num; i++ {
		// append required types of nodes
		for _, devType := range conf.Args.WantDeviceNodes {
			if devName := nodes.DevNames[availableGpus[i]].NodeMap[devType]; devName != "" {
				r.Devices = append(r.Devices, devName)
			}
		}

		//create directory for container
		dpath := filepath.Join(nodes.DirPath[availableGpus[i]], string(conID))
		if err := os.Mkdir(dpath, os.ModePerm); err != nil {
			e.getErrorStruct(ERR_OTHER, fmt.Sprintf("Error: %s", err))
			return r, e
		}
	}

	return r, e
}

//DEL command reply
func del(nodes *GpuInfo, conID string) (DelReply, ErrorReply) {
	var r DelReply
	var e ErrorReply

	// search for container id inside gpu folder then delete
	found := false
	for _, dpath := range nodes.DirPath {
		fpath := filepath.Join(dpath, conID)

		if _, err1 := os.Stat(fpath); err1 == nil {
			if err2 := os.Remove(fpath); err2 != nil {
				e.getErrorStruct(ERR_OTHER, fmt.Sprintf("Error: %s", err2))
				return r, e
			}

			found = true
		}
	}

	if found == false {
		e.getErrorStruct(ERR_REQUESTID_UNKNOWN, "")
	} else {
		r = DelReply{envCdiVersion}
	}

	return r, e
}

// server is used to implement commandProto.CmdProtoServer
type server struct{}

// GetReply implements commandProto.CmdProtoServer
func (s *server) GetReply(ctx context.Context, in *pb.CmdRequest) (*pb.CmdReply,
	error) {
	var retrn pb.CmdReply
	var e ErrorReply

	//get command from client and run required function
	command := in.GetCommand()
	switch command {
	case VERSION:
		//VERSION command should never fail
		reply := version()
		retrn.Message = getJsonStruct(reply)
		//retrn.Message = VERSIONPLUGIN
	case INFO:
		//store device data
		nodes, err := getDevices()
		if err.Code != ERR_NONE {
			retrn.Message = getJsonStruct(err)
		}
		r := info(&nodes)
		retrn.Message = getJsonStruct(r)

	case ADD:
		envCdiVersion = in.GetVersion()
		if err := checkCdiVersions(); err.Code != ERR_NONE {
			retrn.Message = getJsonStruct(err)
			break
		}
		//check if gpu types present
		if len(conf.Args.WantDeviceNodes) == 0 {
			e.getErrorStruct(ERR_OTHER, fmt.Sprintf("want_device_nodes list is empty"))
			retrn.Message = getJsonStruct(e)
			break
		}

		req := in.GetRequest()
		num, err := strconv.Atoi(strings.TrimPrefix(req, "gpu:"))
		if err != nil {
			e.getErrorStruct(ERR_RESOURCE_UNSUPPORT,
				fmt.Sprintf("Unsupported resource request: %s", req))
			retrn.Message = getJsonStruct(e)
			break
		}

		id := in.GetRequestId()
		if id == EMPTY {
			e.getErrorStruct(ERR_REQUESTID_EMPTY, "")
			retrn.Message = getJsonStruct(e)
			break
		}

		//store device data
		if nodes, err1 := getDevices(); err1.Code == ERR_NONE {
			if r, err2 := add(&nodes, num, id); err2.Code == ERR_NONE {
				retrn.Message = getJsonStruct(r)
			} else {
				retrn.Message = getJsonStruct(err2)
			}
		} else {
			retrn.Message = getJsonStruct(err1)
		}

	case DEL:
		envCdiVersion = in.GetVersion()
		if err := checkCdiVersions(); err.Code != ERR_NONE {
			retrn.Message = getJsonStruct(err)
			break
		}

		id := in.GetRequestId()
		if id == EMPTY {
			e.getErrorStruct(ERR_REQUESTID_EMPTY, "")
			retrn.Message = getJsonStruct(e)
			break
		}

		//store device data
		if nodes, err1 := getDevices(); err1.Code == ERR_NONE {
			if r, err2 := del(&nodes, id); err2.Code == ERR_NONE {
				retrn.Message = getJsonStruct(r)
			} else {
				retrn.Message = getJsonStruct(err2)
			}
		} else {
			retrn.Message = getJsonStruct(err1)
		}

	case EMPTY:
		e.getErrorStruct(ERR_CMD_EMPTY, "")
		retrn.Message = getJsonStruct(e)
		//retrn.Message = "=command empty"

	default:
		e.getErrorStruct(ERR_COMMAND_UNSUPPORT,
			fmt.Sprintf("Unsupported command %s", command))
		retrn.Message = getJsonStruct(e)
		//retrn.Message = "=command not supported"
	}

	return &retrn, nil
}

func main() {
	//get conf file from stdin
	dec := json.NewDecoder(os.Stdin)
	if err := dec.Decode(&conf); err != nil {
		log.Fatalf("%s, %d\n", errorMsg[ERR_CONF_READ], ERR_CONF_READ)
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
	s := grpc.NewServer()
	pb.RegisterCmdProtoServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}
