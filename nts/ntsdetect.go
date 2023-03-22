package nts

import (
	"active/datastruct"
	"active/parser"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"time"
)

var (
	variableReq = []byte{
		0x80, 0x01, 0x00, 0x02, 0x00, 0x00, 0x80, 0x04, 0x00, 0x02, 0x00, 0x0F, 0x80, 0x00, 0x00, 0x00,
	}
)

func DetectNTSServer(host, serverName string) (*datastruct.NTSDetectPayload, error) {
	config := new(tls.Config)
	config.NextProtos = []string{alpnID}
	if serverName != "" {
		config.ServerName = serverName
	} else {
		config.InsecureSkipVerify = true
	}

	dialer := &net.Dialer{Timeout: timeout}

	conn, err := tls.DialWithDialer(dialer, "tcp", host+":4460", config)
	if err != nil {
		return nil, fmt.Errorf("cannot dial TLS server %s: %v", host, err)
	}

	info := datastruct.DetectInfo{
		AEADList:      make([]bool, 34),
		ServerPortSet: make(map[string]struct{}),
	}

	res := &datastruct.NTSDetectPayload{
		Host:   host,
		Port:   4460,
		Secure: !config.InsecureSkipVerify,
		Info:   info,
	}

	certs := conn.ConnectionState().PeerCertificates
	if len(certs) > 0 {
		res.CertDomain = certs[0].Subject.CommonName
	}

	err = singleReadWrite(0x0F, conn, info)
	if err != nil {
		return nil, err
	}

	if info.AEADList[0x0F] {
		fmt.Println("- (0F) AEAD_AES_SIV_CMAC_256:        supported")
	}

	for id := byte(0x01); id <= 0x21; id++ {
		if id == 0x0F {
			continue
		}

		<-time.After(haltTime)
		conn, err = tls.DialWithDialer(dialer, "tcp", host+":4460", config)
		if err != nil {
			return nil, fmt.Errorf("cannot dial TLS server %s: %v", host, err)
		}

		err = singleReadWrite(id, conn, info)
		if err != nil {
			return nil, err
		}

		name := datastruct.GetAEADName(id) + ":"
		status := "x"
		if info.AEADList[id] {
			status = "supported"
		}
		fmt.Printf("- (%02X) %-27s   %s\n", id, name, status)
	}
	fmt.Println()

	return res, nil
}

func singleReadWrite(aeadID byte, conn *tls.Conn, info datastruct.DetectInfo) error {
	defer func(conn *tls.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Printf("error closing TLS connect: %v", err)
		}
	}(conn)

	variableReq[11] = aeadID

	_, err := conn.Write(variableReq)
	if err != nil {
		return fmt.Errorf("send NTS-KE request failed: %v", err)
	}

	data, err := io.ReadAll(conn)
	if err != nil {
		return fmt.Errorf("read NTS-KE response failed: %v", err)
	}

	return parser.ParseDetectInfo(data, info)
}
