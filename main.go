package main

import (
	"fmt"
	"github.com/inetant/asn-lookup/whois"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net"
	"regexp"
	"strings"
	"sync"
	"time"
)

func worker(id int, jobs <-chan string, c *mgo.Collection) {
	for j := range jobs {
		fmt.Println("worker", id, "on ip", j)
		var name string
		names, err := net.LookupAddr(j)
		if err != nil {
			name = "null"
		} else {
			name = names[0]
		}
		UpdateReverse(c, j, name)
		time.Sleep(100 * time.Millisecond)
	}
}

func GetIPVersion(ip *net.IP) (version int) {
	if net.IP.To4(*ip) != nil {
		return 4
	} else if net.IP.To16(*ip) != nil {
		return 6
	} else {
		return 0
	}
}

func UpdateReverse(c *mgo.Collection, addr string, reverse string) {
	res := IPAddr{}
	err := c.Find(bson.M{"address": addr}).One(&res)
	if err != nil {
		fmt.Println(err.Error())
	}
	res.Reverse = reverse
	err = c.Update(IPAddr{Address: res.Address}, res)
	if err != nil {
		fmt.Println(err)
	}
}

func AddIPsJobsDB(as *ASN, networks []*net.IPNet, jobs chan<- string, c *mgo.Collection) {
	// analyse all our subnets, add to db with no reverse, add a job to our jobs queue
	for _, subnet := range networks {
		fmt.Println(subnet)
		hostmask := Hostmask(subnet.Mask)
		ips := IPSubnetHosts(&subnet.IP, hostmask)
		for _, v := range ips {
			addr := v.String() // IP address from the subnet
			jobs <- addr
			_ = c.Insert(IPAddr{
				AsnId:   as.Id,
				Address: addr,
				Reverse: "",
			})
		}
	}
}

func main() {

	s, c, _ := Collection("192.168.1.25")
	defer s.Close()

	as_data, _ := whois.GetInfo("whois.radb.net", 46489)
	_, err := c.Upsert(ASN{Number: as_data.Number}, as_data)
	if err != nil {
		fmt.Println(err)
	}

	as := ASN{}
	_ = c.Find(bson.M{"number": 46489}).One(&as)

	subnets, _ := whois.GetSubnets("whois.radb.net", as_data.Number)
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
	fmt.Println(networks)

	var wg sync.WaitGroup
	jobs := make(chan string, 1024)

	for w := 1; w <= 100; w++ {
		fmt.Println("go worker", w)
		go worker(w, jobs, c)
		wg.Add(1)
	}

	go AddIPsJobsDB(&as, networks, jobs, c)

	wg.Wait()

}
