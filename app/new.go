package app

import (
	"fmt"
	"os"
	"wgmanager/infra"
	"wgmanager/util"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

func CreateNewPeer(username string, ipAddress string) *infra.PeerEntity {

	publicKey, privateKey := util.NewWgKeyPairs()
	p := &infra.PeerEntity{
		Username:   username,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
		IpAddress:  ipAddress,
	}
	infra.InsertDbPeer(p)

	return p
}

func SyncWgConfigWithDb() {
	users := infra.SelectAllDbPeer()

	config := infra.ConfigEntity{readConfig(WgConfigFile)[0]}
	for _, u := range users {
		if u.Username == "host" {
			continue
		}
		config = append(config, infra.NewWgConfigPeerSection().Init(
			u.IpAddress,
			"",
			u.PublicKey,
			// infra.DEFAULT_WG_PERSISTENT_KEEPALIVE,
			"",
		))
	}
	CreateConfig(WgConfigFile, config)
}

func CmdNewHandler(cmd *cobra.Command, usernames []string) {

	addrs, err := cmd.Flags().GetStringSlice("address")
	if err != nil {
		fmt.Println(err)
	}
	addrsN := len(addrs)
	var ip string
	for i := len(usernames) - addrsN; i > 0; i-- {
		ip, IpPool = IpPool[0], IpPool[1:]
		addrs = append(addrs, ip)
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Username", "Public Key", "Private Key", "IP Address"})
	for i := range usernames {
		p := CreateNewPeer(usernames[i], addrs[i])
		t.AppendRows([]table.Row{
			{p.Username, p.PublicKey, p.PrivateKey, p.IpAddress},
		})
		t.AppendSeparator()
	}
	SyncWgConfigWithDb()

	t.Render()
}
