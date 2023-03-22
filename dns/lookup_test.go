package dns

import "testing"

const (
	src = "D:/Desktop/Detect/domain/domain6.txt"
	dst = "D:/Desktop/Detect/domain/domain6_ip.txt"
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
