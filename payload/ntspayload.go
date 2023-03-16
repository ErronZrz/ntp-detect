package payload

import (
	"active/utils"
	"bytes"
	"fmt"
)

const (
	keyLength = 32
)

type NTSPayload struct {
	Host       string
	Port       int
	CertDomain string
	Secure     bool
	Err        error
	Len        int
	RcvData    []byte
	C2SKey     []byte
	S2CKey     []byte
}

func (p *NTSPayload) Print() {
	if p.Err != nil {
		fmt.Println(p.Err)
	} else {
		fmt.Printf(p.Lines())
	}
}

func (p *NTSPayload) Lines() string {
	buf := new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("Remote address: %s:%d (%s)\n", p.Host, p.Port, utils.RegionOf(p.Host)))
	if p.Secure {
		buf.WriteString(fmt.Sprintf("Certificate domain: %s (valid)\n", p.CertDomain))
	} else {
		if p.CertDomain != "" {
			buf.WriteString(fmt.Sprintf("Certificate domain: %s (unverified)\n", p.CertDomain))
		} else {
			buf.WriteString(fmt.Sprintf("Certificate not found"))
		}
	}
	if len(p.C2SKey) == keyLength {
		buf.WriteString("C2S key: 0x")
		for _, b := range p.C2SKey {
			buf.WriteString(fmt.Sprintf("%02X", b))
		}
		buf.WriteByte('\n')
	}
	if len(p.S2CKey) == keyLength {
		buf.WriteString("S2C key: 0x")
		for _, b := range p.S2CKey {
			buf.WriteString(fmt.Sprintf("%02X", b))
		}
		buf.WriteByte('\n')
	}
	buf.WriteString(fmt.Sprintf("%d bytes received:\n", p.Len))
	rows := p.Len >> 4
	for i := 0; i < rows; i++ {
		for _, b := range p.RcvData[i<<4 : (i+1)<<4] {
			buf.WriteString(fmt.Sprintf("%02X ", b))
		}
		buf.WriteByte('\n')
	}
	if p.Len > rows<<4 {
		for _, b := range p.RcvData[rows<<4:] {
			buf.WriteString(fmt.Sprintf("%02X ", b))
		}
		buf.WriteByte('\n')
	}
	return buf.String()
}
