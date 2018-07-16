package cni_plugin

import (
	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	"encoding/json"
	"fmt"
	"github.com/containernetworking/cni/pkg/types/current"
	"github.com/containernetworking/cni/pkg/version"
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

func cmdAdd(args *skel.CmdArgs) error {
	conf, err := parseConfig(args.StdinData)

	if err != nil {
		return err
	}

	//TODO : implement

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