package analysis

import "testing"

func TestCountryAvgDelayBarChart(t *testing.T) {
	srcPath := "D:/Desktop/Detect/domain/domain28_sta.csv"
	dstDir := "D:/Desktop/Detect/domain/graphic"
	prefix := "new_syn_"
	err := CountryAvgDelayBarChart(srcPath, dstDir, prefix)
	if err != nil {
		t.Error(err)
	}
}
