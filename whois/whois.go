package whois

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"regexp"
)

type WhoisClient struct {
	Conn *net.TCPConn
}

type AutonomousSystem struct {
	Number int
	Name   string
	Descr  string
}

func Connect(server string) (w *WhoisClient, err error) {
	// Creates the TCP connection to whois server, returns the pointer
	addr, err := net.ResolveTCPAddr("tcp6", server+":43")
	if err != nil {
		fmt.Println(err.Error())
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Println(err.Error())
	}

	whois := WhoisClient{conn}
	return &whois, nil
}

func GetInfo(server string, asn int) (as *AutonomousSystem, err error) {
	w, _ := Connect(server)
	defer w.Conn.Close()
	req := fmt.Sprintf("AS%v\r\n", asn)
	_, err = w.Conn.Write([]byte(req))
	if err != nil {
		fmt.Println(err.Error())
	}

	buf := &bytes.Buffer{}
	for {
		reply := make([]byte, 1514)
		n, err := w.Conn.Read(reply)
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		buf.Write(reply[:n])
		if reply[0] == '\r' && reply[1] == '\n' {
			break
		}
	}

	res := string(buf.String())

	r_as_name := regexp.MustCompile("as-name:\\s+(\\S+)")
	r_descr := regexp.MustCompile("descr:\\s+(\\w+.*)")

	as_name := r_as_name.FindStringSubmatch(res)
	descr := r_descr.FindStringSubmatch(res)

	if as_name == nil {
		panic("asname empty, will not work")
	}

	if descr == nil {
		descr = as_name // dirty but...
	}

	as = &AutonomousSystem{asn, as_name[1], descr[1]}

	return as, nil
}

func GetSubnets(server string, asn int) (response string, err error) {
	w, _ := Connect(server)
	defer w.Conn.Close()
	req := fmt.Sprintf("-i origin AS%v\r\n", asn)
	_, err = w.Conn.Write([]byte(req))
	if err != nil {
		fmt.Println(err.Error())
	}

	buf := &bytes.Buffer{}
	for {
		reply := make([]byte, 1514)
		n, err := w.Conn.Read(reply)
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		buf.Write(reply[:n])
		if reply[0] == '\r' && reply[1] == '\n' {
			break
		}
	}

	return string(buf.String()), nil
}
