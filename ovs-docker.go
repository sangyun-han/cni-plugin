package cni_plugin

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
