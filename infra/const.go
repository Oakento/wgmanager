package infra

const (
	DEFAULT_DIR = "/usr/local/etc/wgmanager"
	// DEFAULT_DIR    = "/home/dev/projects/wgmanager/test"
	DEFAULT_CONFIG = DEFAULT_DIR + "/wgm.conf"
	DEFAULT_DB     = DEFAULT_DIR + "/wgm.db"
	DEFAULT_WG_DIR = "/etc/wireguard"
	// DEFAULT_WG_DIR    = "/home/dev/projects/wgmanager/test/etc/wireguard"
	DEFAULT_WG_CONFIG = DEFAULT_WG_DIR + "/wg0.conf"

	DEFAULT_SUBNET                  = "10.0.8.0/24"
	DEFAULT_WG_IP_ADDRESS           = "10.0.8.1/32"
	DEFAULT_WG_LISTEN_PORT          = "51820"
	DEFAULT_WG_MTU                  = "1420"
	DEFAULT_WG_DNS                  = "8.8.8.8"
	DEFAULT_WG_PREUP                = ""
	DEFAULT_WG_PREDOWN              = ""
	DEFAULT_WG_POSTUP               = "iptables -A FORWARD -i wg0 -j ACCEPT; iptables -A FORWARD -o wg0 -j ACCEPT; iptables -t nat -A POSTROUTING -o enp0s3 -j MASQUERADE"
	DEFAULT_WG_POSTDOWN             = "iptables -D FORWARD -i wg0 -j ACCEPT; iptables -D FORWARD -o wg0 -j ACCEPT; iptables -t nat -D POSTROUTING -o enp0s3 -j MASQUERADE"
	DEFAULT_WG_PERSISTENT_KEEPALIVE = "120"
)
