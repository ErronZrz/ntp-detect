package cmd

import (
	"active/output"
	"active/parser"
	"active/rcvpayload"
	"active/udpdetect"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"strings"
	"time"
)

func executeTimeSync(cmd *cobra.Command, args []string) (string, error) {
	cmdName := cmd.Name()
	if args == nil || len(args) == 0 {
		return "", errors.New("command `timesync` missing arguments")
	}
	address := args[0]
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Ready to run %s.\n    address: %s\n    num of goroutines: %d\n"+
		"    num of printed hosts: %d\n\n\n\n", cmdName, address, nGoroutines, nPrintedHosts))
	var payloads []*rcvpayload.RcvPayload
	var err error
	if nGoroutines <= 0 {
		payloads, err = udpdetect.DialNetworkNTP(address)
	} else {
		payloads, err = udpdetect.DialNetworkNTPWithBatchSize(address, nGoroutines)
	}
	if err != nil {
		builder.WriteString("error: " + err.Error())
		return builder.String(), err
	}
	seqNum := 0
	now := time.Now()
	for _, p := range payloads {
		err := p.Err
		if err != nil {
			builder.WriteString(err.Error())
			continue
		}
		header, err := parser.ParseHeader(p.RcvData)
		if err != nil {
			builder.WriteString(err.Error())
		} else {
			seqNum++
			payloadStr, headerStr := p.Lines(), header.Lines()
			output.WriteToFile(payloadStr, headerStr, seqNum, p.RcvTime, now)
			if seqNum <= nPrintedHosts {
				builder.WriteString(fmt.Sprintf("[Host %d]\n", seqNum))
				builder.WriteString(payloadStr)
				builder.WriteString("[parsed]\n")
				builder.WriteString(headerStr)
			}
		}
	}
	return builder.String(), nil
}
