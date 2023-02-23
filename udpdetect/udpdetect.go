package udpdetect

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"net"
	"sync"
	"time"
)

const (
	configPath       = "../config/"
	timeoutKey       = "detection.rcv_header.timeout"
	batchSizeKey     = "detection.send_udp.batch_size"
	defaultTimeout   = 3000
	defaultBatchSize = 256
)

var (
	timeout time.Duration
	data    = []byte{
		0xDB, 0x00, 0x04, 0xFA, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x00,
	}
)

func init() {
	viper.AddConfigPath(configPath)
	viper.SetConfigType("yaml")
	viper.SetConfigName("properties")
	viper.SetDefault(timeoutKey, defaultTimeout)
	viper.SetDefault(batchSizeKey, defaultBatchSize)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("error reading config file: %s", err)
	}
	var milli = time.Duration(viper.GetInt64(timeoutKey))
	if milli == 0 {
		milli = defaultTimeout
	}
	timeout = time.Millisecond * milli
}

func DialNetworkNTPWithBatchSize(cidr string, batchSize int) ([]*RcvPayload, error) {
	num, err := numOf(cidr)
	if err != nil {
		return nil, err
	}
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	dataCh := make(chan *RcvPayload, num)
	res := make([]*RcvPayload, 0, num)
	ctx, cancel := context.WithCancel(context.Background())
	go handleChan(ctx, dataCh, &res)
	wg := &sync.WaitGroup{}
	fmt.Printf("Ready to detect %d addresses\n", num)
	wg.Add(num)
	host := ip.Mask(ipNet.Mask)
	batchNum := num / batchSize
	for i := 0; i < batchNum; i++ {
		for j := 0; j < batchSize; j++ {
			hostStr := host.String()
			go writeToAddr(hostStr+":123", dataCh, wg)
			inc(host)
		}
		time.Sleep(timeout)
	}
	for ; ipNet.Contains(host); inc(host) {
		hostStr := host.String()
		go writeToAddr(hostStr+":123", dataCh, wg)
	}
	wg.Wait()
	cancel()
	close(dataCh)
	for {
		if payload, ok := <-dataCh; !ok {
			break
		} else {
			res = append(res, payload)
		}
	}
	return res, nil
}

func DialNetworkNTP(cidr string) ([]*RcvPayload, error) {
	return DialNetworkNTPWithBatchSize(cidr, viper.GetInt(batchSizeKey))
}

func writeToAddr(addr string, ch chan<- *RcvPayload, wg *sync.WaitGroup) {
	defer wg.Done()
	payload := &RcvPayload{host: addr[:len(addr)-4], port: 123}
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		payload.err = err
		ch <- payload
		return
	}
	// fmt.Println(udpAddr.String())
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		payload.err = err
		ch <- payload
		return
	}
	defer func() {
		_ = conn.Close()
	}()
	payload.sendTime = time.Now()
	_, err = conn.Write(data)
	if err != nil {
		payload.err = err
		ch <- payload
		return
	}
	buf := make([]byte, 128)
	err = conn.SetDeadline(time.Now().Add(timeout))
	if err != nil {
		payload.err = err
		ch <- payload
		return
	}
	n, _, err := conn.ReadFromUDP(buf)
	if err == nil && n > 0 {
		payload.rcvTime = time.Now()
		payload.len = n
		payload.rcvData = buf[:n]
		ch <- payload
	}
}

func handleChan(ctx context.Context, ch <-chan *RcvPayload, res *[]*RcvPayload) {
	for {
		select {
		case payload := <-ch:
			*res = append(*res, payload)
		case <-ctx.Done():
			break
		}
	}
}

func numOf(cidr string) (int, error) {
	n := len(cidr)
	pow := 32
	val := cidr[n-1]
	if val < 0x30 || val > 0x39 {
		return -1, errors.New("invalid CIDR address")
	}
	pow -= int(val - 0x30)
	val = cidr[n-2]
	if val == 0x2F {
		return 1 << pow, nil
	}
	if cidr[n-3] != 0x2F || val < 0x30 || val > 0x39 {
		return -1, errors.New("invalid CIDR address")
	}
	pow -= 10 * int(val-0x30)
	return 1 << pow, nil
}

func inc(ip []byte) {
	for i := 3; i >= 0; i-- {
		ip[i]++
		if ip[i] > 0 {
			break
		}
	}
}

func TryDialNTP(host string) {
	conn, err := net.DialTimeout("udp", host+":123", timeout)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(conn)
	fmt.Println("Over.")
	n, err := conn.Write(data)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(n)
	what := make([]byte, 48)
	err = conn.SetDeadline(time.Now().Add(timeout))
	if err != nil {
		fmt.Println(err.Error())
	}
	n, err = conn.Read(what)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(n, what)
}
