package parser

import (
	"active/tcp"
	"fmt"
	"testing"
)

var (
	reqBytes = []byte{
		0x80, 0x01, 0x00, 0x02, 0x00, 0x00, 0x80, 0x04, 0x00, 0x02, 0x00, 0x0F, 0x80, 0x00, 0x00, 0x00,
	}
)

func TestParseNTSResponse(t *testing.T) {
	resBytes, err := tcp.WriteReadTLS("194.58.207.74", 4460, "sth2.nts.netnod.se", reqBytes)
	if err != nil {
		t.Error(err)
	}
	n := len(resBytes)
	if n == 0 {
		fmt.Println("empty response")
		return
	}
	fmt.Printf("%d bytes received\n", n)
	res, err := ParseNTSResponse(resBytes)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Print(res.Lines())
	}
}
