package nts

import (
	"crypto/tls"
	"fmt"
	"net"
	"strconv"
)

func DialNTSKE(host string, port int, serverName string) (*KeyExchange, error) {
	config := new(tls.Config)
	if serverName != "" {
		config.ServerName = serverName
	} else {
		config.InsecureSkipVerify = true
	}

	keConn, err := Connect(net.JoinHostPort(host, strconv.Itoa(port)), config, true)
	if err != nil {
		return nil, fmt.Errorf("connection failure to %s: %v", serverName, err)
	}

	err = keConn.Exchange()
	if err != nil {
		return nil, fmt.Errorf("NTS-KE exchange error: %v", err)
	}

	if len(keConn.Meta.Cookie) == 0 {
		return nil, fmt.Errorf("received no cookies")
	}

	if keConn.Meta.Algo != aesSivCmac256 {
		return nil, fmt.Errorf("unknown algorighm in NTS-KE")
	}

	err = keConn.ExportKeys()
	if err != nil {
		return nil, fmt.Errorf("export key failed: %v", err)
	}

	return keConn, nil
}
