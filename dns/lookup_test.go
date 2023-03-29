package dns

import "testing"

const (
	src = "D:/Desktop/TMP/毕设/NTP/第六阶段/official-server/domain-two-backup1.txt"
	dst = "D:/Desktop/TMP/毕设/NTP/第六阶段/official-server/ip-two-hk-backup1.txt"
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
