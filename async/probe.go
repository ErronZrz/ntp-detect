package async

import (
	"active/addr"
	"active/datastruct"
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
	haltTimeKey          = "async.send.halt_time"
	defaultLocalPort     = 11123
	defaultCheckInterval = 1000
	defaultTimeout       = 5000
	defaultHaltTime      = 0
)

var (
	checkInterval time.Duration
	timeout       time.Duration
	haltTime      time.Duration
	localPort     int
	sharedConn    *net.UDPConn
	doneCh        chan struct{}
	localAddr     *net.UDPAddr
)

func init() {
	viper.AddConfigPath(configPath)
	viper.SetConfigType("yaml")
	viper.SetConfigName("properties")
	viper.SetDefault(localPortKey, defaultLocalPort)
	viper.SetDefault(checkIntervalKey, defaultCheckInterval)
	viper.SetDefault(timeoutKey, defaultTimeout)
	viper.SetDefault(haltTimeKey, defaultHaltTime)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("err reading resource file: %v", err)
		return
	}
	localPort = viper.GetInt(localPortKey)
	checkInterval = time.Duration(viper.GetInt64(checkIntervalKey)) * time.Millisecond
	timeout = time.Duration(viper.GetInt64(timeoutKey)) * time.Millisecond
	haltTime = time.Duration(viper.GetInt64(haltTimeKey)) * time.Millisecond

	localAddr = &net.UDPAddr{Port: localPort}
}

func DialNetworkNTP(cidr string) <-chan *datastruct.RcvPayload {
	errChan := make(chan error)
	doneCh = make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
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

	dataCh := make(chan *datastruct.RcvPayload, 1024)

	var err error
	sharedConn, err = net.ListenUDP("udp", localAddr)
	if err != nil {
		errChan <- err
	}

	go writeNetWorkNTP(cidr, errChan)
	go readNetworkNTP(ctx, cidr, dataCh)

	go func() {
		<-doneCh
		<-time.After(timeout)
		cancel()
		<-doneCh
		close(dataCh)
		<-time.After(time.Second)
		_ = sharedConn.Close()
	}()

	return dataCh
}

func writeNetWorkNTP(cidr string, errChan chan<- error) {
	defer func() {
		doneCh <- struct{}{}
	}()

	generator, err := addr.NewAddrGenerator(cidr)
	if err != nil {
		errChan <- err
		return
	}
	for generator.HasNext() {
		probeNext(generator.NextHost(), errChan)
		if haltTime > 0 {
			<-time.After(haltTime)
		}
	}
}

func probeNext(host string, errChan chan<- error) {
	remoteAddr, err := net.ResolveUDPAddr("udp", host+":123")
	if err != nil {
		errChan <- err
		return
	}

	_, err = sharedConn.WriteToUDP(utils.VariableData(), remoteAddr)
	if err != nil {
		errChan <- err
		return
	}
}
