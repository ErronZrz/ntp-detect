package nts

import (
	"fmt"
	"testing"
)

func TestDetectNTSServer(t *testing.T) {
	host := "104.131.155.175"
	serverName := "ntp1.glypnod.com"
	payload, err := DetectNTSServer(host, serverName)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Print(payload.Lines())
}
