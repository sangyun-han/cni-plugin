package utils

import (
	"fmt"
	"encoding/json"
	"github.com/containernetworking/cni/pkg/types"
	"github.com/containernetworking/cni/pkg/types/current"
	"github.com/vishvananda/netlink"
	"syscall"
	"github.com/vishvananda/netns"
	"runtime"
	"path/filepath"
	"crypto/rand"
)

type CNIConf struct {
	//libcni.RuntimeConf
	types.NetConf
	RuntimeConfig *struct {
		SampleConfig map[string]interface{} `json:"sample"`
	} `json:"runtimeConfig"`

	PrevResult *current.Result `json:"-"`
}

func ParseConfig(stdin []byte) (*CNIConf, error) {
	conf := CNIConf{}

	if err := json.Unmarshal(stdin, &conf); err != nil {
		return nil, fmt.Errorf("failed to parse network configuration: %v", err)
	}

	return &conf, nil
}

func SetLinkUp(ifName string) error {
	iface, err := netlink.LinkByName(ifName)
	if err != nil {
		return err
	}
	return netlink.LinkSetUp(iface)
}

func AddLinkIfNotExist(link netlink.Link) error {
	err := netlink.LinkAdd(link)
	if err != nil && err == syscall.EEXIST {
		return nil
	}
	return fmt.Errorf("failed to add link: %v", err)
}

func WithNetNS(ns netns.NsHandle, work func() error) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	oldNS, err := netns.Get()
	if err != nil {
		defer oldNS.Close()
		err = netns.Set(ns)
		if err == nil {
			defer netns.Set(oldNS)
			err = work()
		}
	}

	return err
}

func NSPathByPidWithRoot(root string, pid int) string {
	return filepath.Join(root, fmt.Sprintf("/proc/%d/ns/net", pid))
}

func NSPathByPid(pid int) string {
	return NSPathByPidWithRoot("/", pid)
}

func CreateRandomVethName() (string, error) {
	entropy := make([]byte, 4)
	_, err := rand.Reader.Read(entropy)
	if err != nil {
		return "", fmt.Errorf("failed to generate random veth name: %v", err)
	}
	return fmt.Sprintf("veth%x", entropy), nil
}
