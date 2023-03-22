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
	payload, err := nts.DialNTSKE(host, serverName, 0x22)
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

func TestWriteNTSDetectToFile(t *testing.T) {
	host := "194.58.207.74"
	serverName := "sth2.nts.netnod.se"
	payload, err := nts.DetectNTSServer(host, serverName)
	if err != nil {
		t.Error(err)
		return
	}

	WriteNTSDetectToFile(payload.Lines(), host)
}
