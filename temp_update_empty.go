package main

import (
	"fmt"
	"gopkg.in/mgo.v2"

	"gopkg.in/mgo.v2/bson"
	"net"
	"time"
)

func worker(id int, jobs <-chan string, results chan<- IPAddr) {
	for j := range jobs {
		fmt.Println("worker", id, "on ip", j)
		name := "null"
		names, err := net.LookupAddr(j)
		if err != nil {
			name = "null"
		} else {
			name = names[0]
		}
		results <- IPAddr{Address: j, Reverse: name}
		time.Sleep(100 * time.Millisecond)
	}
}

func UpdateDocReverse(c *mgo.Collection, addr string, reverse string) {
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

func main() {

	s, c, _ := Collection("192.168.1.25")
	defer s.Close()

	var ips []IPAddr
	_ = c.Find(bson.M{"reverse": ""}).All(&ips)
	fmt.Println(len(ips))

	jobs := make(chan string)
	results := make(chan IPAddr)

	for w := 1; w <= 50; w++ {
		go worker(w, jobs, results)
	}

	go func() {
		for result := range results {
			UpdateDocReverse(c, result.Address, result.Reverse)
			fmt.Println("updated", result.Address, "put", result.Reverse)
		}
	}()

	go func() {
		for _, ip := range ips {
			jobs <- ip.Address
			fmt.Println("Added", ip.Address, "to jobs")
		}
	}()

	fmt.Println("Loop ended")
	fmt.Println(len(jobs))

	go func() {
		time.Sleep(time.Second)
	}()
	select {}
}
