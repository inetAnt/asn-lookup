package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ASN struct {
	Id     bson.ObjectId `bson:"_id,omitempty"`
	Number int           `bson:"number"`
	Name   string        `bson:"name"`
	Descr  string        `bson:"descr"`
}

type IPAddr struct {
	AsnId   bson.ObjectId `bson:"asn_id,omitempty"`
	Address string        `bson:"address"`
	Reverse string        `bson:"reverse"`
}

func Collection(server string) (session *mgo.Session, c *mgo.Collection, err error) {
	session, err = mgo.Dial(server)
	if err != nil {
		panic(err)
	}

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c = session.DB("asn-lookup").C("asns")
	c.EnsureIndex(mgo.Index{
		Key:      []string{"number"},
		Unique:   true,
		DropDups: true,
		Sparse:   true,
	})
	c.EnsureIndex(mgo.Index{
		Key:      []string{"address"},
		Unique:   true,
		DropDups: true,
		Sparse:   true,
	})
	return session, c, err
}

/*
	result := ASN{}
	err = c.Find(bson.M{"number": 15169}).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("ASN:", result.Name)
*/
