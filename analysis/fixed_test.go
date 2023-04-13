package analysis

import "testing"

func TestStratumPrecisionBarChart(t *testing.T) {
	srcPath := "D:/Desktop/Detect/domain/domain28_sta.csv"
	dstDir := "D:/Desktop/Detect/domain/graphic"
	prefix := "domain29_13_"
	err := StratumPrecisionBarChart(srcPath, dstDir, prefix)
	if err != nil {
		t.Error(err)
	}
}
