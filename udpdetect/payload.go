package udpdetect

import (
	"active/utils"
	"bytes"
	"fmt"
	"strconv"
	"time"
)

type RcvPayload struct {
	host     string
	port     int
	err      error
	len      int
	sendTime time.Time
	rcvTime  time.Time
	rcvData  []byte
}

func (p *RcvPayload) Error() error {
	return p.err
}

func (p *RcvPayload) Bytes() []byte {
	return p.rcvData
}

func (p *RcvPayload) Print() {
	if p.err != nil {
		fmt.Println(p.err)
	} else {
		fmt.Printf(p.Lines())
	}
}

func (p *RcvPayload) Lines() string {
	s := fmt.Sprintf("%d bytes received from %s:%d (%s):\n", p.len, p.host, p.port, utils.RegionOf(p.host))
	buf := bytes.NewBufferString(s)
	for i := 0; i < 3; i++ {
		for _, b := range p.rcvData[i<<4 : (i+1)<<4] {
			buf.WriteString(fmt.Sprintf("%02X ", b))
		}
		buf.WriteByte('\n')
	}
	// T2 - T1
	sendDelay := utils.CalculateDelay(p.rcvData[32:40], p.sendTime)
	// T4 - T3
	rcvDelay := -utils.CalculateDelay(p.rcvData[40:48], p.rcvTime)
	avgDelay := (sendDelay + rcvDelay) / 2
	offset := (sendDelay - rcvDelay) / 2
	buf.WriteString(fmt.Sprintf("Send delay:    %s\n", durationToStr(sendDelay)))
	buf.WriteString(fmt.Sprintf("Receive delay: %s\n", durationToStr(rcvDelay)))
	buf.WriteString(fmt.Sprintf("Average delay: %s\n", durationToStr(avgDelay)))
	buf.WriteString(fmt.Sprintf("Offset:        %s\n", durationToStr(offset)))
	return buf.String()
}

func durationToStr(d time.Duration) string {
	negative := d < 0
	us := d.Nanoseconds() / 1000
	str := strconv.FormatInt(us, 10)
	n := len(str)
	if n <= 3 || (negative && n <= 4) {
		return str + "μs"
	}
	return str[:n-3] + "." + str[n-3:] + "ms"
}
