package app

import (
	"fmt"
	"os"
	"path/filepath"
	"wgmanager/infra"
	"wgmanager/util"
)

var AppConfigFile string
var DbFile string
var WgConfigFile string
var Subnet string
var IpPool []string
var AppConfig *infra.AppConfigEntity
var WgConfig *infra.WgConfigEntity

func CreateConfig(file string, config infra.ConfigEntity) {
	infra.MkdirIfNotExist(filepath.Dir(file), 0755)
	cf, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		infra.HandlePermissionOrFileSystemError(err)
	}
	defer cf.Close()

	err = infra.EncodeConfig(cf, config)
	if err != nil {
		infra.HandlePermissionOrFileSystemError(err)
	}
}

func createAppConfig() {
	config := infra.ConfigEntity{
		infra.NewAppConfigDefaultSection().Init(
			util.ConvertToAbsolutePath(DbFile),
			util.ConvertToAbsolutePath(WgConfigFile),
			util.GetPublicIP(),
			Subnet,
		),
	}
	CreateConfig(AppConfigFile, config)
}

func createWgConfig(host *infra.PeerEntity) {
	config := infra.ConfigEntity{
		infra.NewWgConfigInterfaceSection().Init(
			infra.DEFAULT_WG_IP_ADDRESS,
			infra.DEFAULT_WG_LISTEN_PORT,
			host.PrivateKey,
			infra.DEFAULT_WG_DNS,
			infra.DEFAULT_WG_MTU,
			infra.DEFAULT_WG_PREUP,
			infra.DEFAULT_WG_POSTUP,
			infra.DEFAULT_WG_PREDOWN,
			infra.DEFAULT_WG_POSTDOWN,
		),
	}
	CreateConfig(WgConfigFile, config)
}

func readConfig(file string) infra.ConfigEntity {
	cf, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	defer cf.Close()
	appcfg := infra.ConfigEntity{}
	return infra.DecodeConfig(cf, appcfg)
}

func readAppConfig() *infra.AppConfigEntity {
	res := readConfig(AppConfigFile)
	for _, s := range res {
		c, ok := s.(*infra.AppConfigDefaultSectionEntity)
		if !ok {
			continue
		}
		return &infra.AppConfigEntity{
			Db:       c.Db,
			WgConfig: c.WgConfig,
			PublicIP: c.PublicIP,
			Subnet:   c.Subnet,
		}
	}
	return &infra.AppConfigEntity{}
}

func readWgConfig() *infra.WgConfigEntity {
	wgConfig := &infra.WgConfigEntity{
		WgPeers: make([]*infra.WgConfigPeerSectionEntity, 0),
	}
	res := readConfig(WgConfigFile)
	for _, s := range res {
		intf, ok := s.(*infra.WgConfigInterfaceSectionEntity)
		if !ok {
			peer, ok := s.(*infra.WgConfigPeerSectionEntity)
			if !ok {
				continue
			}
			wgConfig.WgPeers = append(wgConfig.WgPeers, peer)
			continue
		}
		wgConfig.Address = intf.Address
		wgConfig.ListenPort = intf.ListenPort
		wgConfig.PrivateKey = intf.PrivateKey
		wgConfig.DNS = intf.DNS
		wgConfig.MTU = intf.MTU
		wgConfig.PreUp = intf.PreUp
		wgConfig.PostUp = intf.PostUp
		wgConfig.PreDown = intf.PreDown
		wgConfig.PostDown = intf.PostDown
	}
	return wgConfig
}
func Init(dbChanged bool, wgConfigChanged bool, subnetChanged bool) {
	_, appConfigFileErr := os.Stat(AppConfigFile) // 2 conditions: default or specified by cmd
	if appConfigFileErr != nil {                  // 3 conditions: not found, no permission, file system error
		if os.IsNotExist(appConfigFileErr) {
			createAppConfig()
		} else {
			infra.HandlePermissionOrFileSystemError(appConfigFileErr)
		}
	}

	AppConfig = readAppConfig()
	// if db specified by cmd: use DbFile; else use config value
	// if wg specified by cmd: use WgConfigFile else use config value
	if !dbChanged {
		DbFile = AppConfig.Db
	}
	if !wgConfigChanged {
		WgConfigFile = AppConfig.WgConfig
	}
	if !subnetChanged {
		Subnet = AppConfig.Subnet
	}

	_, dbFileErr := os.Stat(DbFile)
	if dbFileErr != nil {
		if os.IsNotExist(dbFileErr) {
			infra.InitDb(DbFile)
			infra.CreateDb()
			pubkey, prikey := util.NewWgKeyPairs()
			subnetIPs := util.GetIPListFromSubnet(Subnet)
			if len(subnetIPs) < 2 {
				fmt.Println("Subnet not valid")
				return
			}
			infra.InsertDbPeer(&infra.PeerEntity{
				Username:   "host",
				PublicKey:  pubkey,
				PrivateKey: prikey,
				IpAddress:  subnetIPs[0],
			})
		} else {
			infra.HandlePermissionOrFileSystemError(dbFileErr)
		}
	}
	infra.InitDb(DbFile)

	_, wgConfigFileErr := os.Stat(WgConfigFile)
	if wgConfigFileErr != nil {
		if os.IsNotExist(wgConfigFileErr) {
			host := infra.SelectDbPeer("host")
			createWgConfig(host)
		} else {
			infra.HandlePermissionOrFileSystemError(wgConfigFileErr)
		}
	}
	WgConfig = readWgConfig()

	IpPool = util.Difference(util.GetIPListFromSubnet(Subnet), infra.SelectAllDbPeerIpAddress())
}
