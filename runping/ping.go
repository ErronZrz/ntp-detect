package runping

import (
	"github.com/go-ping/ping"
	"time"
)

func RunPing(address string) string {
	p, err := ping.NewPinger(address)
	p.Count = 1
	p.Timeout = 2 * time.Second
	err = p.Run()
	if err != nil {
		return err.Error()
	}
	return p.Statistics().AvgRtt.String()
}
