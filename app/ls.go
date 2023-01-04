package app

import (
	"os"
	"wgmanager/infra"

	"github.com/jedib0t/go-pretty/v6/table"
)

func CmdLsHandler() {

	peers := infra.SelectAllDbPeer()
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Username", "Public Key", "Private Key", "IP Address"})
	for _, p := range peers {
		t.AppendRows([]table.Row{
			{p.Username, p.PublicKey, p.PrivateKey, p.IpAddress},
		})
		t.AppendSeparator()
	}

	t.Render()
}
