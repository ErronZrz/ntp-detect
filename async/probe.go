package async

import (
	"active/addr"
	"active/rcvpayload"
	"active/utils"
	"context"
	"fmt"
	"github.com/spf13/viper"
	"net"
	"time"
)

const (
	configPath           = "../resource/"
	localPortKey         = "async.local_port"
	checkIntervalKey     = "async.read.check_interval"
	timeoutKey           = "async.read.timeout"
	defaultLocalPort     = 11123
	defaultCheckInterval = 1000
	defaultTimeout       = 5000
)

var (
	checkInterval time.Duration
	timeout       time.Duration
	localPort     int
)

func init() {
	viper.AddConfigPath(configPath)
	viper.SetConfigType("yaml")
	viper.SetConfigName("properties")
	viper.SetDefault(localPortKey, defaultLocalPort)
	viper.SetDefault(checkIntervalKey, defaultCheckInterval)
	viper.SetDefault(timeoutKey, defaultTimeout)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("err reading resource file: %s", err)
	}
	localPort = viper.GetInt(localPortKey)
	checkInterval = time.Duration(viper.GetInt64(checkIntervalKey)) * time.Millisecond
	timeout = time.Duration(viper.GetInt64(timeoutKey)) * time.Millisecond
}

func DialNetworkNTP(cidr string) ([]*rcvpayload.RcvPayload, error) {
	errChan := make(chan error)
	done := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		select {
		case <-ctx.Done():
		default:
			cancel()
		}
	}()
	go func(ctx context.Context, errChan <-chan error) {
		for {
			select {
			case <-ctx.Done():
				return
			case err := <-errChan:
				fmt.Println(err)
			}
		}
	}(ctx, errChan)

	res := make([]*rcvpayload.RcvPayload, 0)
	go writeNetWorkNTP(cidr, done, errChan)
	go readNetworkNTP(ctx, cidr, &res, done)
	<-done
	<-time.After(timeout)
	cancel()
	<-done
	return res, nil
}

func writeNetWorkNTP(cidr string, done chan<- struct{}, errChan chan<- error) {
	defer func() {
		done <- struct{}{}
	}()
	generator, err := addr.NewAddrGenerator(cidr)
	if err != nil {
		errChan <- err
		return
	}
	for generator.HasNext() {
		probeNext(generator.NextHost(), errChan)
	}
}

func probeNext(host string, errChan chan<- error) {
	udpAddr, err := net.ResolveUDPAddr("udp", host+":123")
	if err != nil {
		errChan <- err
		return
	}

	conn, err := net.DialUDP("udp", &net.UDPAddr{Port: localPort}, udpAddr)
	if err != nil {
		errChan <- err
		return
	}

	defer func() {
		_ = conn.Close()
	}()

	_, err = conn.Write(utils.VariableData())
	if err != nil {
		errChan <- err
		return
	}
}
