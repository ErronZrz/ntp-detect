package nts

import (
	"active/payload"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"time"
)

const (
	aesSivCmac256 = 0x0F
	alpnID        = "ntske/1"
	exportLabel   = "EXPORTER-network-time-security"
	keyLength     = 32
	timeout       = 5 * time.Second
)

var (
	reqBytes = []byte{
		0x80, 0x01, 0x00, 0x02, 0x00, 0x00, 0x80, 0x04, 0x00, 0x02, 0x00, 0x0F, 0x80, 0x00, 0x00, 0x00,
	}
)

func DialNTSKE(host, serverName string, aeadID byte) (*payload.NTSPayload, error) {
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
	defer func(conn *tls.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Printf("error closing TLS connect: %v", err)
		}
	}(conn)

	res := &payload.NTSPayload{
		Host:   host,
		Port:   4460,
		Secure: !config.InsecureSkipVerify,
	}

	state := conn.ConnectionState()

	certs := state.PeerCertificates
	if len(certs) > 0 {
		res.CertDomain = certs[0].Subject.CommonName
	}
	ctx := make([]byte, 4)
	ctx[3] = aeadID

	res.C2SKey, err = state.ExportKeyingMaterial(exportLabel, append(ctx, 0x00), keyLength)
	if err != nil {
		return nil, fmt.Errorf("export C2S key failed: %v", err)
	}
	res.S2CKey, err = state.ExportKeyingMaterial(exportLabel, append(ctx, 0x01), keyLength)
	if err != nil {
		return nil, fmt.Errorf("export S2C key failed: %v", err)
	}

	if aeadID > 0x00 && aeadID <= 0x21 {
		reqBytes[11] = aeadID
	} else {
		reqBytes[11] = aesSivCmac256
	}

	_, err = conn.Write(reqBytes)
	if err != nil {
		return nil, fmt.Errorf("send NTS-KE request failed: %v", err)
	}

	data, err := io.ReadAll(conn)
	if err != nil {
		return nil, fmt.Errorf("read NTS-KE response failed: %v", err)
	}

	res.Len = len(data)
	res.RcvData = data
	return res, nil
}
