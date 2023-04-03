package dns

import "testing"

const (
	src = "D:/Desktop/Detect/test/domain-test.txt"
	dst = "D:/Desktop/Detect/test/ip-test.txt"
)

func TestOutputDNS(t *testing.T) {
	err := OutputDNS(src, dst)
	if err != nil {
		t.Error(err)
	}
}

func TestDetectAfterDNS(t *testing.T) {
	err := DetectAfterDNS(src, dst)
	if err != nil {
		t.Error(err)
	}
}

func TestAsyncDetectAfterDNS(t *testing.T) {
	err := AsyncDetectAfterDNS(src, dst)
	if err != nil {
		t.Error(err)
	}
}

func TestTLSAfterDNS(t *testing.T) {
	err := TLSAfterDNS(src, dst)
	if err != nil {
		t.Error(err)
	}
}

func TestDetectAEADAfterDNS(t *testing.T) {
	err := DetectAEADAfterDNS(src, dst)
	if err != nil {
		t.Error(err)
	}
}
