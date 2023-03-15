package nts

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	NTPv4ID        uint16 = 0
	criticalBitPos        = 15
)

// NTS-KE record types
const (
	RecEOM uint16 = iota
	RecNextProto
	RecError
	RecWarning
	RecAEAD
	RecCookie
	RecServer
	RecPort
)

// Record is the interface all record types must implement. Header()
// returns the record header. string() returns a printable
// representation of the record type. pack() packs the record into
// wire format.
type Record interface {
	Header() RecordHdr

	string() string
	pack(*bytes.Buffer) error
}

// RecordHdr is the header on all records send in NTS-KE. The first
// bit of the Type is the critical bit.
type RecordHdr struct {
	Type    uint16 // First bit is critical bit
	BodyLen uint16
}

func (h RecordHdr) pack(buf *bytes.Buffer) error {
	err := binary.Write(buf, binary.BigEndian, h)
	return err
}

func (h RecordHdr) Header() RecordHdr {
	return h
}

// NextProto record. Tells the other side we want to speak NTP,
// probably. Set to 0.
type NextProto struct {
	RecordHdr
	NextProto uint16
}

func (n NextProto) string() string {
	return fmt.Sprintf("--NextProto: %v\n", n.NextProto)
}

func (n NextProto) pack(buf *bytes.Buffer) error {
	value := new(bytes.Buffer)
	err := binary.Write(value, binary.BigEndian, n.NextProto)
	if err != nil {
		return err
	}

	n.RecordHdr.Type = RecNextProto
	n.RecordHdr.Type = setBit(n.RecordHdr.Type, criticalBitPos)
	n.RecordHdr.BodyLen = uint16(value.Len())

	err = n.RecordHdr.pack(buf)
	if err != nil {
		return err
	}

	_, err = buf.ReadFrom(value)
	if err != nil {
		return err
	}

	return nil
}

// End is the End of Message record.
type End struct {
	RecordHdr
}

func (e End) pack(buf *bytes.Buffer) error {
	return packHeader(RecEOM, true, buf, 0)
}

func (e End) string() string {
	return fmt.Sprintf("--EOM\n")
}

// Server is the NTP Server record, telling the client to use a
// certain server for the next protocol query. Critical bit is
// optional. Set Critical to true if you want it set.
type Server struct {
	RecordHdr
	Addr     []byte
	Critical bool
}

func (s Server) pack(buf *bytes.Buffer) error {
	return packSimple(RecServer, s.Critical, s.Addr, buf)
}

func (s Server) string() string {
	return fmt.Sprintf("--Server: %v\n", string(s.Addr))
}

// Port is the NTP Port record, telling the client to use this port
// for the next protocol query. Critical bit is optional. Set Critical
// to true if you want it set.
type Port struct {
	RecordHdr
	Port     uint16
	Critical bool
}

func (p Port) pack(buf *bytes.Buffer) error {
	return packSimple(RecPort, p.Critical, p.Port, buf)
}

func (p Port) string() string {
	return fmt.Sprintf("--Port: %v\n", p.Port)
}

// Cookie is an NTS cookie to be used when querying time over NTS.
type Cookie struct {
	RecordHdr
	Cookie []byte
}

func (c Cookie) pack(buf *bytes.Buffer) error {
	return packSimple(RecCookie, false, c.Cookie, buf)
}

func (c Cookie) string() string {
	return fmt.Sprintf("--Cookie: %x\n", c.Cookie)
}

// Warning is the record type to send warnings to the other end. Put
// warning code in Code.
type Warning struct {
	RecordHdr
	Code uint16
}

func (w Warning) pack(buf *bytes.Buffer) error {
	return packSimple(RecWarning, true, w.Code, buf)
}

func (w Warning) string() string {
	return fmt.Sprintf("--Warning: %x\n", w.Code)
}

// Error is the record type to send errors to the other end. Put
// error code in Code.
type Error struct {
	RecordHdr
	Code uint16
}

func (e Error) pack(buf *bytes.Buffer) error {
	return packSimple(RecError, true, e.Code, buf)
}

func (e Error) string() string {
	return fmt.Sprintf("--Error: %x\n", e.Code)
}

// Algorithm is the record type for a list of AEAD algorithm we can use.
type Algorithm struct {
	RecordHdr
	Algo []uint16
}

func (a Algorithm) pack(buf *bytes.Buffer) error {
	return packSimple(RecAEAD, true, a.Algo, buf)
}

func (a Algorithm) string() string {
	var str = "--AEAD: \n"

	for i, alg := range a.Algo {
		algoStr := fmt.Sprintf("  #%v: %v\n", i, alg)
		str += algoStr
	}

	return str
}

func packSimple(t uint16, c bool, v any, buf *bytes.Buffer) error {
	newBuf := new(bytes.Buffer)
	err := binary.Write(newBuf, binary.BigEndian, v)
	if err != nil {
		return err
	}

	err = packHeader(t, c, buf, newBuf.Len())
	if err != nil {
		return err
	}

	_, err = buf.ReadFrom(newBuf)
	if err != nil {
		return err
	}

	return nil
}

func packHeader(t uint16, c bool, buf *bytes.Buffer, bodyLen int) error {
	var hdr RecordHdr

	hdr.Type = t

	if c {
		hdr.Type = setBit(hdr.Type, criticalBitPos)
	}

	hdr.BodyLen = uint16(bodyLen)

	err := hdr.pack(buf)
	if err != nil {
		return err
	}

	return nil
}

func setBit(n uint16, pos uint) uint16 {
	n |= 1 << pos
	return n
}

func checkBit(n uint16, pos uint) bool {
	val := n & (1 << pos)
	return val > 0
}
