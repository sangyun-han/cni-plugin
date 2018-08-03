package pronet

import (
	"github.com/vishvananda/netlink"
	"fmt"
	"net"
	"os"
	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	"github.com/sangyun-han/cni-plugin/utils"
)

func bridgeByName(name string) (*netlink.Bridge, error) {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return nil, fmt.Errorf("could not find %q bridge : %v", name, err)
	}
	br, ok := link.(*netlink.Bridge)
	if !ok {
		return nil, fmt.Errorf("%q already exists but is not a bridge", name)
	}
	return br, nil
}

func createVethPair(name, peer string, mtu int) (netlink.Link, error) {
	veth := &netlink.Veth {
		LinkAttrs: netlink.LinkAttrs{
			Name: name,
			Flags: net.FlagUp,
			MTU: mtu,
		},
		PeerName: peer,
	}
	if err := netlink.LinkAdd(veth); err != nil {
		return nil, err
	}
	return veth, nil
}

func getVethPair(first, second string) (firstLink, secondLink netlink.Link, err error) {
	firstLink, err = createVethPair(first, second, 1500)
	if err != nil {
		switch {
		case os.IsExist(err):
			err = fmt.Errorf("already exists ", first)
			return
		default:
			err = fmt.Errorf("failed to create veth pair: %v", err)
			return
		}
	}
	if secondLink, err = netlink.LinkByName(second); err != nil {
		err = fmt.Errorf("Failed to find %q: %v ", second, err)
	}
	return
}

func addVlanInterface(parentIF string, vlanId int, devName string) (err error) {
	var parentIndex netlink.Link
	if parentIndex, err = netlink.LinkByName(parentIF); err != nil {
		return fmt.Errorf("Failed to get %s: %v", parentIF, err)
	}

	vlanConf := netlink.Vlan{
		LinkAttrs: netlink.LinkAttrs{
			Name: devName,
			ParentIndex: parentIndex.Attrs().Index,
		},
		VlanId:vlanId,
	}

	if err = netlink.LinkAdd(&vlanConf); err != nil {
		return fmt.Errorf("Failed to add vlan %s : %v", devName, err)
	}
	return nil
}



func cmdAdd(args *skel.CmdArgs) error {
	conf, err := utils.ParseConfig(args.StdinData)

	if err != nil {
		return err
	}

	return types.PrintResult(conf.PrevResult, conf.CNIVersion)
}

func cmdDel(args *skel.CmdArgs) error {
	conf, err := utils.ParseConfig(args.StdinData)
	if err != nil {
		return err
	}
	//TODO : implement
	_ = conf
	return nil
}