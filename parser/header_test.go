package parser

import (
	"active/udpdetect"
	"fmt"
	"testing"
)

func TestParseHeader(t *testing.T) {
	payloads, err := udpdetect.DialNetworkNTP("203.107.6.0/24")
	if err != nil {
		t.Error(err)
	}
	for _, p := range payloads {
		err := p.Err
		if err != nil {
			fmt.Println(err)
			continue
		}
		data := p.RcvData
		p.Print()
		header, err := ParseHeader(data)
		if err != nil {
			t.Error(err)
		} else {
			header.Print()
		}
	}
}
