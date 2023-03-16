package output

import (
	"active/nts"
	"active/parser"
	"fmt"
	"testing"
)

func TestWriteNTSToFile(t *testing.T) {
	host := "194.58.207.74"
	serverName := "sth2.nts.netnod.se"
	payload, err := nts.DialNTSKE(host, serverName, 0x0F)
	if err != nil {
		t.Error(err)
		return
	}

	if payload.Len == 0 {
		fmt.Println("empty response")
		return
	}

	res, err := parser.ParseNTSResponse(payload.RcvData)
	if err != nil {
		t.Error(err)
	} else {
		WriteNTSToFile(payload.Lines(), res.Lines(), host)
	}
}
