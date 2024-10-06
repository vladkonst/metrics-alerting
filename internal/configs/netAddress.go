package configs

import (
	"errors"
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
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}
	n.Host = hp[0]
	n.Port = port
	return nil
}
