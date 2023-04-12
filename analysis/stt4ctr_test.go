package analysis

import "testing"

func TestStratumBarChart(t *testing.T) {
	srcPath := "D:/Desktop/Detect/domain/domain28_sta.csv"
	dstDir := "D:/Desktop/Detect/domain/graphic"
	prefix := "domain28_7_"
	err := Stratum4CountryBarChart(srcPath, dstDir, prefix)
	if err != nil {
		t.Error(err)
	}
}
