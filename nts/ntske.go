package nts

import (
	"active/utils"
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// KeyExchange stands for a NTS-KE connection.
type KeyExchange struct {
	hostPort string
	Conn     *tls.Conn
	reader   *bufio.Reader
	Meta     Data
	Debug    bool
}

// Data is the data in the key exchange process.
type Data struct {
	C2SKey []byte
	S2CKey []byte
	Server string
	Port   uint16
	Cookie [][]byte
	Algo   uint16
}

const (
	DefaultNTSKEPort = 4460
	DefaultNTPPort   = 123
	aesSivCmac256    = 0x0F
	alpnID           = "ntske/1"
	exportLabel      = "EXPORTER-network-time-security"
	keyLength        = 32
	timeout          = 5 * time.Second
)

// ExchangeMsg is a representation of a series of records to be sent
// to the peer.
type ExchangeMsg struct {
	Records []Record
}

// Print prints a description of all records in the key exchange
// message.
func (m *ExchangeMsg) Print() {
	for _, r := range m.Records {
		fmt.Print(r.string())
	}
}

// Pack allocates a buffer and packs all records into wire format in
// that buffer.
func (m *ExchangeMsg) Pack() (buf *bytes.Buffer, err error) {
	buf = new(bytes.Buffer)

	for _, r := range m.Records {
		err = r.pack(buf)
		if err != nil {
			return nil, err

		}
	}

	return buf, nil
}

// AddRecord adds a new record type to a key exchange message.
func (m *ExchangeMsg) AddRecord(rec Record) {
	m.Records = append(m.Records, rec)
}

// Connect connects to host:port and establishes an NTS-KE connection.
// If :port is left out, protocol default port is used.
// No further action is done.
func Connect(hostPort string, config *tls.Config, debug bool) (*KeyExchange, error) {
	config.NextProtos = []string{alpnID}

	ke := &KeyExchange{
		Debug:    debug,
		hostPort: hostPort,
	}

	_, _, err := net.SplitHostPort(ke.hostPort)
	if err != nil {
		if !strings.Contains(err.Error(), "missing port in address") {
			return nil, err
		}
		ke.hostPort = net.JoinHostPort(ke.hostPort, strconv.Itoa(DefaultNTSKEPort))
	}

	if ke.Debug {
		fmt.Printf("Connecting to KE server %v\n", ke.hostPort)
	}
	dialer := &net.Dialer{Timeout: timeout}
	ke.Conn, err = tls.DialWithDialer(dialer, "tcp", ke.hostPort, config)
	if err != nil {
		return nil, err
	}

	// Set default NTP server to the IP resolved and connected to for NTS-KE.
	// Handles multiple A records & possible lack of NTPv4ID Server Negotiation.
	ke.Meta.Server, _, err = net.SplitHostPort(ke.Conn.RemoteAddr().String())
	if err != nil {
		return nil, fmt.Errorf("unexpected remoteaddr issue: %s", err)
	}
	ke.Meta.Port = DefaultNTPPort

	if ke.Debug {
		fmt.Printf("Using resolved KE server as NTP default: %v\n",
			net.JoinHostPort(ke.Meta.Server, strconv.Itoa(int(ke.Meta.Port))))
	}
	ke.reader = bufio.NewReader(ke.Conn)

	state := ke.Conn.ConnectionState()
	if state.NegotiatedProtocol != alpnID {
		return nil, fmt.Errorf("server not speaking ntske/1")
	}

	return ke, nil
}

// Exchange initiates a client exchange using sane defaults on a
// connection already established with Connect(). After a successful
// run negotiated data is in ke.Meta.
func (ke *KeyExchange) Exchange() error {
	var msg ExchangeMsg

	var nextProto NextProto
	nextProto.NextProto = NTPv4ID
	msg.AddRecord(nextProto)

	var algo Algorithm
	algo.Algo = []uint16{aesSivCmac256}
	msg.AddRecord(algo)

	var end End
	msg.AddRecord(end)

	buf, err := msg.Pack()
	if err != nil {
		return err
	}

	reqData := buf.Bytes()
	fmt.Print(utils.PrintBytes(reqData, 16))
	_, err = ke.Conn.Write(reqData)
	if err != nil {
		return err
	}

	// Wait for response
	err = ke.Read()
	if err != nil {
		return err
	}

	return nil
}

// ExportKeys exports two extra sessions keys from the already
// established NTS-KE connection for use with NTS.
func (ke *KeyExchange) ExportKeys() error {
	// 4.3 in the RFC file https://tools.ietf.org/html/rfc8915#section-4.3
	//
	// The per-association context value SHALL consist of the following
	// five octets:
	//
	// The first two octets SHALL be zero (the Protocol ID for NTPv4ID).
	//
	// The next two octets SHALL be the Numeric Identifier of the
	// negotiated AEAD Algorithm in network byte order. Typically,
	// 0x0f for AES-SIV-CMAC-256.
	//
	// The final octet SHALL be 0x00 for the C2S key and 0x01 for the
	// S2C key.
	s2cContext := make([]byte, 4)
	binary.BigEndian.PutUint16(s2cContext[2:], ke.Meta.Algo)
	s2cContext = append(s2cContext, 0x01)

	c2sContext := make([]byte, 4)
	binary.BigEndian.PutUint16(c2sContext[2:], ke.Meta.Algo)
	c2sContext = append(c2sContext, 0x00)

	var err error
	state := ke.Conn.ConnectionState()

	ke.Meta.C2SKey, err = state.ExportKeyingMaterial(exportLabel, c2sContext, keyLength)
	if err != nil {
		return err
	}
	ke.Meta.S2CKey, err = state.ExportKeyingMaterial(exportLabel, s2cContext, keyLength)
	if err != nil {
		return err
	}
	return nil
}

// Read reads incoming NTS-KE messages until an End of Message record
// is received or an error occurs. It fills out the ke.Meta structure
// with negotiated data.
func (ke *KeyExchange) Read() error {
	var msg RecordHdr

	for {
		err := binary.Read(ke.reader, binary.BigEndian, &msg)
		if err != nil {
			return err
		}

		// C (Critical Bit): Determines the disposition of
		// unrecognized Records Types. Implementations which
		// receive a record with an unrecognized Records Type
		// MUST ignore the record if the Critical Bit is 0 and
		// MUST treat it as an error if the Critical Bit is 1.
		critical := checkBit(msg.Type, 15)

		// Get rid of Critical bit.
		msg.Type &^= 1 << 15

		if ke.Debug {
			fmt.Printf("Records type %v\n", msg.Type)
			if critical {
				fmt.Printf("Critical set\n")
			}
		}

		switch msg.Type {
		case RecEOM:
			// Check if we have complete data.
			// if len(ke.Meta.Cookie) == 0 || ke.Meta.Algo == 0 {
			// 	return errors.New("incomplete data")
			// }
			return nil

		case RecNextProto:
			var nextProto uint16
			err := binary.Read(ke.reader, binary.BigEndian, &nextProto)
			if err != nil {
				return errors.New("buffer overrun")
			}

		case RecAEAD:
			var aead uint16
			err := binary.Read(ke.reader, binary.BigEndian, &aead)
			if err != nil {
				return errors.New("buffer overrun")
			}

			ke.Meta.Algo = aead

		case RecCookie:
			cookie := make([]byte, msg.BodyLen)
			_, err := ke.reader.Read(cookie)
			if err != nil {
				return errors.New("buffer overrun")
			}

			ke.Meta.Cookie = append(ke.Meta.Cookie, cookie)

		case RecServer:
			address := make([]byte, msg.BodyLen)

			err := binary.Read(ke.reader, binary.BigEndian, &address)
			if err != nil {
				return errors.New("buffer overrun")
			}
			ke.Meta.Server = string(address)
			if ke.Debug {
				fmt.Printf("(got negotiated NTP server: %v)\n", ke.Meta.Server)
			}

		case RecPort:
			err := binary.Read(ke.reader, binary.BigEndian, &ke.Meta.Port)
			if err != nil {
				return errors.New("buffer overrun")
			}
			if ke.Debug {
				fmt.Printf("(got negotiated NTP port: %v)\n", ke.Meta.Port)
			}

		default:
			if critical {
				return fmt.Errorf("unknown record type %v with critical bit set", msg.Type)
			}

			// Swallow unknown record.
			unknownMsg := make([]byte, msg.BodyLen)
			err := binary.Read(ke.reader, binary.BigEndian, &unknownMsg)
			if err != nil {
				return errors.New("buffer overrun")
			}
		}
	}
}
