package addr

import (
	"errors"
	"net"
)

type Generator struct {
	nextIP net.IP
	ipNet  *net.IPNet
	total  int
	used   int
}

func NewAddrGenerator(cidr string) (*Generator, error) {
	num, err := numOf(cidr)
	if err != nil {
		return nil, err
	}
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	host := ip.Mask(ipNet.Mask)
	res := &Generator{
		nextIP: host,
		ipNet:  ipNet,
		total:  num,
		used:   0,
	}
	return res, nil
}

func (g *Generator) HasNext() bool {
	return g.used < g.total
}

func (g *Generator) NextHost() string {
	g.used++
	res := g.nextIP.String()
	inc(g.nextIP)
	return res
}

func (g *Generator) TotalNum() int {
	return g.total
}

func numOf(cidr string) (int, error) {
	n := len(cidr)
	pow := 32
	val := cidr[n-1]
	if val < 0x30 || val > 0x39 {
		return -1, errors.New("invalid CIDR address")
	}
	pow -= int(val - 0x30)
	val = cidr[n-2]
	if val == 0x2F {
		return 1 << pow, nil
	}
	if cidr[n-3] != 0x2F || val < 0x30 || val > 0x39 {
		return -1, errors.New("invalid CIDR address")
	}
	pow -= 10 * int(val-0x30)
	return 1 << pow, nil
}

func inc(ip []byte) {
	for i := 3; i >= 0; i-- {
		ip[i]++
		if ip[i] > 0 {
			break
		}
	}
}
