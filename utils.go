package cni_plugin

import (
	"github.com/vishvananda/netlink"
)

func setLinkUp(ifName string) error {
	iface, err := netlink.LinkByName(ifName)
	if err != nil {
		return err
	}
	return netlink.LinkSetUp(iface)
}
