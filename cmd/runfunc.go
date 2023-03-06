package cmd

import (
	"active/async"
	"active/output"
	"active/parser"
	"active/rcvpayload"
	"active/udpdetect"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
	"time"
)

func executeTimeSync(cmd *cobra.Command, args []string) (string, error) {
	if nPrintedHosts > 128 {
		nPrintedHosts = 128
	}
	cmdName := cmd.Name()
	if args == nil || len(args) == 0 {
		return "", errors.New("command `timesync` missing arguments")
	}
	address := args[0]
	var ngStr string
	if nGoroutines <= 0 {
		ngStr = "auto"
	} else {
		ngStr = strconv.Itoa(nGoroutines)
	}
	_, _ = fmt.Fprintf(os.Stdout, "Ready to run %s.\n    address: %s\n    num of goroutines: %s\n"+
		"    num of printed hosts: %d\n\n\n\n", cmdName, address, ngStr, nPrintedHosts)
	var payloads []*rcvpayload.RcvPayload
	var err error
	if nGoroutines <= 0 {
		payloads, err = udpdetect.DialNetworkNTP(address)
	} else {
		payloads, err = udpdetect.DialNetworkNTPWithBatchSize(address, nGoroutines)
	}
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
	}
	return generateResult(payloads), nil
}

func executeAsync(cmd *cobra.Command, args []string) (string, error) {
	if nPrintedHosts > 128 {
		nPrintedHosts = 128
	}
	cmdName := cmd.Name()
	if args == nil || len(args) == 0 {
		return "", errors.New("command `async` missing arguments")
	}
	address := args[0]
	_, _ = fmt.Fprintf(os.Stdout, "Ready to run %s.\n    address: %s\n    "+
		"num of printed hosts: %d\n\n\n\n", cmdName, address, nPrintedHosts)
	payloads, err := async.DialNetworkNTP(address)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
	}
	return generateResult(payloads), nil
}

func generateResult(payloads []*rcvpayload.RcvPayload) string {
	seqNum := 0
	now := time.Now()
	var builder strings.Builder

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

	return builder.String()
}
