package async

import (
	"active/parser"
	"active/rcvpayload"
	"active/utils"
	"context"
	"errors"
	"fmt"
	"net"
	"time"
)

func readNetworkNTP(ctx context.Context, cidr string, payloads *[]*rcvpayload.RcvPayload, done chan<- struct{}) {
	defer func() {
		done <- struct{}{}
	}()
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: localPort})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	buf := make([]byte, 128)

	for {
		select {
		case <-ctx.Done():
			// fmt.Println("Done!")
			return
		default:
			err := conn.SetReadDeadline(time.Now().Add(checkInterval))
			if err != nil {
				fmt.Println(err)
				continue
			}
			n, udpAddr, err := conn.ReadFromUDP(buf)
			if err != nil {
				continue
			}
			if !ipNet.Contains(udpAddr.IP) {
				fmt.Println("IP out of range: " + udpAddr.IP.String())
				continue
			}
			payload := &rcvpayload.RcvPayload{
				Host:    udpAddr.IP.String(),
				Port:    udpAddr.Port,
				Len:     n,
				RcvTime: time.Now(),
				RcvData: buf[:n],
			}
			if n != parser.HeaderLength {
				payload.Err = errors.New(fmt.Sprintf("header length %d doesn't equal to 48", n))
			} else {
				payload.SendTime = utils.ConvertTimestamp(buf[24:32])
			}
			*payloads = append(*payloads, payload)
		}
	}
}
