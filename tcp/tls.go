package tcp

import (
	"crypto/tls"
	"fmt"
	"net"
	"strconv"
)

func IsTLSEnabled(ip string, port int, serverName string) bool {
	conn, err := net.Dial("tcp", ip+":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(conn)

	config := new(tls.Config)
	if serverName != "" {
		config.ServerName = serverName
	} else {
		config.InsecureSkipVerify = true
	}
	tlsConn := tls.Client(conn, config)
	err = tlsConn.Handshake()
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer func(tlsConn *tls.Conn) {
		err := tlsConn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(tlsConn)
	return true
}
