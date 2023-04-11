package analysis

import "testing"

func TestStratumBarChart(t *testing.T) {
	srcPath := "D:/Desktop/Detect/domain/domain28_sta.csv"
	dstDir := "D:/Desktop/Detect/domain/graphic"
	prefix := "domain28_2_"
	err := StratumBarChart(srcPath, dstDir, prefix)
	if err != nil {
		t.Error(err)
	}
}
