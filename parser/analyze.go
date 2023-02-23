package parser

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

func Analyze(p *Info) {
	p.Leap = analyzeLeap(p.Leap)
	p.Precision = analyzeLog2s(p.Precision)
	p.RootDelay = analyzeMs(p.RootDelay)
	p.RootDispersion = analyzeMs(p.RootDispersion)
	p.ReferenceID = analyzeRefID(p.ReferenceID)
	p.ReferenceTime = analyzeTimeStamp(p.ReferenceTime)
	p.Clock = analyzeTimeStamp(p.Clock)
	p.TimeConstant = analyzeLog2s(p.TimeConstant)
	p.MinTimeConstant = analyzeLog2s(p.MinTimeConstant)
	p.Offset = analyzePPM(p.Offset)
	p.Frequency = analyzePPM(p.Frequency)
	p.TAI = analyzeSecond(p.TAI)
	p.LeapSecond = analyzeSecStr(p.LeapSecond)
	p.Expire = analyzeSecStr(p.Expire)
}

func analyzeLeap(s string) string {
	if s == "" {
		return ""
	}
	switch s[0] {
	case '0':
		return "No warning"
	case '1':
		return "Last minute of the day has 61 seconds"
	case '2':
		return "Last minute of the day has 59 seconds"
	}
	return "Unknown or clock not synchronized"
}

func analyzeLog2s(s string) string {
	if s == "" {
		return ""
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return s
	}
	seconds := math.Pow(2, float64(i))
	if seconds > 0.1 && seconds < 10000 {
		return fmt.Sprintf("%s (%.3fs)", s, seconds)
	}
	return fmt.Sprintf("%s (%.3e s)", s, seconds)
}

func analyzeMs(s string) string {
	if s == "" {
		return ""
	}
	return s + "ms"
}

func analyzeRefID(s string) string {
	// TODO: identify country
	if s == "" {
		return ""
	}
	return s
}

func analyzeTimeStamp(s string) string {
	if s == "" {
		return ""
	}
	if len(s) != 19 {
		return s
	}
	s1, s2 := s[2:10], s[11:19]
	i1, err1 := strconv.ParseInt(s1, 16, 64)
	i2, err2 := strconv.ParseInt(s2, 16, 64)
	if err1 != nil || err2 != nil {
		return s
	}
	t := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	t = t.Add(time.Duration(i1) * time.Second)
	format1 := t.Format("2006-01-02 15:04:05")
	f2 := float64(i2) / float64(int64(1)<<32)
	format2 := strconv.FormatFloat(f2, 'f', 3, 64)
	format2 = format2[1:]
	return fmt.Sprintf("%s (%s%s UTC)", s, format1, format2)
}

func analyzePPM(s string) string {
	if s == "" {
		return ""
	}
	return s + "Part/Million"
}

func analyzeSecond(s string) string {
	if s == "" {
		return ""
	}
	return s + "s"
}

func analyzeSecStr(s string) string {
	if s == "" {
		return ""
	}
	t, err := time.Parse("200601021504", s)
	if err != nil {
		return s
	}
	return t.Format("2006-01-02 15:04:05 UTC")
}
