package tcp

import "testing"

func TestIsTLSEnabled(t *testing.T) {
	var tests = []struct {
		host       string
		port       int
		serverName string
		want       bool
	}{
		{"bilibili.com", 443, "bilibili.com", true},
		{"bilibili.com", 443, "baidu.com", false},
		{"bilibili.com", 444, "bilibili.com", false},
		{"baidu.com", 443, "", true},
		{"192.168.179.129", 4460, "nothing.com", false},
		{"194.58.207.74", 4460, "sth2.nts.netnod.se", true},
		{"194.58.207.74", 4460, "", true},
	}

	for _, test := range tests {
		if got := IsTLSEnabled(test.host, test.port, test.serverName); got != test.want {
			t.Errorf("IsTLSEnabled(%s, %d) = %t", test.host, test.port, got)
		}
	}
}
