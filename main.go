package main

import (
	"fmt"
	"github.com/inetant/asn-lookup/whois"
	"gopkg.in/mgo.v2"
	"net"
	"regexp"
	"strings"
)

func GetIPVersion(ip *net.IP) (version int) {
	if net.IP.To4(*ip) != nil {
		return 4
	} else if net.IP.To16(*ip) != nil {
		return 6
	} else {
		return 0
	}
}

func main() {
	session, err := mgo.Dial("192.168.1.25")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("asn-lookup").C("asns")

	as, _ := whois.GetInfo("whois.radb.net", 16276)
	_ = c.Insert(as)
	fmt.Println(as.Name)

	subnets, _ := whois.GetSubnets("whois.radb.net", as.Number)
	lines := strings.Split(subnets, "\n")

	networks := make([]*net.IPNet, 0) // List of all the networks (CIDRs) owned
	for _, line := range lines {
		prefix := regexp.MustCompile("^route(?:6)?:\\s+(\\S+)")
		m := prefix.FindStringSubmatch(line)
		if m != nil {
			ip, subnet, err := net.ParseCIDR(m[1])
			if err != nil {
				fmt.Printf("net did not find any CIDR in %v\n", m[1])
			}
			if GetIPVersion(&ip) == 4 {
				networks = append(networks, subnet)
			}
		}
	}

	for _, subnet := range networks {
		hostmask := Hostmask(subnet.Mask)
		ips := IPSubnetHosts(&subnet.IP, hostmask)
		for _, v := range ips {
			fmt.Println(v)
		}
		break
	}

}
