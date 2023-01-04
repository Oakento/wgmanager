package app

import (
	"fmt"
	"net"
	"os"
	"wgmanager/infra"
	"wgmanager/util"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

func CmdSetHandler(cmd *cobra.Command, params []string) {

	if len(params) < 1 {
		fmt.Println("Command requires username")
		return
	}
	username := params[0]
	user := infra.SelectDbPeer(username)
	if user.Username == "" {
		fmt.Println("User do not exist")
		return
	}
	isSetKey, err := cmd.Flags().GetBool("key")
	if err != nil {
		fmt.Println(err)
		return
	}
	addr, err := cmd.Flags().GetString("address")
	if err != nil {
		fmt.Println(err)
		return
	}
	if isSetKey {
		pubKey, priKey := util.NewWgKeyPairs()
		user.PublicKey = pubKey
		user.PrivateKey = priKey
		// set new key
	}
	var ipAddress string
	if addr != "" {
		_, ipCIDR, err := net.ParseCIDR(addr)
		if err != nil {
			ip := net.ParseIP(addr)
			if ip == nil {
				fmt.Println("Subnet IP address not valid")
				return
			}
			ipAddress = util.StringConcat(ip.String(), "/32")
		} else {
			ipAddress = ipCIDR.String()
		}
		peer := infra.SelectDbPeerByIpAddress(ipAddress)
		if peer.Username != "" {
			fmt.Println("Subnet IP address already exists")
			return
		}
		user.IpAddress = ipAddress
	}
	infra.UpdateDbPeer(user)
	SyncWgConfigWithDb()
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Username", "Public Key", "Private Key", "IP Address"})
	t.AppendRow(table.Row{user.Username, user.PublicKey, user.PrivateKey, user.IpAddress})
	t.AppendSeparator()

	t.Render()
}
