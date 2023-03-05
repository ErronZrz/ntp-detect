package utils

import (
	"encoding/binary"
	"fmt"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/spf13/viper"
	"math"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	configPath    = "../resource/"
	dbPathKey     = "ip2region.db_path"
	unknownFlag   = "未知地区"
	privateFlag   = "内网地址"
	preciseFormat = "2006-01-02 15:04:05.000000 UTC"
)

var (
	startingPoint = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	searcher      *xdb.Searcher
	fixedData     []byte
	variableData  []byte
)

func init() {
	fixedData = []byte{
		0xDB, 0x00, 0x04, 0xFA, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x00,
	}
	variableData = []byte{
		0xDB, 0x00, 0x04, 0xFA, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x00,
	}
	viper.AddConfigPath(configPath)
	viper.SetConfigType("yaml")
	viper.SetConfigName("properties")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("error reading resource file: %s", err)
	}
	filePath := viper.GetString(dbPathKey)
	buf, err := xdb.LoadContentFromFile(filePath)
	if err != nil {
		fmt.Printf("failed to load content: %s", err)
	}
	searcher, err = xdb.NewWithBuffer(buf)
}

func FixedData() []byte {
	return fixedData
}

func FromInt8(i int8) string {
	val := math.Pow(2, float64(i))
	scientific := FormatScientific(val)
	return fmt.Sprintf("2^%d (%s) sec", i, scientific)
}

func FormatScientific(f float64) string {
	if f == 0 {
		return "0"
	}
	if f >= 0.001 && f <= 1000 {
		return strconv.FormatFloat(f, 'f', 3, 64)
	}
	exp := int(math.Floor(math.Log10(f)))
	mantissa := f / math.Pow10(exp)
	return fmt.Sprintf("%.3fe%d", mantissa, exp)
}

func RegionOf(ipStr string) string {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return unknownFlag
	}
	if ip.IsPrivate() {
		return privateFlag
	}
	region, err := searcher.SearchByStr(ipStr)
	if err != nil {
		fmt.Println(err)
		return unknownFlag
	}
	parts := strings.Split(region, "|")
	country := parts[0]
	if country == "0" {
		return unknownFlag
	}
	if country != "中国" || parts[2] == "0" {
		return country
	}
	if strings.HasPrefix(parts[3], parts[2]) {
		return parts[2]
	}
	res := strings.ReplaceAll(parts[2], "省", "")
	if parts[3] == "0" {
		return res
	}
	return res + strings.ReplaceAll(parts[3], "市", "")
}

func CalculateDelay(timestamp []byte, another time.Time) time.Duration {
	t := binary.BigEndian.Uint64(timestamp)
	seconds := int64(t >> 32)
	nano := (int64(t&0xFFFF_FFFF) * int64(time.Second)) >> 32
	d := startingPoint.Add(time.Duration(seconds) * time.Second).Add(time.Duration(nano))
	delay := d.Sub(another)
	return delay
}

func FormatTimestamp(timestamp []byte) string {
	return ConvertTimestamp(timestamp).Format(preciseFormat)
}

func ConvertTimestamp(timestamp []byte) time.Time {
	intPart := binary.BigEndian.Uint32(timestamp[:4])
	fracPart := binary.BigEndian.Uint32(timestamp[4:])
	intTime := startingPoint.Add(time.Duration(intPart) * time.Second)
	fracDuration := (time.Duration(fracPart) * time.Second) >> 32
	return intTime.Add(fracDuration)
}

func VariableData() []byte {
	d := time.Now().Sub(startingPoint)
	seconds := d / time.Second
	high32 := seconds << 32
	nano := d - seconds*time.Second
	low32 := (nano << 32) / time.Second
	binary.BigEndian.PutUint64(variableData[40:], uint64(high32|low32))
	return variableData
}
