package output

import (
	"active/async"
	"active/parser"
	"active/udpdetect"
	"fmt"
	"testing"
	"time"
)

func TestWriteToFile(t *testing.T) {
	cidr := "203.107.6.0/24"
	payloads, err := udpdetect.DialNetworkNTP(cidr)
	// fmt.Println(payloads)
	if err != nil {
		t.Error(err)
	}
	seqNum := 0
	now := time.Now()
	for _, p := range payloads {
		err := p.Err
		if err != nil {
			fmt.Println(err)
			continue
		}
		header, err := parser.ParseHeader(p.RcvData)
		if err != nil {
			t.Error(err)
		} else {
			seqNum++
			WriteToFile(p.Lines(), header.Lines(), "test timesync "+cidr, seqNum, p.RcvTime, now)
		}
	}
}

func TestAsyncWriteToFile(t *testing.T) {
	cidr := "203.107.6.0/24"
	payloads, err := async.DialNetworkNTP(cidr)
	if err != nil {
		t.Error(err)
	}
	seqNum := 0
	now := time.Now()
	for _, p := range payloads {
		err := p.Err
		if err != nil {
			fmt.Println(err)
			continue
		}
		header, err := parser.ParseHeader(p.RcvData)
		if err != nil {
			t.Error(err)
		} else {
			seqNum++
			WriteToFile(p.Lines(), header.Lines(), "test async "+cidr, seqNum, p.RcvTime, now)
		}
	}
}
