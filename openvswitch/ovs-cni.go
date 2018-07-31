package openvswitch

import (
	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	"encoding/json"
	"fmt"
	"github.com/containernetworking/cni/pkg/types/current"
	"github.com/containernetworking/cni/pkg/version"
	"os/exec"
	"bytes"
	"os"
	"net"
	"github.com/john-lin/ovsdb"
	"time"
	"github.com/vishvananda/netlink"
)

const OVS_CMD_PATH = "/usr/bin"
const OVS_DOCKER_CMD = "ovs-docker"
const (
	ADD_PORT = "add-port"
	DEL_PORT = "del-port"
)

type CNIConf struct {
	//libcni.RuntimeConf
	types.NetConf
	RuntimeConfig *struct {
		SampleConfig map[string]interface{} `json:"sample"`
	} `json:"runtimeConfig"`

	PrevResult *current.Result `json:"-"`
}

type OpenVSwitch struct {
	BridgeName string
	MACAddr    string
	CtrlAddr   net.IP
	CtrlPort   int
	OVSDB      *ovsdb.OvsDriver
}

func NewOpenVSwitch(bridgeName string) (*OpenVSwitch, error) {
	sw := new(OpenVSwitch)
	sw.BridgeName = bridgeName
	sw.OVSDB = ovsdb.NewOvsDriverWithUnix(bridgeName)

	if !sw.OVSDB.IsBridgePresent(bridgeName) {
		err := sw.OVSDB.CreateBridge(bridgeName, "standalone", true)
		if err != nil {
			return nil, err
		}
	}

	time.Sleep(300 * time.Millisecond)

	err := setLinkUp(bridgeName)
	if err != nil {
		return nil, err
	}

	return sw, nil
}

// VLAN will be added
func (sw *OpenVSwitch) addPort(ifName string) error {
	if !sw.OVSDB.IsPortNamePresent(ifName) {
		err := sw.OVSDB.CreatePort(ifName, "", 0)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sw *OpenVSwitch) delPort(ifName string) error {
	if sw.OVSDB.IsPortNamePresent(ifName) {
		err := sw.OVSDB.DeletePort(ifName)
		if err != nil {
			return err
		}
	}
	return nil
}

func setLinkUp(ifName string) error {
	iface, err := netlink.LinkByName(ifName)
	if err != nil {
		return err
	}
	return netlink.LinkSetUp(iface)
}

func parseConfig(stdin []byte) (*CNIConf, error) {
	conf := CNIConf{}

	if err := json.Unmarshal(stdin, &conf); err != nil {
		return nil, fmt.Errorf("failed to parse network configuration: %v", err)
	}

	return &conf, nil
}

// ovs-vsctl add-br br0
// ifconfig br0 10.0.1.1 netmask 255.255.255.0 up
// ovs-docker add-port BRIDGE_NAME ETH CONTAINER_NAME --ipaddress=<ip/subnet>
func cmdAdd(args *skel.CmdArgs) error {
	conf, err := parseConfig(args.StdinData)

	if err != nil {
		return err
	}

	brName := "br0"

	// make command script to add container
	cmd := exec.Command(OVS_DOCKER_CMD, ADD_PORT, brName, args.IfName, args.ContainerID)
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err = cmd.Run()
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	fmt.Println(string(cmdOutput.Bytes()))

	return types.PrintResult(conf.PrevResult, conf.CNIVersion)
}

func cmdDel(args *skel.CmdArgs) error {
	conf, err := parseConfig(args.StdinData)
	if err != nil {
		return err
	}
	//TODO : implement
	_ = conf
	return nil
}

func cmdGet(args *skel.CmdArgs) error {
	return fmt.Errorf("TODO")
}

func main() {
	// init code
	skel.PluginMain(cmdAdd, cmdGet, cmdDel, version.All, "TODO")
}
