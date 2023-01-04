package app

import (
	"fmt"
	"wgmanager/infra"

	"github.com/spf13/cobra"
)

func CmdRmHandler(cmd *cobra.Command, usernames []string) {

	for _, u := range usernames {
		if u == "host" {
			fmt.Println("Deleting host is not allowed")
			continue
		}
		affectedUser := infra.RemoveDbPeer(u)
		fmt.Println(affectedUser)
	}
	SyncWgConfigWithDb()

}
