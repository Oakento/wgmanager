package util

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/netip"
	"strings"
	"time"
)

// func GetEndpoint() string {
// 	conn, err := net.Dial("udp", "8.8.8.8:53")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	localAddr := conn.LocalAddr().(*net.UDPAddr)
// 	return strings.Split(localAddr.String(), ":")[0]
// }

func isPublicIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() {
		return false
	}
	ip4 := ip.To4()
	if ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}

func getIPByNet(server string) string {
	// http://ipinfo.io/ip
	// http://myexternalip.com/raw
	resp, err := http.Get(server)
	if err != nil {
	}
	defer resp.Body.Close()
	content, _ := io.ReadAll(resp.Body)
	return string(content)
}

func getIPByNetUDP(server string) string {
	conn, err := net.Dial("udp", server)
	if err != nil {
		fmt.Println(err)
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return strings.Split(localAddr.String(), ":")[0]
}

// func getIPLocally() string {
// 	addrs, err := net.InterfaceAddrs()
// 	if err != nil {
// 		fmt.Println(err)
// 		return ""
// 	}
// 	for _, address := range addrs {
// 		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
// 			if ipnet.IP.To4() != nil {
// 				return ipnet.IP.String()
// 			}
// 		}
// 	}
// 	return ""
// }

func GetPublicIP() string {

	eChan := make(chan string)
	ddl := make(chan bool)
	go func() {
		time.Sleep(time.Second * 2)
		ddl <- true
	}()
	go func() {
		eChan <- getIPByNet("http://ipinfo.io/ip")
	}()
	go func() {
		eChan <- getIPByNet("http://myexternalip.com/raw")
	}()
	select {
	case <-ddl:
		return getIPByNetUDP("8.8.8.8:53")
	case r := <-eChan:
		return r
	}

}

func GetIPListFromSubnet(subnet string) []string {
	prefix, err := netip.ParsePrefix(subnet)
	if err != nil {
		fmt.Println(err)
	}

	var ips []string
	for addr := prefix.Addr(); prefix.Contains(addr); addr = addr.Next() {
		ips = append(ips, StringConcat(addr.String(), "/32"))
	}
	if len(ips) < 2 {
		return ips
	}
	return ips[1 : len(ips)-1]
}

func Difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
