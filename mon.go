package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type Person struct {
	Name  string
	Phone string
}

type ASN struct {
	Number int
	Name   string
}

type IPAddr struct {
	Address string
	Reverse string
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
	c.EnsureIndex(mgo.Index{
		Key:      []string{"number", "name"},
		Unique:   true,
		DropDups: true,
	})

	result := ASN{}
	err = c.Find(bson.M{"number": 15169}).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("ASN:", result.Name)
}
