package nts

import "testing"

func TestDialNTSKE(t *testing.T) {
	_, err := DialNTSKE("194.58.207.74", 4460, "sth2.nts.netnod.se")
	if err != nil {
		t.Error(err)
	}
}
