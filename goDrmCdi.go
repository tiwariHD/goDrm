package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/tiwariHD/goDrmSys"
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

//jsonOut prints json ouput on stdin
func jsonOut(i interface{}) {
	out, _ := json.MarshalIndent(i, "", "    ")
	fmt.Printf("%s\n", out)
	//os.Stdout.Write(out)
}

//errorMsgOut is wrapper function for displaying error and exiting
func errorMsgOut(err int, msg string) {
	jsonOut(ErrorReply{envCdiVersion, err, errorMsg[err], msg})
	os.Exit(1)
}

//isDirExists checks whether directory already exists
func isDirExists(dpath string) bool {
	if _, err := os.Stat(DRMPATH); err == nil {
		return true
	} else if os.IsNotExist(err) {
		//False case, pass
	} else {
		errorMsgOut(ERR_OTHER, fmt.Sprintf("Error: %s", err))
	}

	return false
}

//isDirEmpty checks whether directory is empty
func isDirEmpty(dpath string) bool {
	f, err := os.Open(dpath)
	if err != nil {
		errorMsgOut(ERR_OTHER, fmt.Sprintf("Error: %s", err))
	}
	defer f.Close()

	// reads atleast 1 name from directory
	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true
	}
	return false
}

//makeDir creates a new directory specified by path, path can be relative also
func makeDir(dpath string) {
	if err := os.Mkdir(dpath, os.ModePerm); err != nil {
		errorMsgOut(ERR_OTHER, fmt.Sprintf("Error: %s", err))
	}
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
func getDevices() GpuInfo {
	//check drm available on host
	if !goDrmSys.DrmAvailable() {
		errorMsgOut(ERR_OTHER, fmt.Sprintf("DRM unavailable!!"))
	}

	var nodes GpuInfo
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

	if isDirExists(DRMPATH) == false {
		makeDir(DRMPATH)
	}
	// creates folder for gpus identified by pci ids
	if isDirEmpty(DRMPATH) == true {
		for i := 0; i < len(dev); i++ {
			// make dirs for devs
			makeDir(nodes.DirPath[i])
		}
	}

	return nodes
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

func checkCdiVersions() {
	//compare version of configuration file with supportedVersions
	if checkVersion(conf.CdiVersion) == false {
		jsonOut(ErrorReply{VERSIONPLUGIN, ERR_VERSION_CONF, errorMsg[ERR_VERSION_CONF],
			fmt.Sprintf("Unsupported version: %s", conf.CdiVersion)})
		os.Exit(1)
	}

	//compare version of environment variable with supportedVersions
	envCdiVersion = os.Getenv("CDI_VERSION")
	if checkVersion(envCdiVersion) == false {
		jsonOut(ErrorReply{VERSIONPLUGIN, ERR_VERSION_ENV, errorMsg[ERR_VERSION_ENV],
			fmt.Sprintf("Unsupported version: %s", envCdiVersion)})
		os.Exit(1)
	}
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
func add(nodes *GpuInfo, num int, conID string) AddReply {
	var r AddReply
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
		if isDirEmpty(dpath) == true {
			if whitelist == true && gpuInWhitelist(nodes.VendorID[i]) == false {
				continue
			} else {
				availableGpus[count] = i
				count++
			}
		}
	}

	// less no of gpu available
	if num > count {
		errorMsgOut(ERR_RESOURCE_UNAVAIL, fmt.Sprintf("No of GPU available: %d", count))
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
			errorMsgOut(ERR_OTHER, fmt.Sprintf("Error: %s", err))
		}
	}

	return r
}

//DEL command reply
func del(nodes *GpuInfo, conID string) DelReply {
	var r DelReply

	// search for container id inside gpu folder then delete
	found := false
	for _, dpath := range nodes.DirPath {
		fpath := filepath.Join(dpath, conID)

		if _, err1 := os.Stat(fpath); err1 == nil {
			if err2 := os.Remove(fpath); err2 != nil {
				errorMsgOut(ERR_OTHER, fmt.Sprintf("Error: %s", err2))
			}

			found = true
		}
	}

	if found == false {
		errorMsgOut(ERR_REQUESTID_UNKNOWN, "")
	} else {
		r = DelReply{envCdiVersion}
	}

	return r
}

func main() {
	//get conf file from stdin
	dec := json.NewDecoder(os.Stdin)
	if err := dec.Decode(&conf); err != nil {
		errorMsgOut(ERR_CONF_READ, "")
	}

	//get command from env and run required function
	command := os.Getenv("CDI_COMMAND")
	switch command {
	case VERSION:
		//VERSION command should never fail
		reply := version()
		jsonOut(reply)
	case INFO:
		//store device data
		nodes := getDevices()
		r := info(&nodes)
		jsonOut(r)

	case ADD:
		checkCdiVersions()
		//check if gpu types present
		if len(conf.Args.WantDeviceNodes) == 0 {
			errorMsgOut(ERR_OTHER, fmt.Sprintf("want_device_nodes list is empty"))
		}

		req := os.Getenv("CDI_REQUEST")
		num, err := strconv.Atoi(strings.TrimPrefix(req, "gpu:"))
		if err != nil {
			errorMsgOut(ERR_RESOURCE_UNSUPPORT,
				fmt.Sprintf("Unsupported resource request: %s", req))
		}

		id := os.Getenv("CDI_REQUEST_ID")
		if id == EMPTY {
			errorMsgOut(ERR_REQUESTID_EMPTY, "")
		}

		//store device data
		nodes := getDevices()
		r := add(&nodes, num, id)
		jsonOut(r)

	case DEL:
		checkCdiVersions()
		id := os.Getenv("CDI_REQUEST_ID")
		if id == EMPTY {
			errorMsgOut(ERR_REQUESTID_EMPTY, "")
		}

		//store device data
		nodes := getDevices()
		r := del(&nodes, id)
		jsonOut(r)

	case EMPTY:
		errorMsgOut(ERR_CMD_EMPTY, "")

	default:
		errorMsgOut(ERR_COMMAND_UNSUPPORT, fmt.Sprintf("Unsupported command %s", command))
	}
}
