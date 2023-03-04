package output

import (
	"active/parser"
	"active/udpdetect"
	"fmt"
	"testing"
	"time"
)

func TestWriteToFile(t *testing.T) {
	payloads, err := udpdetect.DialNetworkNTP("203.107.6.0/24")
	// fmt.Println(payloads)
	if err != nil {
		t.Error(err)
	}
	seqNum := 0
	now := time.Now()
	for _, p := range payloads {
		err := p.Error()
		if err != nil {
			fmt.Println(err)
			continue
		}
		header, err := parser.ParseHeader(p.Bytes())
		if err != nil {
			t.Error(err)
		} else {
			seqNum++
			WriteToFile(p.Lines(), header.Lines(), seqNum, p.RcvTime(), now)
		}
	}
}
