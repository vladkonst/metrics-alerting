package flags

import (
	"errors"
	"flag"
	"os"
	"strconv"
	"strings"
)

type NetAddress struct {
	Host string
	Port int
}

func (n NetAddress) String() string {
	return n.Host + ":" + strconv.Itoa(n.Port)
}

func (n *NetAddress) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		return errors.New("need address in a form host:port")
	}
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}
	n.Host = hp[0]
	n.Port = port
	return nil
}

func GetNetAddress() *NetAddress {
	addr := &NetAddress{Host: "localhost", Port: 8080}
	flag.Var(addr, "a", "Server net address host:port")
	flag.Parse()
	if adr := os.Getenv("ADDRESS"); adr != "" {
		addr.Set(os.Getenv("ADDRESS"))
	}
	return addr
}
