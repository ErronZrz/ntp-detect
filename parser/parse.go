package parser

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	spaceByte byte = 0x20
	quoteByte byte = 0x22
	commaByte byte = 0x2C
	equalByte byte = 0x3D
	crByte    byte = 0x0D
	lfByte    byte = 0x0A
)

type Info struct {
	Version         string `json:"version,omitempty"`
	Processor       string `json:"processor,omitempty"`
	System          string `json:"system,omitempty"`
	Leap            string `json:"leap,omitempty"`
	Stratum         string `json:"stratum,omitempty"`
	Precision       string `json:"precision,omitempty"`
	RootDelay       string `json:"rootdelay,omitempty"`
	RootDispersion  string `json:"rootdisp,omitempty"`
	ReferenceID     string `json:"refid,omitempty"`
	ReferenceTime   string `json:"reftime,omitempty"`
	Clock           string `json:"clock,omitempty"`
	Peer            string `json:"peer,omitempty"`
	TimeConstant    string `json:"tc,omitempty"`
	MinTimeConstant string `json:"mintc,omitempty"`
	Offset          string `json:"offset,omitempty"`
	Frequency       string `json:"frequency,omitempty"`
	SystemJitter    string `json:"sys_jitter,omitempty"`
	ClockJitter     string `json:"clk_jitter,omitempty"`
	ClockWander     string `json:"clk_wander,omitempty"`
	TAI             string `json:"tai,omitempty"`
	LeapSecond      string `json:"leapsec,omitempty"`
	Expire          string `json:"expire,omitempty"`
}

func ParseInfo(s string) (*Info, error) {
	buffer := bytes.NewBuffer([]byte{'{'})
	cur, n := 0, len(s)
	for cur < n {
		next := cur
		for next < n && s[next] != equalByte {
			next++
		}
		if next == n {
			return nil, errors.New(fmt.Sprintf("missing equal sign, parsed data: %s", buffer.Bytes()))
		}
		buffer.WriteByte(quoteByte)
		buffer.WriteString(s[cur:next])
		buffer.WriteString("\": \"")
		next++
		cur = next
		for next < n {
			if s[next] == commaByte {
				skip := 0
				if s[next+1] == spaceByte || s[next+1] == lfByte {
					skip = 2
				} else if s[next+1] == crByte && s[next+2] == lfByte {
					skip = 3
				}
				if skip > 0 {
					if s[cur] == quoteByte {
						buffer.WriteString(s[cur+1 : next-1])
					} else {
						buffer.WriteString(s[cur:next])
					}
					buffer.WriteString("\", ")
					cur = next + skip
					break
				}
			}
			next++
		}
		if next == n {
			buffer.WriteString(s[cur:n])
			buffer.WriteString("\"}")
			break
		}
	}
	jsonBytes := buffer.Bytes()
	// fmt.Println(string(jsonBytes))
	res := new(Info)
	err := json.Unmarshal(jsonBytes, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

//func main() {
//	s := "version=\"ntpd 4.2.8p12@1.3728-o (1)\", processor=\"x86_64\",\r\nsystem=\"Linux/5.15.0-56-generic\", leap=0, stratum=4, precision=-25,\r\nrootdelay=271.993, rootdisp=44.018, refid=78.46.102.180,\r\nreftime=0xe73d7391.a5b76a1a, clock=0xe73d73cc.7224675f, peer=58491,\r\ntc=7, mintc=3, offset=-8.087560, frequency=-0.703, sys_jitter=11.094963,\r\nclk_jitter=5.402, clk_wander=1.956, tai=37, leapsec=201701010000,\r\nexpire=202306280000"
//	res, err := ParseInfo(s)
//	if err != nil {
//		fmt.Println(err.Error())
//	} else {
//		fmt.Printf("%+v", *res)
//		Analyze(res)
//		fmt.Printf("%+v", *res)
//	}
//}
