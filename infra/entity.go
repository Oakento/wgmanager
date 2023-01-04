package infra

import "wgmanager/util"

type ConfigEntity []ConfigSectionEntity
type ConfigSectionEntity interface{}

type AppConfigDefaultSectionEntity struct {
	name     string `section:"default"`
	Db       string `option:"database"`
	WgConfig string `option:"wireguard_config"`
	PublicIP string `option:"public_ip"`
	Subnet   string `option:"subnet"`
}

func NewAppConfigDefaultSection() *AppConfigDefaultSectionEntity {
	return &AppConfigDefaultSectionEntity{
		name: "default",
	}
}
func (s *AppConfigDefaultSectionEntity) Init(db string, wgConfig string, publicIP string, subnet string) *AppConfigDefaultSectionEntity {
	s.Db = db
	s.WgConfig = wgConfig
	s.PublicIP = publicIP
	s.Subnet = subnet
	return s
}

type WgConfigInterfaceSectionEntity struct {
	name       string `section:"Interface"`
	Address    string `option:"Address"`
	ListenPort string `option:"ListenPort"`
	PrivateKey string `option:"PrivateKey"`
	DNS        string `option:"DNS"`
	MTU        string `option:"MTU"`
	PreUp      string `option:"PreUp"`
	PostUp     string `option:"PostUp"`
	PreDown    string `option:"PreDown"`
	PostDown   string `option:"PostDown"`
}

func NewWgConfigInterfaceSection() *WgConfigInterfaceSectionEntity {
	return &WgConfigInterfaceSectionEntity{
		name: "Interface",
	}
}
func (s *WgConfigInterfaceSectionEntity) Init(
	address string, listenPort string, privateKey string, dns string, mtu string,
	preUp string, postUp string, preDown string, postDown string) *WgConfigInterfaceSectionEntity {
	s.Address = address
	s.ListenPort = listenPort
	s.PrivateKey = privateKey
	s.DNS = dns
	s.MTU = mtu
	s.PreUp = preUp
	s.PostUp = postUp
	s.PreDown = preDown
	s.PostDown = postDown
	return s
}

type WgConfigPeerSectionEntity struct {
	name                string `section:"Peer"`
	AllowedIPs          string `option:"AllowedIPs"`
	Endpoint            string `option:"Endpoint"`
	PublicKey           string `option:"PublicKey"`
	PersistentKeepalive string `option:"PersistentKeepalive"`
}

func NewWgConfigPeerSection() *WgConfigPeerSectionEntity {
	return &WgConfigPeerSectionEntity{
		name: "Peer",
	}
}
func (s *WgConfigPeerSectionEntity) Init(allowIPs string, endpoint string, publicKey string, persistentKeepalive string) *WgConfigPeerSectionEntity {
	s.AllowedIPs = allowIPs
	s.Endpoint = endpoint
	s.PublicKey = publicKey
	s.PersistentKeepalive = persistentKeepalive
	return s
}

func NewConfigSection(sectionName string) ConfigSectionEntity {
	switch sectionName {
	case "default":
		return NewAppConfigDefaultSection()
	case "Interface":
		return NewWgConfigInterfaceSection()
	case "Peer":
		return NewWgConfigPeerSection()
	}
	return nil
}

type AppConfigEntity struct {
	Db       string
	WgConfig string
	PublicIP string
	Subnet   string
}

type WgConfigEntity struct {
	Address    string
	ListenPort string
	PrivateKey string
	DNS        string
	MTU        string
	PreUp      string
	PostUp     string
	PreDown    string
	PostDown   string
	WgPeers    []*WgConfigPeerSectionEntity
}

type PeerEntity struct {
	Username   string
	PublicKey  string
	PrivateKey string
	IpAddress  string
}

func NewPeer(username string, ipAddress string) *PeerEntity {
	p := &PeerEntity{
		Username:  username,
		IpAddress: ipAddress,
	}
	pub, pri := util.NewWgKeyPairs()
	p.PublicKey = pub
	p.PrivateKey = pri
	return p
}
