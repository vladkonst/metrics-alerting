package configs

import (
	"errors"
	"net"
	"strconv"
	"strings"
)

type NetAddressCfg struct {
	Host string
	Port int
}

func (n NetAddressCfg) String() string {
	return n.Host + ":" + strconv.Itoa(n.Port)
}

func (n *NetAddressCfg) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		return errors.New("need address in a form host:port")
	}
	var ip string
	if hp[0] == "localhost" || hp[0] == "" {
		ip = hp[0]
	} else {
		ip = net.ParseIP(hp[0]).String()
		if ip == "<nil>" {

			return errors.New("can't validate provided IP")
		}
	}
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}
	if port < 0 || port > 65535 {
		return errors.New("can't validate provided port number")
	}

	n.Host = ip
	n.Port = port
	return nil
}
