package analysis

import "testing"

func TestCountry4StratumBarChart(t *testing.T) {
	srcPath := "D:/Desktop/Detect/domain/domain28_sta.csv"
	dstDir := "D:/Desktop/Detect/domain/graphic"
	prefix := "domain29_14_"
	err := Country4StratumBarChart(srcPath, dstDir, prefix, true)
	if err != nil {
		t.Error(err)
	}
}
