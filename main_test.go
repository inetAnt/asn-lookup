package main_test

import (
	"fmt"
	"github.com/inetant/asn-lookup"
	"net"
	"testing"
)

func Test_GetIPv4(t *testing.T) {
	ipv4 := net.IPv4(0x08, 0x08, 0x08, 0x08)
	v := main.GetIPVersion(&ipv4)
	ex := 4
	if v != ex {
		e := fmt.Sprintf("GetIPVersion(%v) should return %v", ipv4, ex)
		t.Error(e)
	}
}

func Test_GetIPv4Nope(t *testing.T) {
	ipv4 := net.ParseIP("256.257.258.259")
	v := main.GetIPVersion(&ipv4)
	ex := 0
	if v != ex {
		e := fmt.Sprintf("GetIPVersion(%v) should return %v", ipv4)
		t.Error(e)
	}
}

func Test_GetIPv4Short(t *testing.T) {
	ipv4 := net.ParseIP("1.2.3")
	v := main.GetIPVersion(&ipv4)
	ex := 0
	if v != ex {
		e := fmt.Sprintf("GetIPVersion(%v) should return %v", ipv4)
		t.Error(e)
	}
}
func Test_GetIPv6(t *testing.T) {
	ipv6 := net.ParseIP("dead:acab::1312")
	v := main.GetIPVersion(&ipv6)
	ex := 6
	if v != ex {
		e := fmt.Sprintf("GetIPVersion(%v) should return %v", ipv6, ex)
		t.Error(e)
	}
}

func Test_GetIPv6Nope(t *testing.T) {
	ipv6 := net.ParseIP("ipv6_lol_nope::1312")
	v := main.GetIPVersion(&ipv6)
	ex := 0
	if v != ex {
		e := fmt.Sprintf("GetIPVersion(%v) should return %v", ipv6, ex)
		t.Error(e)
	}
}
