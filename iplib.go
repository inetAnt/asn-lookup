package main

import (
	"encoding/binary"
	"errors"
	"net"
)

func ReverseBits(b byte) (d byte) {
	d = b ^ 0xff
	return d
}

func Hostmask(mask []byte) []byte {
	hostmask := make([]byte, len(mask))
	for j, v := range mask {
		hostmask[j] = ReverseBits(v)
	}
	return hostmask
}

func IPSubnetHosts(firstIP *net.IP, hostmask []byte) (ips []net.IP) {
	// Returns all the IPs for a given range
	bs := make([]byte, 4) // bs will be our []byte
	mask := binary.BigEndian.Uint32(hostmask) + 1
	ips = make([]net.IP, mask)
	for i := uint32(0x0); i < mask; i++ {
		binary.BigEndian.PutUint32(bs, i)
		s, _ := AddBytesSlices(*firstIP, bs)
		ips[i] = net.IP{s[0], s[1], s[2], s[3]}
	}
	return ips
}

func AddBytesSlices(a []byte, b []byte) (c []byte, err error) {
	l := len(a)
	if l != len(b) {
		return nil, errors.New("Slices must be of the same length")
	}
	c = make([]byte, l)
	for i := 0; i < l; i++ {
		c[i] = a[i] + b[i]
	}
	return c, nil
}
