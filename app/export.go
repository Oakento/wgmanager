package app

import (
	"fmt"
	"wgmanager/infra"
	"wgmanager/util"

	"github.com/spf13/cobra"
)

func CmdExportHandler(cmd *cobra.Command, usernames []string) {
	host := infra.SelectDbPeer("host")
	for _, u := range usernames {
		if u == "host" {
			continue
		}
		peer := infra.SelectDbPeer(u)
		if peer.Username == "" {
			continue
		}
		config := infra.ConfigEntity{
			infra.NewWgConfigInterfaceSection().Init(
				peer.IpAddress,
				infra.DEFAULT_WG_LISTEN_PORT,
				peer.PrivateKey,
				infra.DEFAULT_WG_DNS,
				infra.DEFAULT_WG_MTU,
				"",
				"",
				"",
				"",
			),
			infra.NewWgConfigPeerSection().Init(
				Subnet,
				util.StringConcat(AppConfig.PublicIP, ":", WgConfig.ListenPort),
				host.PublicKey,
				infra.DEFAULT_WG_PERSISTENT_KEEPALIVE,
			),
		}
		outdir, err := cmd.Flags().GetString("output")
		if err != nil {
			fmt.Println(err)
			return
		}
		filename := util.StringConcat(outdir, "/", u, ".conf")
		CreateConfig(filename, config)
	}
}
